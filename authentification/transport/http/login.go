package http

import (
	"context"
	"net/http"
	"time"

	"github.com/There-is-Go-alternative/GoMicroServices/authentification/usecase"
	"github.com/gofiber/fiber"
)

type conf struct {
	LoginExp time.Duration
}

var defaultConf = &conf{
	LoginExp: 100 * time.Second,
}

func LoginHandler(cmd usecase.LoginProto) fiber.Handler {
	return func(c *fiber.Ctx) {
		var dto usecase.LoginDTO

		ctx, cancel := context.WithTimeout(c.Context(), defaultConf.LoginExp)
		err := c.BodyParser(&dto)
		defer cancel()
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(err)
			return
		}

		payload, err := cmd(ctx, dto)
		if err != nil {
			c.Status(http.StatusNotFound).JSON(err)
			return
		}
		c.Status(http.StatusCreated).JSON(payload)
	}
}
