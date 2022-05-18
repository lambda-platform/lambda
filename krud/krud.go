package krud

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/krud/handlers"
	"github.com/lambda-platform/lambda/krud/krudMW"
	"github.com/lambda-platform/lambda/krud/utils"
)

func Set(e *echo.Echo, GetGridMODEL func(schema_id string) datagrid.Datagrid, GetMODEL func(schema_id string) dataform.Dataform, krudMiddleWares []echo.MiddlewareFunc, KrudWithPermission bool) {
	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

	g := e.Group("/lambda/krud")
	//g.Use(agentMW.IsLoggedInCookie)

	if len(krudMiddleWares) >= 1 {
		for _, krudMiddleWare := range krudMiddleWares {
			g.Use(krudMiddleWare)
		}
	}

	if KrudWithPermission {
		g.POST("/update-row/:schemaId", handlers.UpdateRow(GetGridMODEL), agentMW.IsLoggedInCookie, krudMW.PermissionDelete)
		g.POST("/:schemaId/:action", handlers.Crud(GetMODEL), agentMW.IsLoggedInCookie, krudMW.PermissionCreate)
		g.POST("/:schemaId/:action/:id", handlers.Crud(GetMODEL), agentMW.IsLoggedInCookie, krudMW.PermissionEdit)
		g.DELETE("/delete/:schemaId/:id", handlers.Delete(GetGridMODEL), agentMW.IsLoggedInCookie, krudMW.PermissionDelete)

	} else {
		g.POST("/update-row/:schemaId", handlers.UpdateRow(GetGridMODEL), agentMW.IsLoggedInCookie)
		g.POST("/:schemaId/:action", handlers.Crud(GetMODEL), agentMW.IsLoggedInCookie)
		g.POST("/:schemaId/:action/:id", handlers.Crud(GetMODEL), agentMW.IsLoggedInCookie)
		g.DELETE("/delete/:schemaId/:id", handlers.Delete(GetGridMODEL), agentMW.IsLoggedInCookie)
	}

	/*
		OTHER
	*/
	g.POST("/upload", handlers.Upload)
	g.OPTIONS("/upload", handlers.Upload)
	//g.POST("/upload", handlers.Upload, agentMW.IsLoggedInCookie)
	//g.OPTIONS("/upload", handlers.Upload, agentMW.IsLoggedInCookie)
	g.POST("/unique", handlers.CheckUnique)
	g.POST("/check_current_password", handlers.CheckCurrentPassword, agentMW.IsLoggedInCookie)
	g.POST("/excel/:schemaId", handlers.ExportExcel(GetGridMODEL), agentMW.IsLoggedInCookie)

	/*
		PUBLIC CURDS
	*/
	public := e.Group("/lambda/krud-public")
	public.POST("/:schemaId/:action", handlers.Crud(GetMODEL))
	p := e.Group("lambda/krud-public")
	p.POST("/:schemaId/:action", handlers.Crud(GetMODEL))
}
