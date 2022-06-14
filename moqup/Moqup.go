package moqup

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/moqup/handler"
	lambdaUtils "github.com/lambda-platform/lambda/utils"
	"github.com/lambda-platform/lambda/moqup/utils"
	"html/template"
)

func Set(e *echo.Echo) {

	templates := lambdaUtils.GetTemplates(e)
	AbsolutePath := utils.AbsolutePath()

	templates["moqup.html"] = template.Must(template.ParseFiles(
		AbsolutePath + "views/moqup.html",
	))
	e.GET("/pages/moqup/:id", handler.Moqup)

}
