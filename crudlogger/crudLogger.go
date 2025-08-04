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

	a.Get("/api/logs", agentMW.IsLoggedIn(), handler.GetLogs)
	a.Get("/api/log-categories", agentMW.IsLoggedIn(), handler.GetLogCategories)
	a.Get("/api/log-actions", agentMW.IsLoggedIn(), handler.GetLogActions)
	a.Get("/api/log-chart-user", agentMW.IsLoggedIn(), handler.GetLogChartUser)
	a.Get("/api/log-chart-action", agentMW.IsLoggedIn(), handler.GetLogChartAction)
	a.Get("/api/log-chart-hourly", agentMW.IsLoggedIn(), handler.GetLogChartHourly)
	a.Get("/api/log-chart-monthly", agentMW.IsLoggedIn(), handler.GetLogChartMonthly)
	a.Get("/api/log-chart-type", agentMW.IsLoggedIn(), handler.GetLogChartType)
	a.Get("/api/log-chart-category", agentMW.IsLoggedIn(), handler.GetLogChartCategory)
}

func MW(GetGridMODEL func(schemaId string) datagrid.Datagrid, GetMODEL func(schemaId string) dataform.Dataform) fiber.Handler {

	return func(c *fiber.Ctx) error {
		return handler.BodyDump(c, GetGridMODEL, GetMODEL)
	}

}
