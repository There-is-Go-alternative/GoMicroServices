package http

import (
	"context"
	"fmt"
	netHTTP "net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Engine *netHTTP.Server
	logger zerolog.Logger
}

type useCase interface {
	CreateAd() usecase.CreateAdCmd
	GetAdById() usecase.GetAdByIdCmd
	GetAllAds() usecase.GetAllAdsCmd
	DeleteAd() usecase.DeleteAdCmd
}

// TODO: change database by future Database interface
func NewHttpServer(uc useCase, conf *config.Config) *Server {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.Status(netHTTP.StatusOK)
	})
	adHandler := NewAdHandler()
	// Grouping Ad routes with url specified in config (I.E: 'ad')
	ad := router.Group(fmt.Sprintf("/%s", conf.AdEndpoint))
	{
		ad.POST("/", adHandler.CreateAdHandler(uc.CreateAd()))
		ad.GET("/", adHandler.GetAdsHandler(uc.GetAllAds()))
		ad.GET("/:id", adHandler.GetAdsByIDHandler(uc.GetAdById()))
		ad.DELETE("/:id", adHandler.DeleteAdHandler(uc.DeleteAd()))
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
