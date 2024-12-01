package agent

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/agent/handlers"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
)

func Set(e *fiber.App) {

	if config.Config.App.Migrate == "true" {
		agentUtils.AutoMigrateSeed()
	}

	//e.Use(jwtware.New(jwtware.Config{
	//	KeyFunc: agentMW.KeyFunc(),
	//}))

	e.Get("/lambda-config", handlers.LambdaConfig)
	/* ROUTES */
	a := e.Group("/auth")
	//a.Get("/", handlers.LoginPage)
	a.Get("/login", handlers.LoginPage)
	a.Get("/forgot", handlers.LoginPage)
	a.Post("/login", handlers.Login)
	a.Get("/check", agentMW.IsLoggedIn(), handlers.CheckAuth)

	e.Get("/get-permissions", agentMW.IsLoggedIn(), handlers.GetPermissions)
	a.Post("/logout", handlers.Logout)

	/*PASSWORD RESET*/
	a.Post("/send-forgot-mail", handlers.SendForgotMail)
	a.Post("/password-reset", handlers.PasswordReset)

	u := e.Group("/agent")
	if config.LambdaConfig.ProjectKey != "" && config.LambdaConfig.LambdaMainServicePath != "" {
		u.Get("/users", agentMW.IsLoggedIn(), agentMW.IsCloudUser, handlers.GetUsers)
		u.Get("/search/:q", agentMW.IsLoggedIn(), agentMW.IsCloudUser, handlers.SearchUsers)
		u.Get("/users/deleted", agentMW.IsLoggedIn(), agentMW.IsCloudUser, handlers.GetDeletedUsers)
		u.Get("/delete/:id", agentMW.IsLoggedIn(), agentMW.IsCloudUser, handlers.DeleteUser)
		u.Get("/roles", agentMW.IsLoggedIn(), agentMW.IsCloudUser, handlers.GetRoles)
	} else {
		u.Get("/users", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetUsers)
		u.Get("/search/:q", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.SearchUsers)
		u.Get("/users/deleted", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetDeletedUsers)
		u.Get("/delete/:id", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.DeleteUser)
		u.Get("/roles", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetRoles)
	}

}
