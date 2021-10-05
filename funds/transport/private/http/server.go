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
	Set() usecase.SetCmd
	IncreaseByUser() usecase.IncreaseByUserCmd
	DecreaseByUser() usecase.DecreaseByUserCmd
	SetByUser() usecase.SetByUserCmd
}

func NewHttpServer(uc useCase, conf *config.Config, auth *usecase.AuthService) *Server {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.Status(netHTTP.StatusOK)
	})
	fundsHandler := NewFundsHandler(auth)
	api := router.Group("/api/v1/")
	{
		api.Use(fundsHandler.ValidateToken)

		api.POST("/:id", fundsHandler.CreateFundsHandler(uc.Create()))
		api.GET("/", fundsHandler.GetAllFundsHandler(uc.All()))
		api.GET("/:id", fundsHandler.GetFundsByIDHandler(uc.GetByID()))
		api.GET("/user/:id", fundsHandler.GetFundsByUserIDHandler(uc.GetByUserID()))
		api.DELETE("/:id", fundsHandler.DeleteFundsByIDHandler(uc.DeleteByID()))
		api.DELETE("/user/:id", fundsHandler.DeleteFundsByUserIDHandler(uc.DeleteByUserID()))

		balance := api.Group("/balance/")
		{
			balance.POST("/increase/:id", fundsHandler.IncreaseFundsHandler(uc.Increase()))
			balance.POST("/decrease/:id", fundsHandler.DecreaseFundsHandler(uc.Decrease()))
			balance.POST("/set/:id", fundsHandler.SetFundsHandler(uc.Set()))
			balance.POST("/increase/user/:id", fundsHandler.IncreaseFundsByUserHandler(uc.IncreaseByUser()))
			balance.POST("/decrease/user/:id", fundsHandler.DecreaseFundsByUserHandler(uc.DecreaseByUser()))
			balance.POST("/set/user/:id", fundsHandler.SetFundsByUserHandler(uc.SetByUser()))
		}
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
