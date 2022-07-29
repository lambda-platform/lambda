package moqup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/moqup/handler"
)

func Set(e *fiber.App) {

	e.Get("/pages/moqup/:id", handler.Moqup)

}
