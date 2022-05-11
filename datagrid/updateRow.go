package datagrid

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
)

func UpdateRow(c echo.Context, datagrid Datagrid) error {

	RowUpdateData := new(RowUpdateData)

	if err := c.Bind(RowUpdateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
			"error":  err.Error(),
		})
	}
	if len(RowUpdateData.Ids) >= 1 && RowUpdateData.Model != "" && RowUpdateData.Value != nil {
		for _, id_ := range RowUpdateData.Ids {

			DB.DB.Model(datagrid.MainModel).Where(datagrid.Identity+" = ?", id_).Update(RowUpdateData.Model, RowUpdateData.Value)
		}

	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "true",
	})
}
