package http

import (
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewHttpServer(collection *mongo.Collection) *fiber.App {
	app := fiber.New()

	app.Post("/login", LoginHandler(usecase.Login(collection)))

	return app
}
