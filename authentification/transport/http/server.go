package http

import (
	"os"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

func NewHttpServer(u usecase.LoginProto) {
	app := fiber.New()
	//httpHandler := http.Handler{Logger: log.New(os.Stdout).With().Logger()}
	app.Post("/login", LoginHandler(usecase.Login()))

	app.Listen(os.Getenv("PORT"))
}
