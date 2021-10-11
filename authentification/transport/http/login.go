package http

import (
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

func LoginHandler(cmd usecase.LoginProto) fiber.Handler {
	return func(c *fiber.Ctx) {
		var dto usecase.LoginDTO

		err := c.BodyParser(&dto)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(err)
			return
		}

		payload, err := cmd(c.Context(), dto)
		if err != nil {
			c.Status(http.StatusNotFound).JSON(err)
			return
		}
		c.Status(http.StatusCreated).JSON(payload)
	}
}
