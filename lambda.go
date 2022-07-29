package lambda

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/lambda-platform/lambda/config"
)

type Lambda struct {
	App        *fiber.App
	ModuleName string
}

func (lambda *Lambda) Start() {
	lambda.App.Listen(":" + config.Config.App.Port)
	//defer DB.DB.Close()
}

type Settings struct {
	ModuleName string
}

func New(lambdaSettings ...*Settings) *Lambda {

	if len(lambdaSettings) == 0 {
		panic("Lambda settings required")
	}

	engine := html.New("./views", ".html")

	lambda := &Lambda{
		App: fiber.New(fiber.Config{
			Views: engine,
		}),
		ModuleName: lambdaSettings[0].ModuleName,
	}

	lambda.App.Static("/", "public")

	return lambda
}
