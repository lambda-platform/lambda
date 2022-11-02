package agent

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/agent/handlers"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
)

func Set(e *fiber.App) {
	fmt.Println(config.Config.App.Migrate)
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

	e.Get("/get-permissions", agentMW.IsLoggedIn(), handlers.GetPermissions)
	a.Post("/logout", handlers.Logout)

	/*PASSWORD RESET*/
	a.Post("/send-forgot-mail", handlers.SendForgotMail)
	a.Post("/password-reset", handlers.PasswordReset)

	u := e.Group("/agent")
	u.Get("/users", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetUsers)
	u.Get("/search/:q", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.SearchUsers)
	u.Get("/users/deleted", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetDeletedUsers)
	u.Get("/delete/:id", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.DeleteUser)
	u.Get("/roles", agentMW.IsLoggedIn(), agentMW.IsAdmin, handlers.GetRoles)

}
