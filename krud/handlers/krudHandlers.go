package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
)

func Crud(GetMODEL func(schema_id string) dataform.Dataform) echo.HandlerFunc {
	return func(c echo.Context) error {
		schemaId := c.Param("schemaId")
		action := c.Param("action")
		id := c.Param("id")

		return dataform.Exec(c, schemaId, action, id, GetMODEL)
	}
}

func CheckUnique(c echo.Context) error {
	return dataform.CheckUnique(c)
}
func Upload(c echo.Context) error {

	return dataform.Upload(c)

}
func CheckCurrentPassword(c echo.Context) error {
	return utils.CheckCurrentPassword(c)
}

func UpdateRow(GetGridMODEL func(schema_id string) datagrid.Datagrid) echo.HandlerFunc {
	return func(c echo.Context) error {
		schemaId := c.Param("schemaId")

		return datagrid.Exec(c, schemaId, "update-row", "", GetGridMODEL)
	}
}

func Delete(GetGridMODEL func(schema_id string) datagrid.Datagrid) echo.HandlerFunc {
	return func(c echo.Context) error {
		schemaId := c.Param("schemaId")
		id := c.Param("id")

		return datagrid.Exec(c, schemaId, "delete", id, GetGridMODEL)
	}
}
func ExportExcel(GetGridMODEL func(schema_id string) datagrid.Datagrid) echo.HandlerFunc {
	return func(c echo.Context) error {
		schemaId := c.Param("schemaId")

		return datagrid.Exec(c, schemaId, "excel", "", GetGridMODEL)
	}
}

func dieIF(err error) {
	if err != nil {
		panic(err)
	}
}
