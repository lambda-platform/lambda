package datagrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/tealeg/xlsx"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

func ExportExcel(c *fiber.Ctx, datagrid Datagrid) error {
	sortColumn := c.Query("sort")
	order := c.Query("order")
	name := cleanSheetName(trim(datagrid.Name, 21))
	query := DB.DB.Table(datagrid.DataTable)

	if sortColumn != "null" && sortColumn != "undefined" {
		if order == "asc" || order == "desc" || order == "ASC" || order == "DESC" {
			query = query.Order(sortColumn + " " + order)
		}
	}

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

			rowColumns := ""

			for _, column := range datagrid.Columns {

				value := getCellValue(rows[i][column.Model], column.GridType)
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

		// Getting max widths for each column
		maxWidths := getWidths(rows, datagrid.Columns)

		for i := range rows {

			dataRow := sheet.AddRow()

			for _, column := range datagrid.Columns {

				dataCell := dataRow.AddCell()

				setCellValue(rows[i][column.Model], column.GridType, dataCell)

			}

		}

		// Apply the widths
		for idx, width := range maxWidths {

			sheet.Col(idx).Width = width
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

func setCellValue(rawValue interface{}, GridType string, cell *xlsx.Cell) {

	if rawValue != nil {

		switch v := rawValue.(type) {
		case time.Time:
			if GridType == "Datetime" {
				cell.SetDateTime(v)
			} else {
				cell.SetDate(v)
			}
		case float32:
			vFloat := rawValue.(float32)
			if vFloat == float32(int(vFloat)) {

				cell.SetInt(int(vFloat))

			} else {
				cell.SetFloat(float64(vFloat))
			}
		case float64:
			vFloat := rawValue.(float64)
			if vFloat == float64(int(vFloat)) {

				cell.SetInt(int(vFloat))
			} else {

				cell.SetFloat(vFloat)
			}

		case string:
			if t, err3 := time.Parse(time.RFC3339, v); err3 == nil {
				if GridType == "Datetime" {
					cell.SetDateTime(t)
				} else {
					cell.SetDate(t)
				}
			} else {
				cell.SetString(StripTags(rawValue.(string)))
			}
		case int:
			cell.SetInt(v)

		case int8:
			cell.SetInt(int(v))

		case int16:
			cell.SetInt(int(v))

		case int32:
			cell.SetInt(int(v))

		case int64:
			cell.SetInt(int(v))

		default:

			cell.SetString(fmt.Sprintf("%v", rawValue))
		}

	}

}
func getCellValue(rawValue interface{}, GridType string) string {

	value := ""
	if rawValue != nil {

		switch v := rawValue.(type) {
		case time.Time:
			if GridType == "Datetime" {
				value = v.Format("2006-01-02 15:04:05")
			} else {
				value = v.Format("2006-01-02")
			}
		case float32:
			vFloat := rawValue.(float32)
			if vFloat == float32(int(vFloat)) {

				value = fmt.Sprintf("%d", int(vFloat))
			} else {

				value = fmt.Sprintf("%f", vFloat)
			}
		case float64:
			vFloat := rawValue.(float64)
			if vFloat == float64(int(vFloat)) {

				value = fmt.Sprintf("%d", int(vFloat))
			} else {

				value = fmt.Sprintf("%f", vFloat)
			}

		case string:
			if t, err3 := time.Parse(time.RFC3339, v); err3 == nil {
				if GridType == "Datetime" {
					value = t.Format("2006-01-02 15:04:05")
				} else {
					value = t.Format("2006-01-02")
				}
			} else {
				value = StripTags(rawValue.(string))
			}
		case int, int8, int16, int32, int64:

			value = fmt.Sprintf("%d", rawValue)

		default:

			value = fmt.Sprintf("%v", rawValue)
		}

		if value == "<nil>" {
			value = ""
		}
	}

	return value
}

func getWidths(data []map[string]interface{}, columns []Column) []float64 {
	maxWidths := make([]float64, len(columns))
	for _, row := range data {
		for idx, col := range columns {
			cellValue := getCellValue(row[col.Model], col.GridType)

			LabelLength := float64(utf8.RuneCountInString(col.Label))
			cellLength := float64(utf8.RuneCountInString(cellValue))

			if cellLength > maxWidths[idx] {
				maxWidths[idx] = cellLength
			}

			if LabelLength > maxWidths[idx] {
				maxWidths[idx] = LabelLength
			}
		}
	}
	return maxWidths
}

func trim(s string, length int) string {
	var size, x int

	for i := 0; i < length && x < len(s); i++ {
		_, size = utf8.DecodeRuneInString(s[x:])
		x += size
	}

	return s[:x]
}

func StripTags(html string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(html, "")
}
func cleanSheetName(name string) string {
	name = strings.ReplaceAll(name, "\\", " ")
	name = strings.ReplaceAll(name, "/", " ")
	name = strings.ReplaceAll(name, "?", " ")
	name = strings.ReplaceAll(name, "*", " ")
	name = strings.ReplaceAll(name, "[", " ")
	name = strings.ReplaceAll(name, "]", " ")
	return name
}
