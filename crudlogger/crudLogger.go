package crudlogger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/crudlogger/handler"
	"github.com/lambda-platform/lambda/crudlogger/utils"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
)

func Set(a *fiber.App) {

	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

	a.Post("/crud_log/history", agentMW.IsLoggedIn(), handler.CrudLogHistory)
}

func MW(GetGridMODEL func(schemaId string) datagrid.Datagrid, GetMODEL func(schemaId string) dataform.Dataform) fiber.Handler {

	return func(c *fiber.Ctx) error {
		return handler.BodyDump(c, GetGridMODEL, GetMODEL)
	}

}
