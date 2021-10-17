package http

import (
	"net/http"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/domain"
	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

func DeleteAccountHandler(cmd usecase.DeleteAccountProto) fiber.Handler {
	return func(c *fiber.Ctx) {
		var account domain.Auth

		err := c.BodyParser(&account)
		if err != nil {
			ResponseError(c, http.StatusBadRequest, err)
			return
		}

		payload, err := cmd(c.Context(), account)
		if err != nil {
			c.Status(http.StatusNotFound).JSON(err)
			return
		}
		ResponseSuccess(c, http.StatusCreated, payload)
	}
}
