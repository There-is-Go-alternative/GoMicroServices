package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/infra/database"
	"github.com/There-is-Go-alternative/GoMicroServices/funds/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/funds/transport/private/http"
	"github.com/There-is-Go-alternative/GoMicroServices/funds/usecase"
	log "github.com/sirupsen/logrus"
)

var (
	confPath        = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
	shutdownTimeOut = flag.Int("timeout", 2, "Time out for graceful shutdown, in seconds.")
)

func main() {
	flag.Parse()
	if *confPath == "" {
		log.Fatal("config path is empty")
	}

	conf, err := config.ParseConfigFromPath(*confPath)
	if err != nil {
		log.Fatalf("problem when parsing config: %v", err)
	}

	fmt.Println(conf)

	signalCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	ctx, ctxCancel := context.WithCancel(signalCtx)

	storage := database.NewPrismaDB()

	if err := storage.Connect(); err != nil {
		log.Fatalf("problem while connecting to the DB: %v", err)
	}

	defer func() {
		if err := storage.Disconnect(); err != nil {
			panic(err)
		}
	}()

	fundsUseCase := usecase.NewUseCase(&http.AuthHTTP{}, storage)
	ginServer := http.NewHttpServer(fundsUseCase, conf)

	type service struct {
		name string
		fct  func(context.Context) error
	}
	services := []service{
		{
			name: "Http Server",
			fct:  ginServer.Run,
		},
	}

	errChan := make(chan error)

	for _, fct := range services {
		// Launching the go routine and logging.
		go func(s service) {
			log.WithFields(log.Fields{
				"stage": "runner",
			}).Infof("Running %s", s.name)
			errChan <- s.fct(ctx)
		}(fct)
	}

	// Waiting for a channel to receive something
	select {
	case <-ctx.Done():
		log.WithFields(log.Fields{
			"stage": "runner",
		}).Info("Context Canceled. Shutdown ...")
		time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
		return
	case err := <-errChan:
		log.WithFields(log.Fields{
			"stage": "runner",
		}).Errorf("An Error happend in a service: %s", err)
		// Cancel context to shut down blocking services.
		ctxCancel()
		time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
		os.Exit(1)
	}
}
