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
			ResponseError(c, http.StatusBadRequest, err)
			return
		}

		payload, err := cmd(c.Context(), token)
		if err != nil {
			c.Status(http.StatusNotFound).JSON(err)
			return
		}
		ResponseSuccess(c, http.StatusCreated, payload)
	}
}
