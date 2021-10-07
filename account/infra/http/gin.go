package http

import (
	"context"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/config"
	transportPublicHTTP "github.com/There-is-Go-alternative/GoMicroServices/account/transport/public/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	netHTTP "net/http"
)

type Server struct {
	Engine *netHTTP.Server
	logger *logrus.Entry
}

func NewHttpServer(uc transportPublicHTTP.AccountUseCase, conf *config.Config, logger *logrus.Logger) *Server {
	gin.DefaultWriter = logger.Writer()
	router := gin.Default()

	// Configuring CORS
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.Status(netHTTP.StatusOK)
	})

	transportPublicHTTP.ApplyAccountRoutes(router, uc, conf)

	// Grouping Account routes with url specified in config (I.E: 'account')
	return &Server{
		Engine: &netHTTP.Server{
			Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.HTTPPort),
			Handler:  router,
			ErrorLog: log.New(logger.Writer(), "", 0),
		},
		logger: logger.WithFields(logrus.Fields{"service": "HTTP gin server"}),
	}
}

func (s Server) Run(ctx context.Context) (err error) {
	s.logger.Infof("Running gin HTTP server on %v ...", s.Engine.Addr)
	errc := make(chan error)
	go func() {
		errc <- s.Engine.ListenAndServe()
	}()
	select {
	case err = <-errc:
		return
	case <-ctx.Done():
		if err = s.Engine.Shutdown(ctx); err != nil && err != context.Canceled {
			s.logger.Errorf("Error happened when server forced to shutdown: %v", err)
			return
		}
	}
	return
}
