package http

import (
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

func RegisterHandler(cmd usecase.RegisterProto) fiber.Handler {
	return func(c *fiber.Ctx) {
		var register domain.Auth

		err := c.BodyParser(&register)
		if err != nil {
			ResponseError(c, http.StatusBadRequest, err)
			return
		}

		payload, err := cmd(c.Context(), register)
		if err != nil {
			c.Status(http.StatusNotFound).JSON(err)
			return
		}
		ResponseSuccess(c, http.StatusCreated, payload)
	}
}
