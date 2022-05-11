package crudlogger

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/crudlogger/handler"
	"github.com/lambda-platform/lambda/crudlogger/utils"
)

func Set(e *echo.Echo) {

	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

}

func MW() echo.MiddlewareFunc {

	return middleware.BodyDump(handler.CrudLogger)

}
