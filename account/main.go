package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"

	infraDB "github.com/There-is-Go-alternative/GoMicroServices/account/infra/database"
	infraHTTP "github.com/There-is-Go-alternative/GoMicroServices/account/infra/http"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/account/tests"
	privateHTTP "github.com/There-is-Go-alternative/GoMicroServices/account/transport/private/http"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
	_ "github.com/joho/godotenv/autoload"
)

var (
	confPath        = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
	shutdownTimeOut = flag.Int("timeout", 2, "Time out for graceful shutdown, in seconds.")
	loadFixture     = flag.Bool("fixtures", false, "Time out for graceful shutdown, in seconds.")
	logFormatter    = flag.String("log-formatter", "text", "Which formatter the logger must use")
)

func main() {
	flag.Parse()
	if *confPath == "" {
		log.Fatal("Config path is empty")
	}

	logger := log.New()

	switch *logFormatter {
	case "json":
		logger.SetFormatter(&log.JSONFormatter{})
	case "text":
		break
	default:
		log.Fatalf("Log formatter is not one of possible, got: %s", *logFormatter)
	}

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
	accountStorage, err := infraDB.NewPrismaDB()
	if err != nil {
		setupContext.Fatal("When initialising Acccount storage: %v", err)
	}

	// Initialising Account UseCase
	setupContext.Info("Setting up Account UseCase ...")
	accountUseCase := usecase.NewUseCase(&privateHTTP.AuthHTTP{}, accountStorage, logger)

	// Initialising Gin Server
	setupContext.Info("Setting up Account Http handler ...")
	ginServer := infraHTTP.NewHttpServer(accountUseCase, conf, logger)

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

	if *loadFixture {
		if err = tests.DefaultFixtures(ctx, accountStorage); err != nil {
			runnerContext.Fatalf("Error when loading fixture: %v", err)
		}
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
