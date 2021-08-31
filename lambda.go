package lambda


import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/agent"
	"github.com/lambda-platform/krud"
	"github.com/labstack/echo/v4"
)

type Lambda struct {
	Echo         *echo.Echo
	ModuleName   string
	GetGridMODEL func(schemaId string) (interface{}, interface{}, string, string, interface{}, string)
	GetFormMODEL     func(schemaId string) (string, interface{})
	GetFormMessages  func(schemaId string) map[string][]string
	GetFormRules     func(schemaId string) map[string][]string
	echoWrapHandler     echo.HandlerFunc
	KrudMiddleWares []echo.MiddlewareFunc
	KrudWithPermission bool

}

func (app *Lambda) Start() {
	app.Echo.Logger.Fatal(app.Echo.Start(":"+config.Config.App.Port))
	defer DB.DB.Close()
}
type Settings struct {
	ModuleName string
	GetGridMODEL func(schemaId string) (interface{}, interface{}, string, string, interface{}, string)
	GetFormMODEL func(schemaId string) (string, interface{})
	GetFormMessages func(schemaId string) map[string][]string
	GetFormRules func(schemaId string) map[string][]string
	KrudMiddleWares []echo.MiddlewareFunc
	KrudWithPermission bool
}


func New(lambdaSettings ...*Settings) *Lambda {

	if len(lambdaSettings) == 0 {
		panic("Lambda settings required")
	}
	lambda := &Lambda{
		Echo:         echo.New(),
		ModuleName:   lambdaSettings[0].ModuleName,
		GetGridMODEL: lambdaSettings[0].GetGridMODEL,
		GetFormMODEL: lambdaSettings[0].GetFormMODEL,
		GetFormMessages: lambdaSettings[0].GetFormMessages,
		GetFormRules: lambdaSettings[0].GetFormRules,
		KrudWithPermission: lambdaSettings[0].KrudWithPermission,

	}
	if(len(lambdaSettings[0].KrudMiddleWares) >= 1){
		lambda.KrudMiddleWares = lambdaSettings[0].KrudMiddleWares
	}

	agent.Set(lambda.Echo)
	krud.Set(lambda.Echo, lambda.GetGridMODEL, lambda.GetFormMODEL, lambda.GetFormMessages, lambda.GetFormRules, lambda.KrudMiddleWares, lambda.KrudWithPermission)



	lambda.Echo.Static("/", "public")

	return lambda
}
