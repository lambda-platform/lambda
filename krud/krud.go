package krud

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/krud/handlers"
	"github.com/lambda-platform/lambda/krud/krudMW"
	"github.com/lambda-platform/lambda/krud/utils"
)

func Set(e *fiber.App, GetGridMODEL func(schema_id string) datagrid.Datagrid, GetMODEL func(schema_id string) dataform.Dataform, krudMiddleWares []fiber.Handler, KrudWithPermission bool) {
	if config.Config.App.Migrate == "true" {
		utils.AutoMigrateSeed()
	}

	g := e.Group("/lambda/krud")
	if len(krudMiddleWares) >= 1 {
		for _, krudMiddleWare := range krudMiddleWares {
			g.Use(krudMiddleWare)
		}
	}
	kp := e.Group("/lambda/krud-public")
	if len(krudMiddleWares) >= 1 {
		for _, krudMiddleWare := range krudMiddleWares {
			kp.Use(krudMiddleWare)
		}
	}
	p := e.Group("/lambda/puzzle")
	if len(krudMiddleWares) >= 1 {
		for _, krudMiddleWare := range krudMiddleWares {
			p.Use(krudMiddleWare)
		}
	}

	g.Post("/excel/:schemaId", agentMW.IsLoggedIn(), handlers.ExportExcel(GetGridMODEL))
	g.Post("/import-excel/:schemaId", agentMW.IsLoggedIn(), handlers.ImportExcel(GetGridMODEL))
	g.Post("/print/:schemaId", agentMW.IsLoggedIn(), handlers.Print(GetGridMODEL))
	if KrudWithPermission {
		g.Post("/:schemaId/filter-options", agentMW.IsLoggedIn(), krudMW.PermissionDelete, handlers.FilterOptions(GetGridMODEL))
		g.Post("/update-row/:schemaId", agentMW.IsLoggedIn(), krudMW.PermissionDelete, handlers.UpdateRow(GetGridMODEL))
		g.Post("/:schemaId/:action", agentMW.IsLoggedIn(), krudMW.PermissionCreate, handlers.Crud(GetMODEL))
		g.Post("/:schemaId/:action/:id", agentMW.IsLoggedIn(), krudMW.PermissionEdit, handlers.Crud(GetMODEL))
		g.Delete("/delete/:schemaId/:id", agentMW.IsLoggedIn(), krudMW.PermissionDelete, handlers.Delete(GetGridMODEL))

	} else {
		g.Post("/:schemaId/filter-options", agentMW.IsLoggedIn(), handlers.FilterOptions(GetGridMODEL))
		g.Post("/update-row/:schemaId", agentMW.IsLoggedIn(), handlers.UpdateRow(GetGridMODEL))
		g.Post("/:schemaId/:action", agentMW.IsLoggedIn(), handlers.Crud(GetMODEL))
		g.Post("/:schemaId/:action/:id", agentMW.IsLoggedIn(), handlers.Crud(GetMODEL))
		g.Delete("/delete/:schemaId/:id", agentMW.IsLoggedIn(), handlers.Delete(GetGridMODEL))

	}
	kp.Post("/:schemaId/:action", krudMW.PermissionCreate, handlers.Crud(GetMODEL))
	/*
		OTHER
	*/
	g.Post("/upload", handlers.Upload)
	g.Options("/upload", handlers.Upload)
	//g.Post("/upload", handlers.Upload, agentMW.IsLoggedIn())
	//g.OPTIONS("/upload", handlers.Upload, agentMW.IsLoggedIn())
	g.Post("/unique", handlers.CheckUnique)
	g.Get("/today", handlers.Now)
	g.Get("/now", handlers.Now)
	g.Post("/check_current_password", agentMW.IsLoggedIn(), handlers.CheckCurrentPassword)

	/*
		PUBLIC CURDS
	*/
	public := e.Group("/lambda/krud-public")
	public.Post("/:schemaId/:action", handlers.Crud(GetMODEL))
	public.Post("/:schemaId/:action/:id", handlers.Crud(GetMODEL))

}
