package datagrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/tealeg/xlsx/v3"
	"io"
	"net/http"
	"os"
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

	if len(datagrid.Relations) >= 1 {
		datagrid.Data = datagrid.FillVirtualColumns(datagrid.Data)
	}

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
		headerRow.SetHeight(30)
		for _, k := range keys {

			headerCell := headerRow.AddCell()
			headerCell.Value = datagrid.Columns[k].Label
			headerCell.GetStyle().Font.Bold = true                                      // Make header bold
			headerCell.GetStyle().Font.Bold = true                                      // Make header bold
			headerCell.GetStyle().Fill = *xlsx.NewFill("solid", "0099CCFF", "0099CCFF") // Set background color
			headerCell.GetStyle().Alignment.Horizontal = "center"
			headerCell.GetStyle().Alignment.Vertical = "center"

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
		// Dynamic aggregation with uppercase SUM, AVG, COUNT, MIN, MAX formulas using SUBTOTAL
		if datagrid.Aggregation != "" {
			// Add a blank row to separate data from aggregation
			sheet.AddRow()

			// Add aggregation row with a "Totals" label
			aggRow := sheet.AddRow()
			aggRow.AddCell().SetString("") // Label in first column
			rowCount := len(rows)          // Number of data rows (excluding header)

			// Split the aggregation string into individual parts (e.g., "SUM(qty) as SUM_qty")
			aggParts := strings.Split(datagrid.Aggregation, ",")

			// Create a map of aggregation functions to column names
			aggMap := make(map[string]string) // e.g., "qty" -> "SUM", "total_price" -> "AVG"
			for _, part := range aggParts {
				part = strings.TrimSpace(part)
				// Match patterns like "SUM(qty)", "AVG(total_price)", etc.
				re := regexp.MustCompile(`(SUM|AVG|COUNT|MIN|MAX)\(([a-zA-Z_]+)\)`)
				matches := re.FindStringSubmatch(part)
				if len(matches) == 3 {
					funcName := matches[1] // e.g., "SUM", "AVG", "COUNT", "MIN", "MAX"
					colName := matches[2]  // e.g., "qty", "total_price"
					aggMap[colName] = funcName
				}
			}

			// Apply formulas based on column.Model, starting from second column
			for colIdx, column := range datagrid.Columns[1:] { // Skip first column (label)
				aggCell := aggRow.AddCell()

				if funcName, exists := aggMap[column.Model]; exists {
					// Calculate the range for the column (e.g., B2:B5)
					colLetter := string('A' + colIdx + 1) // Offset by 1 for label column
					formulaRange := fmt.Sprintf("%s2:%s%d", colLetter, colLetter, rowCount+1)

					// Set the appropriate Excel SUBTOTAL formula for filtered data
					var formula string
					switch funcName {
					case "SUM":
						formula = fmt.Sprintf("SUBTOTAL(9,%s)", formulaRange) // 9 = SUM
					case "AVG":
						formula = fmt.Sprintf("SUBTOTAL(1,%s)", formulaRange) // 1 = AVERAGE
					case "COUNT":
						formula = fmt.Sprintf("SUBTOTAL(2,%s)", formulaRange) // 2 = COUNT
					case "MIN":
						formula = fmt.Sprintf("SUBTOTAL(5,%s)", formulaRange) // 5 = MIN
					case "MAX":
						formula = fmt.Sprintf("SUBTOTAL(4,%s)", formulaRange) // 4 = MAX
					}

					if formula != "" {
						aggCell.SetFormula(formula)
						aggCell.GetStyle().Font.Bold = true
					}
				}
			}
		}

		// Apply the widths
		for idx, width := range maxWidths {
			sheet.SetColWidth(idx+1, idx+1, width)
		}
		// Add AutoFilter to the header row
		lastRow := len(rows) + 1 // +1 for header row

		colCount := len(datagrid.Columns)
		if colCount > 0 && lastRow > 1 {
			startCell := "A1"
			endColLetter := string('A' + colCount - 1)
			endCell := fmt.Sprintf("%s%d", endColLetter, lastRow)
			sheet.AutoFilter = &xlsx.AutoFilter{
				TopLeftCell:     startCell,
				BottomRightCell: endCell,
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
				if GridType == "Image" {
					imageURL := fmt.Sprintf("%s%s", config.LambdaConfig.Domain, StripTags(v))
					addImageToCell(imageURL, cell)
				} else {
					cell.SetString(StripTags(v))
				}
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
func addImageToCell(imageURL string, cell *xlsx.Cell) {
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		cell.SetString("Image Not Found")
		return
	}
	defer resp.Body.Close()

	// Read image data
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading image data:", err)
		cell.SetString("Image Error")
		return
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "image_*.png")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		cell.SetString("File Error")
		return
	}
	defer os.Remove(tmpFile.Name())

	// Write image to temp file
	_, err = tmpFile.Write(imgData)
	if err != nil {
		fmt.Println("Error writing image file:", err)
		cell.SetString("Write Error")
		return
	}
	tmpFile.Close()

	if err != nil {
		fmt.Println("Error adding image to Excel:", err)
		cell.SetString("Insert Error")
		return
	}

	// Set cell value to empty since image is placed
	cell.Value = ""
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

			LabelLength := float64(utf8.RuneCountInString(col.Label)) + 8
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
