package datagrid

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/dataform"
)

func FilterOptions(c *fiber.Ctx, datagrid Datagrid) error {

	var optionsData = map[string][]map[string]interface{}{}
	fmt.Println(datagrid.FilterRelations)
	for table, relation := range datagrid.FilterRelations {
		data := dataform.OptionsData(relation, c)

		optionsData[table] = data

	}

	return c.JSON(optionsData)
}
