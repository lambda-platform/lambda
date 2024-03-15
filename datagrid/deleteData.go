package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
)

func DeleteData(c *fiber.Ctx, datagrid Datagrid, id string) error {

	//fmt.Println(Identity, id, "Identity, id")
	qr := DB.DB.Where(datagrid.Identity+" = ?", id)
	err := qr.Delete(datagrid.MainModel).Error

	ExecTrigger("beforeDelete", id, datagrid, qr, c)

	if err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})

	} else {

		ExecTrigger("afterDelete", id, datagrid, qr, c)

		return c.JSON(map[string]interface{}{
			"status": true,
		})
	}
}
