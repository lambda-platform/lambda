package moqup

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/moqup/handler"
	lambdaUtils "github.com/lambda-platform/lambda/utils"
	"html/template"
)

func Set(e *echo.Echo) {

	templates := lambdaUtils.GetTemplates(e)

	templates["moqup.html"] = template.Must(template.ParseFiles(
		"views/moqup.html",
	))
	e.GET("/pages/moqup/:id", handler.Moqup)

}
