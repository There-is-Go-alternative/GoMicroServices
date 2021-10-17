package http

import (
	"context"
	"fmt"
	netHTTP "net/http"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/ads/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/ads/usecase"
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
	CreateAd() usecase.CreateAdCmd
	UpdateAd() usecase.UpdateAdCmd
	GetAdById() usecase.GetAdByIdCmd
	GetAllAds() usecase.GetAllAdsCmd
	DeleteAd() usecase.DeleteAdCmd
	SearchAd() usecase.SearchAdCmd
}

// TODO: change database by future Database interface
func NewHttpServer(uc useCase, conf *config.Config) *Server {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		AllowWildcard:    true,
		MaxAge:           12 * time.Hour,
	}))

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
		ad.PATCH("/:id", adHandler.UpdateAdHandler(uc.UpdateAd()))
		ad.GET("/search/", adHandler.SearchAdHandler(uc.SearchAd()))
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
