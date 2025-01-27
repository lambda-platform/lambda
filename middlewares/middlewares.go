package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func Set(app *fiber.App) {
	/*
		|----------------------------------------------
		| Useful MIDDLEWARES
		|----------------------------------------------
	*/
	ThumbMiddleware(app)

}
