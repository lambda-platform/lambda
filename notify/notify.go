package notify

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/notify/handler"
	"github.com/lambda-platform/lambda/notify/utils"
)

func Set(e *echo.Echo) {

	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

	g := e.Group("/lambda/notify")
	/* ROUTES */
	g.GET("/new/:user_id", handler.GetNewNotifications, agentMW.IsLoggedInCookie)
	g.GET("/all/:user_id", handler.GetAllNotifications, agentMW.IsLoggedInCookie)
	g.GET("/seen/:id", handler.SetSeen, agentMW.IsLoggedInCookie)
	g.GET("/token/:user_id/:token", handler.SetToken, agentMW.IsLoggedInCookie)
	g.GET("/token", handler.SetTokenUrlParam)
	g.GET("/fcm", handler.Fcm)

}

func MW() echo.MiddlewareFunc {

	return middleware.BodyDump(handler.BodyDump)

}
