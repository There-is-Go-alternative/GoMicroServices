package http

//func ApplyAdsRoutes(router *gin.Engine, uc AccountUseCase, conf *config.Config) {
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
