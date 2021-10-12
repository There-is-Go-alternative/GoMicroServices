package http

import (
	"github.com/gofiber/fiber"
)

func ResponseError(c *fiber.Ctx, code int, message error) {
	c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message.Error(),
	})
}

func ResponseSuccess(c *fiber.Ctx, code int, message interface{}) {
	c.Status(code).JSON(fiber.Map{
		"success": true,
		"data":    message,
	})
}
