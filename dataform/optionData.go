package dataform

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func Options(c *fiber.Ctx) error {
	r := new(RalationOption)
	if err := c.BodyParser(r); err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
	}

	Relation := Ralation_{}

	Relation.Fields = r.Fields
	Relation.Filter = r.Filter
	Relation.Key = r.Key
	Relation.SortField = r.SortField
	Relation.SortOrder = r.SortOrder
	Relation.Table = r.Table
	Relation.FilterWithUser = r.FilterWithUser
	Relation.ParentFieldOfForm = r.ParentFieldOfForm
	Relation.ParentFieldOfTable = r.ParentFieldOfTable

	data := OptionsData(Relation, c)
	return c.JSON(data)
}

func OptionsData(relation Ralation_, c *fiber.Ctx) []map[string]interface{} {

	table := relation.Table
	labels := strings.Join(relation.Fields[:], ",', ',")
	key := relation.Key
	sortField := relation.SortField
	sortOrder := relation.SortOrder
	parentFieldOfTable := relation.ParentFieldOfTable
	filter := relation.Filter
	FilterWithUser := relation.FilterWithUser

	//fmt.Println(FilterWithUser)
	data := []map[string]interface{}{}

	if table == "" || len(labels) < 1 || key == "" {
		return data
	}
	var parent_column string
	if parentFieldOfTable != "" {
		parent_column = ", " + parentFieldOfTable + " as parent_value"
	}
	var order_value string
	if sortField != "" && sortOrder != "" {
		order_value = sortField + " " + sortOrder
	}
	var where_value string
	if filter != "" {
		where_value = filter
	}

	if len(FilterWithUser) >= 1 {

		User := agentUtils.AuthUserObject(c)
		for _, userCon := range FilterWithUser {

			tableField := userCon["tableField"]
			userField := User[userCon["userField"]]

			//if userField
			if userField != nil {
				userFieldValue := strconv.FormatInt(reflect.ValueOf(userField).Int(), 10)

				userDataFilter := tableField + " = '" + userFieldValue + "'"

				if userFieldValue != "" && userFieldValue != "0" {
					if where_value == "" {
						where_value = "WHERE " + userDataFilter
					} else {
						where_value = where_value + " AND " + userDataFilter
					}
				}
			}
		}
	}

	//fmt.Println("SELECT " + key + " as value, concat(" + labels + ") as label " + parent_column + "  FROM " + table + " " + where_value + " " + order_value)

	concatTxt := "CONCAT"
	if config.Config.Database.Connection == "mssql" {

	}

	//rows, _ := DB.Query("SELECT " + key + " as value, "+concatTxt+"(" + labels + ") as label " + parent_column + "  FROM " + table + " " + where_value + " " + order_value)
	//
	//fmt.Println("SELECT " + key + " as value, "+concatTxt+"(" + labels + ") as label " + parent_column + "  FROM " + table + " " + where_value + " " + order_value)
	//

	if config.Config.Database.Connection == "oracle" {
		labels = strings.Join(relation.Fields[:], " || ', ' || ")

		if parentFieldOfTable != "" {
			parent_column = ", " + parentFieldOfTable + " as \"parent_value\""
		}

		return GetTableData(key+" as \"value\", "+labels+" as \"label\" "+parent_column, table, where_value, order_value)
	} else {
		return GetTableData(key+" as value, "+concatTxt+"("+labels+") as label "+parent_column, table, where_value, order_value)
	}

	///*start*/
	//
	//columns, _ := rows.Columns()
	//count := len(columns)
	//values := make([]interface{}, count)
	//valuePtrs := make([]interface{}, count)
	//
	///*end*/
	//
	//for rows.Next() {
	//
	//	/*start */
	//
	//	for i := range columns {
	//		valuePtrs[i] = &values[i]
	//	}
	//
	//	rows.Scan(valuePtrs...)
	//
	//	var myMap = make(map[string]interface{})
	//	for i, col := range columns {
	//		val := values[i]
	//
	//		b, ok := val.([]byte)
	//
	//		if (ok) {
	//
	//			v, error := strconv.ParseInt(string(b), 10, 64)
	//			if error != nil {
	//				stringValue := string(b)
	//
	//				myMap[col] = stringValue
	//			} else {
	//				myMap[col] = v
	//			}
	//
	//		}
	//
	//	}
	//	/*end*/
	//
	//	data = append(data, myMap)
	//
	//}
	return data
}

type FormOption struct {
	Label       interface{} `gorm:"column:label" json:"label"`
	Value       int         `gorm:"column:value;type:uuid" json:"value"`
	ParentValue interface{} `gorm:"column:parent_value" json:"parent_value"`
}

func GetTableData(query string, table string, where_value string, order_value string) []map[string]interface{} {
	var data []map[string]interface{}

	err := DB.DB.Table(table).Select(query).Where(where_value).Order(order_value).Find(&data).Error

	if err != nil {
		fmt.Println(err.Error())
	}
	//start
	//
	//columns, _ := rows.Columns()
	//
	//count := len(columns)
	//values := make([]interface{}, count)
	//valuePtrs := make([]interface{}, count)
	//
	///*end*/
	//
	//for rows.Next() {
	//
	//	/*start */
	//
	//	for i := range columns {
	//		valuePtrs[i] = &values[i]
	//	}
	//
	//	rows.Scan(valuePtrs...)
	//
	//	var myMap = make(map[string]interface{})
	//	for i, col := range columns {
	//
	//		val := values[i]
	//
	//		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
	//			if col == "id"{
	//				b, ok := val.([]byte)
	//				if ok {
	//					stringValue := string(b)
	//					myMap[col] = stringValue
	//				} else {
	//					myMap[col] = val
	//				}
	//			} else {
	//				myMap[col] = val
	//			}
	//
	//		} else {
	//			b, ok := val.([]byte)
	//
	//			if ok {
	//
	//				v, error := strconv.ParseInt(string(b), 10, 64)
	//				if error != nil {
	//					stringValue := string(b)
	//
	//					myMap[col] = stringValue
	//				} else {
	//					myMap[col] = v
	//				}
	//
	//			}
	//		}
	//
	//	}
	//	/*end*/
	//
	//	data = append(data, myMap)
	//
	//}

	if config.Config.Database.Connection == "oracle" {
		if len(data) >= 1 {
			switch data[0]["value"].(type) {
			case string:
				_, parseError := strconv.ParseInt(data[0]["value"].(string), 10, 64)
				if parseError == nil {
					for i, row := range data {
						v, parseErr := strconv.ParseInt(row["value"].(string), 10, 64)
						if parseErr == nil {
							data[i]["value"] = v
						}
					}
				}
			default:
			}
		}
	}

	return data

}
