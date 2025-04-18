package lambda

import (
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/generator"
	"github.com/lambda-platform/lambda/middlewares"
	"github.com/lambda-platform/lambda/puzzle/views"
	"os"
)

var viewsfs embed.FS

type Lambda struct {
	App          *fiber.App
	ModuleName   string
	IgnoreStatic bool
}

func (lambda *Lambda) Start() {
	if len(os.Args) < 2 {

		lambda.App.Listen(":" + config.Config.App.Port)

	} else {
		command := os.Args[1]

		switch command {

		case "table":
			if len(os.Args) < 3 {
				fmt.Println("Please provide table name: table your-table-name")
			} else {
				table := os.Args[2]
				generator.GetStruct(table)
			}
		case "proto":
			if len(os.Args) < 3 {
				fmt.Println("Please provide table name: table your-table-name")
			} else {
				table := os.Args[2]
				generator.GetProtobuf(table)
			}
		default:
			fmt.Printf("Unknown command: %s\n", command)
			os.Exit(1)
		}
	}

	//defer DB.DB.Close()
}

type Settings struct {
	ModuleName   string
	IgnoreStatic bool
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
			Views:          engine,
			BodyLimit:      100 * 1024 * 1024, // this is the default limit of 100MB
			ReadBufferSize: 1024 * 1024,       // Set this to the desired size in  1MB
			//JSONEncoder: json.Marshal,
			//JSONDecoder: json.Unmarshal,
		}),
		ModuleName:   lambdaSettings[0].ModuleName,
		IgnoreStatic: lambdaSettings[0].IgnoreStatic,
	}
	middlewares.Set(lambda.App)
	if !lambdaSettings[0].IgnoreStatic {
		lambda.App.Static("/", "public")
	}

	return lambda
}
