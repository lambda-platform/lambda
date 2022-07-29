package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetBody(c *fiber.Ctx) []byte {
	var bodyBytes []byte

	if c.Request().Body != nil {
		bodyBytes = c.Body()
	}

	return bodyBytes
}
