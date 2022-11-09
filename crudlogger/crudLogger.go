package crudlogger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/crudlogger/handler"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"

	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/crudlogger/utils"
)

func Set(c *fiber.App) {

	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

}

func MW(GetGridMODEL func(schema_id string) datagrid.Datagrid, GetMODEL func(schema_id string) dataform.Dataform) fiber.Handler {

	return func(c *fiber.Ctx) error {
		return handler.BodyDump(c, GetGridMODEL, GetMODEL)
	}

}
