package datagrid

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
)

func DeleteData(c echo.Context, datagrid Datagrid, id string) error {

	//fmt.Println(Identity, id, "Identity, id")
	qr := DB.DB.Where(datagrid.Identity+" = ?", id)
	err := qr.Delete(datagrid.MainModel).Error

	if err != nil {

		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false",
		})

	} else {

		CallTrigger("afterDelete", datagrid, []map[string]interface{}{}, id, qr, c)

		return c.JSON(http.StatusOK, map[string]string{
			"status": "true",
		})
	}
}
