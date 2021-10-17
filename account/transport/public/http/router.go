package http

import (
	"fmt"
	netHTTP "net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AccountUseCase interface {
	CreateAccount() usecase.CreateAccountCmd
	GetAccountByID() usecase.GetAccountByIDCmd
	GetAllAccounts() usecase.GetAllAccountsCmd
	DeleteAccount() usecase.DeleteAccountCmd
	PatchAccount() usecase.PatchAccountCmd
	UpdateAccount() usecase.UpdateAccountCmd
}

func ApplyAccountRoutes(router *gin.Engine, uc AccountUseCase, conf *config.Config) {
	accountHandler := NewAccountHandler(conf.APIKey)

	// Configuring CORS
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.Status(netHTTP.StatusOK)
	})
	// Grouping Account routes with url specified in config (I.E: 'account')
	account := router.Group(fmt.Sprintf("/%s", conf.AccountEndpoint))
	{
		account.POST("/", accountHandler.CreateAccountHandler(uc.CreateAccount()))
		account.GET("/", accountHandler.Authorize(), accountHandler.GetAllAccountsHandler(uc.GetAllAccounts()))
		account.GET("/:id", accountHandler.Authorize(), accountHandler.GetAccountsByIDHandler(uc.GetAccountByID()))
		account.PATCH("/:id", accountHandler.Authorize(), accountHandler.PatchAccountHandler(uc.PatchAccount()))
		account.PUT("/:id", accountHandler.Authorize(), accountHandler.PutAccountHandler(uc.UpdateAccount()))
		account.DELETE("/:id", accountHandler.Authorize(), accountHandler.DeleteAccountHandler(uc.DeleteAccount()))
	}
}
