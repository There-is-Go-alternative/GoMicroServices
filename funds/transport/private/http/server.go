package http

import (
	"context"
	"fmt"
	netHTTP "net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/funds/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/funds/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Engine *netHTTP.Server
	logger zerolog.Logger
}

type useCase interface {
	Create() usecase.CreateFundsCmd
	All() usecase.AllCmd
	GetByID() usecase.GetByIDCmd
	GetByUserID() usecase.GetByUserIDCmd
	DeleteByID() usecase.DeleteByIDCmd
	DeleteByUserID() usecase.DeleteByUserIDCmd
	Increase() usecase.IncreaseCmd
	Decrease() usecase.DecreaseCmd
}

func NewHttpServer(uc useCase, conf *config.Config) *Server {
	router := gin.Default()

	// Configuring CORS
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.Status(netHTTP.StatusOK)
	})
	fundsHandler := NewFundsHandler()
	account := router.Group("/funds")
	{
		account.POST("/:id", fundsHandler.CreateFundsHandler(uc.Create()))
		account.GET("/", fundsHandler.GetAllFundsHandler(uc.All()))
		account.GET("/:id", fundsHandler.GetFundsByIDHandler(uc.GetByID()))
		account.GET("/user/:id", fundsHandler.GetFundsByUserIDHandler(uc.GetByUserID()))
		account.DELETE("/:id", fundsHandler.DeleteFundsByIDHandler(uc.DeleteByID()))
		account.DELETE("/user/:id", fundsHandler.DeleteFundsByUserIDHandler(uc.DeleteByUserID()))
		account.POST("/i/:id", fundsHandler.IncreaseFundsHandler(uc.Increase()))
		account.POST("/d/:id", fundsHandler.DecreaseFundsHandler(uc.Decrease()))
	}

	return &Server{
		Engine: &netHTTP.Server{
			Addr:    fmt.Sprintf("%s:%s", conf.Host, conf.Port),
			Handler: router,
		},
		logger: log.With().Str("service", "HTTP gin server").Logger(),
	}
}

func (s Server) Run(ctx context.Context) (err error) {
	s.logger.Info().Msg("Running gin HTTP server ...")
	errc := make(chan error)
	go func() {
		errc <- s.Engine.ListenAndServe()
	}()
	select {
	case err = <-errc:
		return
	case <-ctx.Done():
		if err = s.Engine.Shutdown(ctx); err != nil && err != context.Canceled {
			s.logger.Error().Msgf("Error happened when server forced to shutdown: %v", err)
			return
		}
	}
	return
}
