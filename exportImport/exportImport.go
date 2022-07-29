package exportImport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/exportImport/handlers"
)

func Set(e *fiber.App) {
	e.Get("crud/export", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.Export)
	e.Get("crud/import/:file", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.Import)
}
