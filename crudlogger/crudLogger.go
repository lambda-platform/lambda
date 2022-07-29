package crudlogger

import (
	"github.com/gofiber/fiber/v2"

	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/crudlogger/utils"
)

func Set(c *fiber.Ctx) {

	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

}

//
//func MW() echo.MiddlewareFunc {
//
//	return middleware.BodyDump(handler.CrudLogger)
//
//}
