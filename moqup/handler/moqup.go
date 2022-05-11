package handler

import (
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Moqup(c echo.Context) error  {
	id := c.Param("id")
	//csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	csrfToken := ""
	return c.Render(http.StatusOK, "moqup.html", map[string]interface{}{
		"title":                     config.LambdaConfig.Title,
		"favicon":                   config.LambdaConfig.Favicon,
		"id":                   id,
		"csrfToken":                   csrfToken,
		"mix":                       utils.Mix,
	})
}
