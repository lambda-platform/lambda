package datagrid

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
	"strconv"
)

func Aggregation(c echo.Context, datagrid Datagrid) error {

	//GetAggregations := reflect.ValueOf(GridModel).MethodByName("GetAggregations")
	//aggregationsRes := GetAggregations.Call([]reflect.Value{})
	//aggregations := aggregationsRes[0].Interface().([]map[string]string)
	//
	//
	//func (v *DSIrtsiinBurtgel524) GetAggregations() []map[string]string {
	//	//[{"column":"tetgeleg_dun","aggregation":"SUM","symbol":"₮"},{"column":"id","aggregation":"COUNT","symbol":"Нийт "}]
	//
	//	aggregations := []map[string]string{
	//	map[string]string{
	//	"column": "tetgeleg_dun",
	//	"aggregation": "SUM",
	//	"symbol": "₮",
	//},
	//	map[string]string{
	//	"column": "id",
	//	"aggregation": "COUNT",
	//	"symbol": "Нийт ",
	//},
	//}
	//	return aggregations
	//}



	query := DB.DB.Table(datagrid.DataTable)

	query, _ = Filter(c, datagrid, query)

	if len(datagrid.Condition) > 0 {
		query = query.Where(datagrid.Condition)
	}

	query = Search(c, datagrid.DataModel, query)



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

	return c.JSON(http.StatusOK, data)
}
