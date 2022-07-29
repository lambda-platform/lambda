package notify

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/notify/handler"
	"github.com/lambda-platform/lambda/notify/utils"
)

func Set(c *fiber.App) {

	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

	g := c.Group("/lambda/notify")
	/* ROUTES */
	g.Get("/new/:user_id", agentMW.IsLoggedIn(), handler.GetNewNotifications)
	g.Get("/all/:user_id", agentMW.IsLoggedIn(), handler.GetAllNotifications)
	g.Get("/seen/:id", agentMW.IsLoggedIn(), handler.SetSeen)
	g.Get("/token/:user_id/:token", agentMW.IsLoggedIn(), handler.SetToken)
	g.Get("/token", handler.SetTokenUrlParam)
	g.Get("/fcm", handler.Fcm)

}

//func MW() echo.MiddlewareFunc {
//
//	return middleware.BodyDump(handler.BodyDump)
//
//}
