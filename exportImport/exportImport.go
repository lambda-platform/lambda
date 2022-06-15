package exportImport

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/exportImport/handlers"
)

func Set(e *echo.Echo) {
	e.GET("crud/export", handlers.Export, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	e.GET("crud/import/:file", handlers.Import, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
}
