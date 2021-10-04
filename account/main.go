//go:generate go run github.com/prisma/prisma-client-go generate --schema infra/database/prisma/schema.prisma

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/account/infra/database/prisma"
	infraHTTP "github.com/There-is-Go-alternative/GoMicroServices/account/infra/http"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/account/tests"
	privateHTTP "github.com/There-is-Go-alternative/GoMicroServices/account/transport/private/http"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	confPath        = flag.String("conf", os.Getenv("CONF_PATH"), "path to the json config file.")
	shutdownTimeOut = flag.Int("timeout", 2, "Time out for graceful shutdown, in seconds.")
	loadFixture     = flag.Bool("fixtures", false, "Time out for graceful shutdown, in seconds.")
	noLogColor      = flag.Bool("no-log-color", false, "if the logger should not print with color")
)

func main() {
	flag.Parse()
	if *confPath == "" {
		log.Fatal().Msg("Config path is empty")
	}
	if !*noLogColor {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger := zerolog.New(log.Logger)

	// Reading config from json file
	logger.Info().Str("stage", "setup").Msg("Parsing config ...")
	conf, err := config.ParseConfigFromPath(*confPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("problem when parsing config")
	}

	// Setup context to be notified when the program receive a signal
	logger.Info().Str("stage", "setup").Msg("Setting up context ...")
	signalCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	ctx, ctxCancel := context.WithCancel(signalCtx)

	// Initialising Account Database
	logger.Info().Str("stage", "setup").Msg("Setting up Account Database ...")
	//accountStorage := database.NewAccountMemMapStorage()
	//accountStorage, err := database.NewFirebaseRealTimeDB(ctx, database.DefaultConf)
	accountStorage, err := prisma.NewPrismaDB()
	if err != nil {
		log.Fatal().Err(err).Msg("When initialis")
	}

	// Initialising Account UseCase
	logger.Info().Str("stage", "setup").Msg("Setting up Account UseCase ...")
	accountUseCase := usecase.NewUseCase(&privateHTTP.AuthHTTP{}, accountStorage)

	// Initialising Gin Server
	logger.Info().Str("stage", "setup").Msg("Setting up Account Http handler ...")
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

	// launching each service in goroutine and catching error if any in errChan
	for _, fct := range services {
		// Launching the go routine and logging.
		go func(s service) {
			logger.Info().Str("stage", "runner").Msgf("Running %s", s.name)
			errChan <- s.fct(ctx)
		}(fct)
	}

	if *loadFixture {
		if err = tests.DefaultFixtures(ctx, accountStorage); err != nil {
			logger.Fatal().Err(err).Msg("Error when loading fixture.")
		}
	}

	// Waiting for a channel to receive something
	select {
	case <-ctx.Done():
		logger.Info().Str("stage", "runner").Msg("Context Canceled. Shutdown ...")
		time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
		return
	case err := <-errChan:
		logger.Error().Str("stage", "runner").Err(err).Msg(
			"An Error happened in a service, shutting down ...",
		)
		// Cancel context to shut down blocking services.
		ctxCancel()
		time.Sleep(time.Second * time.Duration(*shutdownTimeOut))
		os.Exit(1)
	}
}
