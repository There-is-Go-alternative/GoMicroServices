package http

import (
	"context"
	"net/http"
	"time"

	"log"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

func LoginHandler(cmd usecase.LoginProto) fiber.Handler {
	return func(c *fiber.Ctx) {
		var dto usecase.LoginDTO

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		err := c.BodyParser(&dto)
		defer cancel()
		if err != nil {
			log.Printf("LoginInput invalid: %v\n", dto)
			c.Status(http.StatusNotFound).JSON(err)
			return
		}

		payload, err := cmd(ctx, dto)
		if err != nil {
			log.Printf("POST error login: %v\n", err)
			c.Status(http.StatusNotFound).JSON(err)
			return
		}
		c.Status(http.StatusCreated).JSON(payload)
	}
}
