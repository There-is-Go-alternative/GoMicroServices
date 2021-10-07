package http

import (
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

type UseCase interface {
	Login() usecase.LoginProto
}

func NewHttpServer(u UseCase) *fiber.App {
	app := fiber.New()

	app.Post("/login", LoginHandler(u.Login()))

	return app
}
