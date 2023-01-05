package lambda

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/puzzle/views"
)

var viewsfs embed.FS

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

	engine.Reload(false)

	err := engine.Load()

	if err != nil {
		panic(err)
	}
	_, err = engine.Templates.New("puzzle").Parse(views.PuzzleTemplate)
	if err != nil {
		panic(err)
	}

	lambda := &Lambda{
		App: fiber.New(fiber.Config{
			Views:     engine,
			BodyLimit: 100 * 1024 * 1024, // this is the default limit of 100MB
			//JSONEncoder: json.Marshal,
			//JSONDecoder: json.Unmarshal,
		}),
		ModuleName: lambdaSettings[0].ModuleName,
	}
	lambda.App.Static("/", "public")

	return lambda
}
