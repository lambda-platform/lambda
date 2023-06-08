package puzzle

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/puzzle/handlers"
	"github.com/lambda-platform/lambda/puzzle/utils"
)

func Set(e *fiber.App, moduleName string, GetGridMODEL func(schema_id string) datagrid.Datagrid, isMicroservice bool, withUserRole bool) {

	if isMicroservice {

	} else {
		if config.Config.App.Migrate == "true" {
			utils.AutoMigrateSeed()
		}
	}
	//if isMicroservice && withUserRole{
	//	handlers.GetRoleData()
	//}
	//if withUserRole || !isMicroservice {
	//	templates := lambdaUtils.GetTemplates(e)
	//
	//	//* REGISTER VIEWS */
	//	AbsolutePath := utils.AbsolutePath()
	//	TemplatePath := templateUtils.AbsolutePath()
	//
	//	templates["puzzle.html"] = template.Must(template.ParseFiles(
	//		TemplatePath + "views/paper.html",
	//	))
	//
	//	template.Must(templates["puzzle.html"].ParseFiles(
	//		AbsolutePath + "views/puzzle.html",
	//	))
	//}

	/*ROUTES */
	e.Get("/build-me", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.BuildMe)

	g := e.Group("/lambda")

	//g.Get("/puzzle", handlers.Index, agentMW.IsLoggedIn())
	g.Get("/puzzle", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.Index)

	//Puzzle
	g.Get("/puzzle/schema/:type", agentMW.IsLoggedIn(), handlers.GetVB)
	g.Get("/puzzle/schema/:type/:id", agentMW.IsLoggedIn(), handlers.GetVB)
	g.Get("/puzzle/schema-public/:type/:id", handlers.GetVB)
	g.Get("/puzzle/schema/:type/:id/:condition", agentMW.IsLoggedIn(), handlers.GetVB)

	//VB SCHEMA
	g.Get("/puzzle/table-schema/:table", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetTableSchema)
	g.Post("/puzzle/schema/:type", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.SaveVB(moduleName))
	g.Post("/puzzle/schema/:type/:id", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.SaveVB(moduleName))
	g.Delete("/puzzle/delete/vb_schemas/:type/:id", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.DeleteVB)
	//MENU SHOW
	e.Get("/lambda/krud/menu_form/edit/:id", agentMW.IsLoggedIn(), handlers.GetMenuVB)

	//GRID
	g.Post("/puzzle/grid/:action/:schemaId", agentMW.IsLoggedIn(), handlers.GridVB(GetGridMODEL))
	g.Post("/puzzle/grid-public/:action/:schemaId", handlers.GridVB(GetGridMODEL))

	//Get From Options
	g.Post("/puzzle/get_options", agentMW.IsLoggedIn(), handlers.GetOptions)
	g.Post("/puzzle/get_options-public", handlers.GetOptions)

	//Roles
	g.Get("/puzzle/roles-menus", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetRolesMenus)
	g.Get("/puzzle/roles-menus/:microserviceID", agentMW.IsLoggedIn(), agentMW.IsCloudUser, handlers.GetRolesMenus)
	g.Get("/puzzle/get-krud-fields/:id", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetKrudFields)
	g.Get("/puzzle/get-krud-fields-micro/:id", agentMW.IsLoggedIn(), handlers.GetKrudFieldsConsole)
	if isMicroservice {
		g.Post("/puzzle/save-role", agentMW.IsLoggedIn(), agentMW.IsCloudUser, handlers.SaveRole)
	} else {
		g.Post("/puzzle/save-role", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.SaveRole)
	}

	g.Post("/puzzle/roles/create", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.CreateRole)
	g.Post("/puzzle/roles/store/:id", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.UpdateRole)
	g.Delete("/puzzle/roles/destroy/:id", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.DeleteRole)

	//Puzzle with project
	g.Get("/puzzle/projects/:type", agentMW.IsLoggedIn(), handlers.GetProjectVBs)
	g.Get("/puzzle/projects/:type/:id", agentMW.IsLoggedIn(), handlers.GetProjectVBs)
	g.Get("/puzzle/project/:pid/:type", agentMW.IsLoggedIn(), handlers.GetProjectVB)
	g.Get("/puzzle/project/:pid/:type/:id", agentMW.IsLoggedIn(), handlers.GetProjectVB)
	g.Get("/puzzle/project/:pid/:type/:id/builder", agentMW.IsLoggedIn(), handlers.GetProjectVB)
	g.Post("/puzzle/project/:pid/:type", agentMW.IsLoggedIn(), handlers.SaveProjectVB(moduleName))
	g.Post("/puzzle/project/:pid/:type/:id", agentMW.IsLoggedIn(), handlers.SaveProjectVB(moduleName))
	g.Delete("/puzzle/delete/project/vb_schemas/:pid/:type/:id", agentMW.IsLoggedIn(), handlers.DeleteProjectVB)

}
