package http

import (
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
)

type UseCase interface {
	Login() usecase.LoginProto
	Authorize() usecase.AuthorizeProto
	Register() usecase.RegisterProto
	DeleteAccount() usecase.DeleteAccountProto
}

func NewHttpServer(u UseCase) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/login", LoginHandler(u.Login()))
	app.Post("/authorize", AuthorizeHandler(u.Authorize()))
	app.Post("/register", RegisterHandler(u.Register()))
	app.Delete("/unregister", DeleteAccountHandler(u.DeleteAccount()))

	return app
}
