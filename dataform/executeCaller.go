package dataform

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Exec(c echo.Context, schemaId string, action string, id string, GetMODEL func(schemaId string) Dataform) error {
	datform := GetMODEL(schemaId)

	switch action {
	case "store":
		return Store(c, datform, action, id)
	case "update":
		return Store(c, datform, action, id)
	case "edit":
		return Edit(c, datform, id)
	case "options":
		return Options(c)
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"status": "false",
	})

}




