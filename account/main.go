package main

import (
	"context"
	"flag"
	"github.com/There-is-Go-alternative/GoMicroServices/account/infra/database"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/account/tests"
	"github.com/There-is-Go-alternative/GoMicroServices/account/transport/public/http"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	confPath        = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
	shutdownTimeOut = flag.Int("timeout", 2, "Time out for graceful shutdown, in seconds.")
	loadFixture     = flag.Bool("fixtures", false, "Time out for graceful shutdown, in seconds.")
)

func main() {
	flag.Parse()
	if *confPath == "" {
		log.Fatal("Config path is empty")
	}

	// Reading config from json file
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Parsing config ...")
	conf, err := config.ParseConfigFromPath(*confPath)
	if err != nil {
		log.Fatalf("probleme when parsing config: %v", err)
	}

	// Setup context to be notified when the program receive a signal
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Setting up context ...")
	signalCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	ctx, ctxCancel := context.WithCancel(signalCtx)

	// Initialising Account Database
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Setting up Account Database ...")
	//accountStorage := database.NewAccountMemMapStorage()
	accountStorage, err := database.NewFirebaseRealTimeDB(ctx, database.DefaultConf)
	if err != nil {
		log.Fatal(err)
	}

	// Initialising Account UseCase
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Setting up Account UseCase ...")
	accountUseCase := usecase.NewUseCase(accountStorage)

	// Initialising Gin Server
	log.WithFields(log.Fields{
		"stage": "setup",
	}).Info("Setting up Account Http handler ...")
	ginServer := http.NewHttpServer(accountUseCase, conf)

	// Setup blocking service that must be run in parallel inside a go routine
	//  I.E: Http server, kafka consumer, ...
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

	// Setup an error channel
	errChan := make(chan error)

	// launching each service in goroutine and catching error if any in errChan
	for _, fct := range services {
		// Launching the go routine and logging.
		go func(s service) {
			log.WithFields(log.Fields{
				"stage": "runner",
			}).Infof("Running %s", s.name)
			errChan <- s.fct(ctx)
		}(fct)
	}

	if *loadFixture {
		if err = tests.DefaultFixtures(ctx, accountStorage); err != nil {
			log.Fatal(err)
		}
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
