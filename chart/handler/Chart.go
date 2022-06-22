package handler

import (
	"fmt"
	echo "github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
	"strconv"
	"unicode"
)

func CountData(c echo.Context) (err error) {
	request := new(CountRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	if len(request.CountFields) == 1 {
		var count int64
		DB.DB.Table(request.CountFields[0].Table).Count(&count)
		return c.JSON(http.StatusOK, count)
	} else {
		return c.JSON(http.StatusOK, "0")
	}
}

func PieData(c echo.Context) (err error) {
	request := new(PieRequest)

	if err = c.Bind(request); err != nil {
		return
	}

	if len(request.Value) >= 1 && len(request.Title) >= 1 {
		var groups string
		var columns string
		for _, col := range request.Value {
			if columns == "" {
				columns = getColumn(col)
			} else {
				columns = columns + ", " + getColumn(col)
			}

		}
		for _, col := range request.Title {
			if columns == "" {
				columns = getColumn(col)
			} else {
				columns = columns + ", " + getColumn(col)
			}
			if col.GroupBy {
				if groups == "" {
					groups = col.Name
				} else {
					groups = groups + ", " + col.Name
				}
			}
		}

		conditions := ""
		for _, filter := range request.Filters {
			if conditions == "" {
				conditions = getFilter(filter)
			} else {
				conditions = conditions + " AND " + getFilter(filter)
			}
		}

		limitStr := ""
		if len(request.Limit) >= 1 {
			limitStr = request.Limit
		}

		data := GetTableData(request.Value[0].Table, columns, conditions, groups, limitStr, "")
		return c.JSON(http.StatusOK, data)
	} else {
		return c.JSON(http.StatusOK, "[]")
	}
}

func TableData(c echo.Context) (err error) {
	request := new(TableRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	if len(request.Values) >= 1 {
		var columns string
		for _, col := range request.Values {
			if columns == "" {
				columns = col.Name
			} else {
				columns = columns + ", " + col.Name
			}
		}

		data := GetTableData(request.Values[0].Table, columns, "", "", "", "")
		return c.JSON(http.StatusOK, data)
	} else {
		return c.JSON(http.StatusOK, "[]")
	}
}
func getColumn(column Column) string {
	if column.Aggregate == "count" {
		return "count(" + column.Name + ") as " + column.Name
	} else if column.Aggregate == "max" {
		return "max(" + column.Name + ") as " + column.Name
	} else if column.Aggregate == "min" {
		return "min(" + column.Name + ") as " + column.Name
	} else if column.Aggregate == "avg" {
		return "avg(" + column.Name + ") as " + column.Name
	} else if column.Aggregate == "sum" {
		return "sum(" + column.Name + ") as " + column.Name
	} else {
		return column.Name
	}
}

func getValue(value string) string {
	if isNumber(value) {
		return value
	} else {
		return fmt.Sprintf("'%s'", value)
	}
}

func isNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func getFilter(filter Filter) string {
	switch filter.Condition {
	case "equals":
		return filter.Column + " = " + getValue(filter.Value)
	case "notEqual":
		return filter.Column + " != " + getValue(filter.Value)
	case "contains":
		return filter.Column + " LIKE ? " + "%" + filter.Value + "%"
	case "notContains":
		return filter.Column + " not LIKE ? " + "%" + filter.Value + "%"
	case "startsWith":
		return filter.Column + " not LIKE ? " + "%" + filter.Value
	case "endsWith":
		return filter.Column + " not LIKE ? " + filter.Value + "%"
	case "greaterThan":
		return filter.Column + " > " + getValue(filter.Value)
	case "greaterThanOrEqual":
		return filter.Column + " >= " + getValue(filter.Value)
	case "lessThan":
		return filter.Column + " < " + getValue(filter.Value)
	case "lessThanOrEqual":
		return filter.Column + " <= " + getValue(filter.Value)
	case "isNull":
		return filter.Column + "IS NULL"
	case "notNull":
		return filter.Column + "IS NOT NULL"
	case "whereIn":
		return filter.Column + " IN (" + filter.Value + ")"
	default:
		return ""
	}

	return ""
}

func LineData(c echo.Context) (err error) {
	request := new(LineRequest)
	if err = c.Bind(request); err != nil {
		return
	}

	var groups string
	if len(request.Axis) >= 1 && len(request.Lines) >= 1 {
		var columns string
		for _, col := range request.Axis {

			if columns == "" {
				columns = getColumn(col)
			} else {
				columns = columns + ", " + getColumn(col)
			}
			if col.GroupBy {
				if groups == "" {
					groups = col.Name
				} else {
					groups = groups + ", " + col.Name
				}
			}
		}

		for _, col := range request.Lines {
			if columns == "" {
				columns = getColumn(col)
			} else {
				columns = columns + ", " + getColumn(col)
			}
		}

		conditions := ""
		for _, filter := range request.Filters {
			if conditions == "" {
				conditions = getFilter(filter)
			} else {
				conditions = conditions + " AND " + getFilter(filter)
			}
		}

		limitStr := ""
		if len(request.Limit) >= 1 {
			limitStr = request.Limit
		}
		OrderStr := ""
		if len(request.Order) >= 1 {
			OrderStr = request.Order
		}

		data := GetTableData(request.Axis[0].Table, columns, conditions, groups, limitStr, OrderStr)
		return c.JSON(http.StatusOK, data)
	} else {
		return c.JSON(http.StatusOK, "[]")
	}
}

func GetTableData(Table string, Columns string, Condition string, GroupBy string, limitStr string, OrderStr string) []map[string]interface{} {
	data := []map[string]interface{}{}
	filter := ""
	if Condition != "" {
		filter = " WHERE " + Condition
	}

	GroupBySting := ""
	if GroupBy != "" {
		GroupBySting = " GROUP BY " + GroupBy
	}

	LimitQr := ""
	if limitStr != "" {
		LimitQr = " LIMIT " + limitStr
	}
	OrderQr := ""
	if OrderStr != "" {
		OrderQr = " ORDER BY " + OrderStr
	}

	//fmt.Println("SELECT " + Columns + "  FROM " + Table + filter + GroupBySting + OrderQr + LimitQr)
	qr := "SELECT " + Columns + "  FROM " + Table + filter + GroupBySting + OrderQr + LimitQr
	rows, _ := DB.DB.Raw(qr).Rows()

	/*start*/
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
					//	fmt.Println(stringValue)

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
	return data
}

type Column struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	Table      string `json:"table"`
	Alias      string `json:"alias"`
	Output     bool   `json:"output"`
	SortType   string `json:"sortType"`
	SortOrder  int    `json:"sortOrder"`
	GroupBy    bool   `json:"groupBy"`
	GroupOrder int    `json:"groupOrder"`
	Aggregate  string `json:"aggregate"`
	Color      string `json:"color"`
}

type Filter struct {
	Column    string `json:"column"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type LineRequest struct {
	Axis    []Column `json:"axis"`
	Lines   []Column `json:"lines"`
	Filters []Filter `json:"filters"`
	Limit   string   `json:"limit"`
	Order   string   `json:"order"`
}

type CountRequest struct {
	CountFields []Column `json:"countFields"`
	Filters     []Filter `json:"filters"`
}

type PieRequest struct {
	Value   []Column `json:"value"`
	Title   []Column `json:"title"`
	Filters []Filter `json:"filters"`
	Limit   string   `json:"limit"`
}

type TableRequest struct {
	Values  []Column `json:"values"`
	Filters []Filter `json:"filters"`
	Limit   string   `json:"limit"`
	Order   string   `json:"order"`
}
