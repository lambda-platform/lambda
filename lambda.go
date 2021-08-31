package lambda

import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/labstack/echo/v4"
)

type Lambda struct {
	Echo       *echo.Echo
	ModuleName string
}

func (app *Lambda) Start() {
	app.Echo.Logger.Fatal(app.Echo.Start(":" + config.Config.App.Port))
	defer DB.DB.Close()
}

type Settings struct {
	ModuleName string
}

func New(lambdaSettings ...*Settings) *Lambda {

	if len(lambdaSettings) == 0 {
		panic("Lambda settings required")
	}
	lambda := &Lambda{
		Echo:       echo.New(),
		ModuleName: lambdaSettings[0].ModuleName,
	}

	lambda.Echo.Static("/", "public")

	return lambda
}
