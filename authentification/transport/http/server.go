package http

import (
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

type UseCase interface {
	Login() usecase.LoginProto
	Authorize() usecase.AuthorizeProto
	Register() usecase.RegisterProto
}

func NewHttpServer(u UseCase) *fiber.App {
	app := fiber.New()

	app.Post("/login", LoginHandler(u.Login()))
	app.Post("/authorize", AuthorizeHandler(u.Authorize()))
	app.Post("/register", RegisterHandler(u.Register()))
	//app.Post("/unregister", RegisterHandler(u.Register()))

	return app
}
