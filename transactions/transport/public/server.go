package transport

import (
	"context"
	"fmt"
	netHTTP "net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/transactions/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/transactions/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Engine *netHTTP.Server
	logger zerolog.Logger
}

func NewHttpServer(uc *usecase.UseCase, conf *config.Config) *Server {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.Status(netHTTP.StatusOK)
	})
	handler := NewTransactionHandler(uc)
	api := router.Group("/api/v1/")
	{
		api.GET("/", handler.GetAll(uc.GetAll()))
		api.OPTIONS("/", handler.GetByDateRange(uc.GetByDateRange()))

		api.GET("/:id", handler.GetById(uc.GetById()))
		api.POST("/buy/:id", handler.Register(uc.Register()))

		buyed := api.Group("/buyed/")
		{
			buyed.GET("/:user_id", handler.GetAllBuyedByUser(uc.GetByBuyerId()))
			buyed.OPTIONS("/:user_id", handler.GetBuyedByDateRangeUser(uc.GetByBuyerIdDateRange()))
		}

		selled := api.Group("/selled/")
		{
			selled.GET("/:user_id", handler.GetAllSelledByUser(uc.GetBySellerId()))
			selled.OPTIONS("/:user_id", handler.GetSelledByDateRangeUser(uc.GetBySellerIdDateRange()))
		}
	}

	return &Server{
		Engine: &netHTTP.Server{
			Addr:    fmt.Sprintf("%s:%s", conf.Host, conf.HTTPPort),
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
