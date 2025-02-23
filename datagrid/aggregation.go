package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"strconv"
)

func Aggregation(c *fiber.Ctx, datagrid Datagrid) error {

	var data []interface{}
	aggregationData := GetAggregationData(c, datagrid)
	data = append(data, aggregationData)
	return c.JSON(data)
}

// Fetch Aggregation Data
func GetAggregationData(c *fiber.Ctx, datagrid Datagrid) map[string]interface{} {
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
	query = query.Select(datagrid.Aggregation)
	rows, _ := query.Rows()

	aggregationData := make(map[string]interface{})
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		for i, col := range columns {
			val := values[i]

			// Handle byte slices directly as numbers
			if b, ok := val.([]byte); ok {
				strValue := string(b)

				// First, try parsing as an integer
				if intVal, err := strconv.ParseInt(strValue, 10, 64); err == nil {
					aggregationData[col] = intVal
					continue
				}

				// If integer parsing fails, try parsing as a float
				if floatVal, err := strconv.ParseFloat(strValue, 64); err == nil {
					aggregationData[col] = floatVal
					continue
				}

				// If parsing fails, default to 0
				aggregationData[col] = 0
				continue
			}

			// Ensure known numeric types are stored correctly
			switch v := val.(type) {
			case int, int8, int16, int32, int64:
				aggregationData[col] = v
			case float32, float64:
				aggregationData[col] = v
			default:
				// If it's a string that should be numeric, try converting
				if str, ok := v.(string); ok {
					if intVal, err := strconv.ParseInt(str, 10, 64); err == nil {
						aggregationData[col] = intVal
					} else if floatVal, err := strconv.ParseFloat(str, 64); err == nil {
						aggregationData[col] = floatVal
					} else {
						aggregationData[col] = 0 // Default to 0 if conversion fails
					}
				}
			}
		}
	}

	return aggregationData
}
