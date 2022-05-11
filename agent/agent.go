package agent

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/agent/handlers"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	vpUtils "github.com/lambda-platform/lambda/utils"
	"html/template"
)

func Set(e *echo.Echo) {

	if config.Config.App.Migrate == "true" {
		agentUtils.AutoMigrateSeed()
	}
	templates := vpUtils.GetTemplates(e)

	/* REGISTER VIEWS */
	AbsolutePath := agentUtils.AbsolutePath()

	templates["login.html"] = template.Must(template.ParseFiles(AbsolutePath + "templates/login.html"))
	templates["forgot.html"] = template.Must(template.ParseFiles(AbsolutePath + "templates/email/forgot.html"))

	/* ROUTES */
	a := e.Group("/auth")
	a.GET("/", handlers.LoginPage)
	a.GET("/login", handlers.LoginPage)
	a.POST("/login", handlers.Login)
	a.POST("/logout", handlers.Logout)

	/*PASSWORD RESET*/
	a.POST("/send-forgot-mail", handlers.SendForgotMail)
	a.POST("/password-reset", handlers.PasswordReset)

	u := e.Group("/agent")
	u.GET("/users", handlers.GetUsers, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	u.GET("/search/:q", handlers.SearchUsers, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	u.GET("/users/deleted", handlers.GetDeletedUsers, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	u.GET("/delete/:id", handlers.DeleteUser, agentMW.IsLoggedInCookie, agentMW.IsAdmin)
	u.GET("/roles", handlers.GetRoles, agentMW.IsLoggedInCookie, agentMW.IsAdmin)

}
