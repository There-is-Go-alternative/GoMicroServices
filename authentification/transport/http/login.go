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
			ResponseError(c, http.StatusBadRequest, err)
			return
		}

		payload, err := cmd(c.Context(), dto)
		if err != nil {
			ResponseError(c, http.StatusNotFound, err)
			return
		}
		ResponseSuccess(c, http.StatusCreated, payload)
	}
}
