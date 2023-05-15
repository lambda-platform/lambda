package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"strconv"
)

func Aggregation(c *fiber.Ctx, datagrid Datagrid) error {

	query := DB.DB.Table(datagrid.DataTable)

	query, _ = Filter(c, datagrid, query)

	if len(datagrid.Condition) > 0 {
		query = query.Where(datagrid.Condition)
	}

	query = Search(c, datagrid.DataModel, query)

	// Check if the "DeletedAt" column exists
	exists := DB.DB.Migrator().HasColumn(datagrid.DataModel, "DeletedAt")
	if exists {
		if config.Config.Database.Connection == "oracle" {
			query = query.Where("DELETED_AT IS NULL")
		} else {
			query = query.Where("deleted_at IS NULL")
		}
	}
	query = query.Select(datagrid.Aggergation)
	rows, _ := query.Rows()

	data := []interface{}{}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	/*end*/

	for rows.Next() {

		/*start */

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		var myMap = make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			b, ok := val.([]byte)

			if ok {

				v, error := strconv.ParseInt(string(b), 10, 64)
				if error != nil {
					stringValue := string(b)

					myMap[col] = stringValue
				} else {
					myMap[col] = v
				}

			} else {
				myMap[col] = val
			}

		}
		/*end*/

		data = append(data, myMap)

	}

	return c.JSON(data)
}
