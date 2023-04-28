package datagrid

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func Filter(c *fiber.Ctx, datagrid Datagrid, query *gorm.DB) (*gorm.DB, string) {

	customHeader := ""
	bodyBytes := utils.GetBody(c)
	var filterData map[string]interface{}
	json.Unmarshal([]byte(bodyBytes), &filterData)

	if len(filterData) >= 1 {

		if val, ok := filterData["customHeader"]; ok {
			customHeader = val.(string)
		}

		for k, v := range filterData {
			if k == "user_condition" {

				for _, userCondition := range v.([]interface{}) {
					codintion := reflect.ValueOf(userCondition).Interface().(map[string]interface{})
					User, err := agentUtils.AuthUserObject(c)

					if err != nil {
						c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
							"error":  err.Error(),
							"status": false,
						})
					}

					query = query.Where(codintion["grid_field"].(string)+" = ?", User[codintion["user_field"].(string)])
				}

			} else if k == "custom_condition" {

				if reflect.TypeOf(v).String() == "map[string]interface {}" {
					for kc, vc := range v.(map[string]interface{}) {
						query = query.Where(kc+" = ?", fmt.Sprintf("%v", vc))
					}
				} else {
					query = query.Where(fmt.Sprintf("%v", v))
				}

			} else {
				filterType := datagrid.Filters[k]

				if filterType != "" {
					switch filterType {
					case "Select":
						query = query.Where(k+" = ?", v)
					case "Tag":
						query = query.Where(k+" IN (?)", v)
					case "Date":
						if config.Config.Database.Connection == "oracle" {
							query = query.Where(k+" = TO_DATE(?,'YYYY-MM-DD')", fmt.Sprintf("%v", v))
						} else {
							query = query.Where(k+" = ?", fmt.Sprintf("%v", v))
						}

					case "DateRange":

						if config.Config.Database.Connection == "oracle" {
							query = query.Where(k+" BETWEEN TO_DATE(?,'YYYY-MM-DD') AND TO_DATE(?,'YYYY-MM-DD')", reflect.ValueOf(v).Index(0).Interface().(string), reflect.ValueOf(v).Index(1).Interface().(string))
						} else {
							query = query.Where(k+" BETWEEN ? AND ?", reflect.ValueOf(v).Index(0).Interface().(string), reflect.ValueOf(v).Index(1).Interface().(string))
						}
					case "DateRangeDouble":
						start := reflect.ValueOf(v).Index(0).Interface().(string)
						end := reflect.ValueOf(v).Index(1).Interface().(string)
						if start != "" && end != "" {
							query = query.Where(k+" BETWEEN ? AND ?", start, end)
						} else if start != "" && end == "" {
							query = query.Where(k+" >= ?", start)
						} else if start == "" && end != "" {
							query = query.Where(k+" <= ?", end)
						}

					default:
						switch vtype := v.(type) {
						case map[string]interface{}:
							fmt.Println(vtype)
							vmap := v.(map[string]interface{})
							switch vmap["type"] {
							case "contains":
								query = query.Where("LOWER("+k+") LIKE ?", "%"+strings.ToLower(fmt.Sprintf("%v", vmap["filter"]))+"%")
							case "equals":
								query = query.Where(k+" = ?", fmt.Sprintf("%v", vmap["filter"]))
							case "lessThan":
								query = query.Where(k+" <= ?", fmt.Sprintf("%v", vmap["filter"]))
							case "greaterThan":
								query = query.Where(k+" >= ?", fmt.Sprintf("%v", vmap["filter"]))
							case "notContains":
								query = query.Where(k+" != ?", fmt.Sprintf("%v", vmap["filter"]))
							default:
								query = query.Where(k+" = ?", fmt.Sprintf("%v", vmap["filter"]))
								//query = query.Where("LOWER("+k+") LIKE ?", "%"+strings.ToLower(fmt.Sprintf("%v", v))+"%")
							}
						default:
							if filterType == "Text" {

								query = query.Where("LOWER("+k+") LIKE ?", "%"+strings.ToLower(fmt.Sprintf("%v", v))+"%")

							} else {
								query = query.Where(k+" = ?", fmt.Sprintf("%v", v))

							}

						}

					}
				} else {
					query = query.Where(k+" = ?", fmt.Sprintf("%v", v))
				}
			}

		}

	}

	return query, customHeader
}

func Search(c *fiber.Ctx, GridModel interface{}, query *gorm.DB) *gorm.DB {

	search := c.Query("search")

	if search != "" {

		GetColumns := reflect.ValueOf(GridModel).MethodByName("GetColumns")
		columnsPre := GetColumns.Call([]reflect.Value{})
		columns := columnsPre[0].Interface().(map[int]map[string]string)

		i := 0
		for _, c := range columns {
			if i <= 0 {
				query = query.Where(c["column"]+" LIKE ?", "%"+search+"%")
			} else {
				//query = query.Or(c+" LIKE ?", "%"+search+"%")
				//query = query.Where(c+" LIKE ?", "%"+search+"%")
			}
			i++
		}

	}

	return query
}
