package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
)

func Print(c *fiber.Ctx, datagrid Datagrid) error {

	query := DB.DB.Table(datagrid.DataTable)

	//DB.DB.LogMode(true)
	query, _ = Filter(c, datagrid, query)

	if len(datagrid.Condition) > 0 {
		query = query.Where(datagrid.Condition)
	}
	query = Search(c, datagrid.DataModel, query)

	_, query, _, _ = ExecTrigger("beforePrint", []interface{}{}, datagrid, query, c)

	query.Find(datagrid.Data)

	return c.JSON(datagrid.Data)

}
