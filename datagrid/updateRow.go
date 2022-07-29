package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
)

func UpdateRow(c *fiber.Ctx, datagrid Datagrid) error {

	RowUpdateData := new(RowUpdateData)

	if err := c.BodyParser(RowUpdateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status": "false",
			"error":  err.Error(),
		})
	}
	if len(RowUpdateData.Ids) >= 1 && RowUpdateData.Model != "" && RowUpdateData.Value != nil {
		for _, id_ := range RowUpdateData.Ids {

			DB.DB.Model(datagrid.MainModel).Where(datagrid.Identity+" = ?", id_).Update(RowUpdateData.Model, RowUpdateData.Value)
		}

	}
	return c.JSON(map[string]string{
		"status": "true",
	})
}
