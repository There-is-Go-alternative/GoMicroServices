package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	infraDB "github.com/There-is-Go-alternative/GoMicroServices/transactions/infra/database"
	"github.com/There-is-Go-alternative/GoMicroServices/transactions/internal/config"
	privateHTTP "github.com/There-is-Go-alternative/GoMicroServices/transactions/transport/private/http"
	infraHTTP "github.com/There-is-Go-alternative/GoMicroServices/transactions/transport/public"
	"github.com/There-is-Go-alternative/GoMicroServices/transactions/usecase"
	log "github.com/sirupsen/logrus"
)

var (
	confPath        = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
	loadFixture     = flag.Bool("fixtures", false, "Time out for graceful shutdown, in seconds.")
	shutdownTimeOut = flag.Int("timeout", 2, "Time out for graceful shutdown, in seconds.")
)

func main() {
	flag.Parse()

	if *confPath == "" {
		log.Debug("Config path is empty")
	}

	logger := log.New()

	setupContext := logger.WithFields(log.Fields{"stage": "setup"})

	// Reading config from json file
	setupContext.Info("Parsing config ...")
	conf, err := config.NewConfig(*confPath)
	if err != nil {
		setupContext.Fatalf("problem when parsing config: %v", err)
	}

	// Setup context to be notified when the program receive a signal
	setupContext.Info("Setting up context ...")
	signalCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	ctx, ctxCancel := context.WithCancel(signalCtx)

	// Initialising Account Database
	setupContext.Info("Setting up Account Database ...")
	//accountStorage := database.NewAccountMemMapStorage()
	//accountStorage, err := database.NewFirebaseRealTimeDB(ctx, database.DefaultConf)
	storage := infraDB.NewPrismaDB()

	// Initialising Auth Service connector
	setupContext.Info("Setting up Auth service ...")
	accountService := privateHTTP.NewAccountHTTP(conf.AccountURL)

	// Initialising Balance Service connector
	setupContext.Info("Setting up Balance service ...")
	fundsService := privateHTTP.NewFundsHTTP(conf.FundsURL, conf.APIKey)

	// Initialising Balance Service connector
	setupContext.Info("Setting up Balance service ...")
	adsService := privateHTTP.NewAdsHTTP(conf.AdsURL, conf.APIKey)

	// Initialising Account UseCase
	setupContext.Info("Setting up Account UseCase ...")
	accountUseCase := usecase.NewUseCase(accountService, adsService, fundsService, storage)

	// Initialising Gin Server
	setupContext.Info("Setting up Account Http handler ...")
	ginServer := infraHTTP.NewHttpServer(accountUseCase, conf)

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
	runnerContext := logger.WithFields(log.Fields{"stage": "setup"})

	// launching each service in goroutine and catching error if any in errChan
	for _, fct := range services {
		// Launching the go routine and logging.
		go func(s service) {
			runnerContext.Infof("Running %s", s.name)
			errChan <- s.fct(ctx)
		}(fct)
	}

	// Waiting for a channel to receive something
	select {
	case <-ctx.Done():
		runnerContext.Info("Context Canceled. Shutdown ...")
		time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
		return
	case err := <-errChan:
		runnerContext.Errorf("An Error happened in a service, shutting down ... (%v)", err)
		// Cancel context to shut down blocking services.
		ctxCancel()
		time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
		os.Exit(1)
	}
}
