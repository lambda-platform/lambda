package datagrid

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Exec(c echo.Context, schemaId string, action string, id string, GetGridMODEL func(schema_id string) Datagrid) error {

	datagrid := GetGridMODEL(schemaId)

	switch action {
	case "data":
		return FetchData(c, datagrid)
	case "aggergation":
		return Aggregation(c, datagrid)
	case "delete":
		return DeleteData(c, datagrid, id)
	case "excel":
		return ExportExcel(c, datagrid)
	case "update-row":
		return UpdateRow(c, datagrid)
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"status": "false",
	})
}






