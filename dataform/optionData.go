package dataform

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/models"
	"strconv"
	"strings"
)

func Options(c *fiber.Ctx, dataform Dataform) error {

	var optionsData = map[string][]map[string]interface{}{}

	for table, relation := range dataform.Relations {
		data := OptionsData(relation, c)

		optionsData[table] = data

	}

	return c.JSON(optionsData)
}

func OptionsData(relation models.Relation, c *fiber.Ctx) []map[string]interface{} {

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

	if FilterWithUser != nil {

		User, err := agentUtils.AuthUserObject(c)

		if err != nil {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  err.Error(),
				"status": false,
			})
		}
		for _, userCon := range *FilterWithUser {
			fmt.Println(userCon)
			tableField := userCon.TableField
			userFieldValue := User[userCon.UserField]

			//if userField
			if userFieldValue != nil {

				userDataFilter := tableField + " = '" + fmt.Sprintf("%v", userFieldValue) + "'"

				if where_value == "" {
					where_value = userDataFilter
				} else {
					where_value = where_value + " AND " + userDataFilter
				}
			}
		}
	}

	//fmt.Println("SELECT " + key + " as value, concat(" + labels + ") as label " + parent_column + "  FROM " + table + " " + where_value + " " + order_value)

	concatTxt := "CONCAT"
	if config.Config.Database.Connection == "mssql" {
		if len(relation.Fields) == 1 {
			concatTxt = ""
		}
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
