package http

import (
	"fmt"
	netHTTP "net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/config"
	"github.com/There-is-Go-alternative/GoMicroServices/account/usecase"
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

// TODO: change database by future Database interface
func ApplyAccountRoutes(router *gin.Engine, uc AccountUseCase, conf *config.Config) {
	accountHandler := NewAccountHandler(conf.APIKey)
	// Grouping Account routes with url specified in config (I.E: 'account')
	account := router.Group(fmt.Sprintf("/%s", conf.AccountEndpoint))
	{
		account.POST("/", accountHandler.CreateAccountHandler(uc.CreateAccount()))
		account.GET("/", accountHandler.GetAllAccountsHandler(uc.GetAllAccounts()))
		account.GET("/:id", accountHandler.GetAccountsByIDHandler(uc.GetAccountByID()))
		account.PATCH("/:id", accountHandler.PatchAccountHandler(uc.PatchAccount()))
		account.PUT("/:id", accountHandler.PutAccountHandler(uc.UpdateAccount()))
		account.DELETE("/:id", accountHandler.DeleteAccountHandler(uc.DeleteAccount()))
		account.GET("/test", accountHandler.Authorize(), func(c *gin.Context) {
			uuid, _ := domain.NewAccountID()
			c.JSON(netHTTP.StatusOK, gin.H{
				"success": true,
				"data": domain.Account{
					ID: *uuid,
					Firstname: "COUCOU",
					Lastname:  "ENTHONNE",
				},
			})
		})
	}
}

//// TODO: change database by future Database interface
//func Apply(router *gin.Engine, uc AccountUseCase, conf *config.Config) {
//	accountHandler := NewAccountHandler()
//	// Grouping Account routes with url specified in config (I.E: 'account')
//	account := router.Group(fmt.Sprintf("/%s", conf.AccountEndpoint))
//	{
//		account.POST("/", accountHandler.CreateAccountHandler(uc.CreateAccount()))
//		account.GET("/", accountHandler.GetAllAccountsHandler(uc.GetAllAccounts()))
//		account.GET("/:id", accountHandler.GetAccountsByIDHandler(uc.GetAccountByID()))
//		account.PATCH("/:id", accountHandler.PatchAccountHandler(uc.PatchAccount()))
//		account.PUT("/:id", accountHandler.PutAccountHandler(uc.UpdateAccount()))
//		account.DELETE("/:id", accountHandler.DeleteAccountHandler(uc.DeleteAccount()))
//	}
//}
