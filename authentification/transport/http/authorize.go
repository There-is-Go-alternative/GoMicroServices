package http

import (
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

func AuthorizeHandler(cmd usecase.AuthorizeProto) fiber.Handler {
	return func(c *fiber.Ctx) {
		var token domain.Token

		err := c.BodyParser(&token)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(err.Error())
			return
		}

		payload, err := cmd(c.Context(), token)
		if err != nil {
			c.Status(http.StatusNotFound).JSON(err.Error())
			return
		}
		c.Status(http.StatusCreated).JSON(payload)
	}
}
