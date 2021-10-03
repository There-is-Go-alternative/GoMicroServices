package http

import (
	"context"
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	netHTTP "net/http"
)

type Server struct {
	Engine *netHTTP.Server
	logger zerolog.Logger
}

type useCase interface {
	CreateAccount() usecase.CreateAccountCmd
	GetAccountByID() usecase.GetAccountByIDCmd
	GetAllAccounts() usecase.GetAllAccountsCmd
	DeleteAccount() usecase.DeleteAccountCmd
	PatchAccount() usecase.PatchAccountCmd
	UpdateAccount() usecase.UpdateAccountCmd
}

// TODO: change database by future Database interface
func NewHttpServer(uc useCase, conf *config.Config) *Server {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.Status(netHTTP.StatusOK)
	})
	accountHandler := NewAccountHandler()
	// Grouping Account routes with url specified in config (I.E: 'account')
	account := router.Group(fmt.Sprintf("/%s", conf.AccountEndpoint))
	{
		account.POST("/", accountHandler.CreateAccountHandler(uc.CreateAccount()))
		account.GET("/", accountHandler.GetAllAccountsHandler(uc.GetAllAccounts()))
		account.GET("/:id", accountHandler.GetAccountsByIDHandler(uc.GetAccountByID()))
		account.PATCH("/:id", accountHandler.PatchAccountHandler(uc.PatchAccount()))
		account.PUT("/:id", accountHandler.PutAccountHandler(uc.UpdateAccount()))
		account.DELETE("/:id", accountHandler.DeleteAccountHandler(uc.DeleteAccount()))
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
