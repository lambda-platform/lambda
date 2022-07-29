package datagrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/tealeg/xlsx"
	"net/http"
	"reflect"
	"sort"
	"unicode/utf8"
)

func ExportExcel(c *fiber.Ctx, datagrid Datagrid) error {

	name := trim(datagrid.Name, 21)

	query := DB.DB.Table(datagrid.DataTable)
	customHeader := ""
	query, customHeader = Filter(c, datagrid, query)

	if len(datagrid.Condition) > 0 {
		query = query.Where(datagrid.Condition)
	}

	_, query, _, _ = ExecTrigger("beforeFetch", []interface{}{}, datagrid, query, c)

	query.Find(datagrid.Data)

	keys := make([]int, 0, len(datagrid.Columns))
	for k := range datagrid.Columns {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	if customHeader != "" {
		data := ""
		/*HEADER*/

		rowTemplate := `
<tr>%s</tr>
`
		colTemplate := `
<td>%s</td>
`
		headerRow := ""

		if customHeader == "" {

			headerColumns := ""
			for _, column := range datagrid.Columns {

				headerColumns = headerColumns + fmt.Sprintf(colTemplate, column.Label)
			}

			headerRow = fmt.Sprintf(rowTemplate, headerColumns)
		}

		rows_json, _ := json.Marshal(datagrid.Data)

		var rows []map[string]interface{}
		json.Unmarshal(rows_json, &rows)

		data = headerRow
		for i := range rows {

			row_ := rows[i]

			rowColumns := ""

			for _, column := range datagrid.Columns {

				value := getCellValue(row_[column.Model])

				rowColumns = rowColumns + fmt.Sprintf(colTemplate, value)
			}
			data = data + fmt.Sprintf(rowTemplate, rowColumns)

		}

		return c.JSON(map[string]interface{}{
			"name":      name,
			"tableRows": data,
		})
	} else {

		var file *xlsx.File
		var sheet *xlsx.Sheet

		var err error

		file = xlsx.NewFile()
		sheet, err = file.AddSheet(name)
		if err != nil {
			fmt.Printf(err.Error())
		}

		headerRow := sheet.AddRow()
		for _, k := range keys {

			headerCell := headerRow.AddCell()
			headerCell.Value = datagrid.Columns[k].Label
		}
		/*HEADER*/

		rows_json, _ := json.Marshal(datagrid.Data)

		var rows []map[string]interface{}
		json.Unmarshal(rows_json, &rows)

		for i := range rows {

			row_ := rows[i]

			dataRow := sheet.AddRow()

			for _, column := range datagrid.Columns {

				dataCell := dataRow.AddCell()

				dataCell.Value = getCellValue(row_[column.Model])
			}

		}

		var b bytes.Buffer
		if err := file.Write(&b); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{
				"status": "false",
			})
		}

		return c.JSON(map[string]interface{}{
			"name": name + ".xlsx",
			"file": b.Bytes(),
		})
	}

}
func getCellValue(rawValue interface{}) string {
	value := ""
	if reflect.TypeOf(rawValue).String() == "float64" {

		value = fmt.Sprintf("%.3f", rawValue)

	} else if reflect.TypeOf(rawValue).String() == "float32" {

		value = fmt.Sprintf("%.3f", rawValue)

	} else if reflect.TypeOf(rawValue).String() == "string" {

		value = reflect.ValueOf(rawValue).Interface().(string)

	} else if reflect.TypeOf(rawValue).String() == "int" {

		value = reflect.ValueOf(rawValue).Interface().(string)

	} else if reflect.TypeOf(rawValue).String() == "Int" {

		value = fmt.Sprintf("%d", rawValue)

	} else if reflect.TypeOf(rawValue).String() == "Int" {

		value = fmt.Sprintf("%d", rawValue)

	} else if reflect.TypeOf(rawValue).String() == "Int8" {

		value = fmt.Sprintf("%d", rawValue)

	} else if reflect.TypeOf(rawValue).String() == "Int16" {

		value = fmt.Sprintf("%d", rawValue)

	} else if reflect.TypeOf(rawValue).String() == "Int32" {

		value = fmt.Sprintf("%d", rawValue)

	} else if reflect.TypeOf(rawValue).String() == "Int64" {

		value = fmt.Sprintf("%d", rawValue)

	} else {
		value = fmt.Sprintf("%v", rawValue)
	}

	if value == "<nil>" {
		value = ""
	}
	return value
}
func trim(s string, length int) string {
	var size, x int

	for i := 0; i < length && x < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[x:])
		x += size
	}

	return s[:x]
}
