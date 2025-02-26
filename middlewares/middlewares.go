package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func Set(app *fiber.App) {
	// Rate limiter middleware
	//app.Use(limiter.New(limiter.Config{
	//	Expiration: 30 * time.Second, //
	//	Max:        1000,             //
	//}))

	/*
		|----------------------------------------------
		| Useful MIDDLEWARES
		|----------------------------------------------
	*/
	ThumbMiddleware(app)

}
