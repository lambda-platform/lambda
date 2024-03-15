package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ExecTrigger(action string, data interface{}, datagrid Datagrid, query *gorm.DB, c *fiber.Ctx) (interface{}, *gorm.DB, bool, bool) {

	switch action {
	case "afterDelete":
		if datagrid.AfterDelete != nil {
			return datagrid.AfterDelete(data, datagrid, query, c)
		}
	case "beforeFetch":
		if datagrid.BeforeFetch != nil {
			return datagrid.BeforeFetch(data, datagrid, query, c)
		}
	case "beforeDelete":
		if datagrid.BeforeDelete != nil {
			return datagrid.BeforeDelete(data, datagrid, query, c)
		}
	case "beforePrint":
		if datagrid.BeforePrint != nil {
			return datagrid.BeforePrint(data, datagrid, query, c)
		}

	}

	return []interface{}{}, query, false, false

}
