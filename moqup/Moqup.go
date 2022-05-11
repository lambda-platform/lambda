package moqup

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/moqup/handler"
	templateUtils "github.com/lambda-platform/lambda/moqup/utils"
	lambdaUtils "github.com/lambda-platform/lambda/utils"
	"html/template"
)

func Set(e *echo.Echo) {

	templates := lambdaUtils.GetTemplates(e)

	TemplatePath := templateUtils.AbsolutePath()

	templates["moqup.html"] = template.Must(template.ParseFiles(
		TemplatePath + "templates/moqup.html",
	))
	e.GET("/pages/moqup/:id", handler.Moqup)

}
