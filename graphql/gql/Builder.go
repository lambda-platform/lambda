package gql

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"gorm.io/gorm"
	"strings"
)

func Order(sortsPre interface{}, query *gorm.DB, columns []string) (*gorm.DB, error) {
	var sorts []map[string]string
	order, err := json.Marshal(sortsPre)
	if err != nil {
		return query, errors.New("Please insert correct sort values")
	}
	err2 := json.Unmarshal(order, &sorts)
	if err2 != nil {
		return query, errors.New("Please insert correct sort values")
	}
	for _, sort := range sorts {

		errCol := CheckColumns(sort["column"], columns)
		if errCol != nil {
			return query, errCol
		}
		query = query.Order(sort["column"] + " " + sort["order"])

	}
	return query, nil

}
func Filter(filtersPre interface{}, query *gorm.DB, columns []string) (*gorm.DB, error) {

	var filters []map[string]string
	order, err := json.Marshal(filtersPre)
	if err != nil {
		return query, errors.New("Please insert correct filter value")
	}
	err2 := json.Unmarshal(order, &filters)

	if err2 != nil {
		return query, errors.New("Please insert correct filter value")
	}

	if len(filters) >= 1 {
		for _, filter := range filters {

			errCol := CheckColumns(filter["column"], columns)

			if errCol != nil {
				return query, errCol
			}
			k := filter["column"]
			v := filter["value"]

			if v != "" {
				switch filter["condition"] {

				case "equals":
					query = query.Where(k+" = ?", v)
				case "notEqual":
					query = query.Where(k+" != ?", v)
				case "contains":
					query = query.Where("LOWER("+k+") LIKE ?", "%"+strings.ToLower(v)+"%")
				case "notContains":
					query = query.Where("LOWER("+k+") not LIKE ?", "%"+strings.ToLower(v)+"%")
				case "startsWith":
					query = query.Where("LOWER("+k+")  LIKE ?", strings.ToLower(v)+"%")
				case "endsWith":
					query = query.Where("LOWER("+k+")  LIKE ?", "%"+strings.ToLower(v))
				case "greaterThan":
					query = query.Where(k+" > ?", v)
				case "greaterThanOrEqual":
					query = query.Where(k+" >= ?", v)
				case "lessThan":
					query = query.Where(k+" < ?", v)
				case "lessThanOrEqual":
					query = query.Where(k+" <= ?", v)
				case "isNull":
					query = query.Where(k + " IS NULL")
				case "notNull":
					query = query.Where(k + " IS NOT NULL")
				case "whereIn":
					query = query.Where(k+" IN (?)", strings.Split(v, ","))

				default:
					return query, errors.New(filter["condition"] + ": is wrong condition")
				}

			}

		}

	}

	return query, nil
}
func GroupFilter(filterGroupsPre interface{}, queryMain *gorm.DB, columns []string) (*gorm.DB, error) {

	var filterGroups []map[string]interface{}
	order, err := json.Marshal(filterGroupsPre)
	if err != nil {
		return queryMain, errors.New("Please insert correct filter group value")
	}
	err2 := json.Unmarshal(order, &filterGroups)

	if err2 != nil {
		return queryMain, errors.New("Please insert correct filter group value")
	}

	if len(filterGroups) >= 1 {
		query := DB.DB
		for _, filterGroup := range filterGroups {

			var combine string = filterGroup["combine"].(string)

			if combine == "or" || combine == "and" {

				var filters []map[string]string
				order, err := json.Marshal(filterGroup["filters"])
				if err != nil {
					return queryMain, errors.New("Please insert correct filter value")
				}
				err2 := json.Unmarshal(order, &filters)

				if err2 != nil {
					return queryMain, errors.New("Please insert correct filter value")
				}

				if len(filters) >= 1 {
					for filterIndex, filter := range filters {

						errCol := CheckColumns(filter["column"], columns)

						if errCol != nil {
							return query, errCol
						}
						k := filter["column"]
						v := filter["value"]

						if v != "" {
							switch filter["condition"] {

							case "equals":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k+" = ?", v)
								} else {
									query = query.Where(k+" = ?", v)
								}

							case "notEqual":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k+" != ?", v)
								} else {
									query = query.Where(k+" != ?", v)
								}

							case "contains":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or("LOWER("+k+") LIKE ?", "%"+strings.ToLower(v)+"%")
								} else {
									query = query.Where("LOWER("+k+") LIKE ?", "%"+strings.ToLower(v)+"%")
								}

							case "notContains":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or("LOWER("+k+") not LIKE ?", "%"+strings.ToLower(v)+"%")
								} else {
									query = query.Where("LOWER("+k+") not LIKE ?", "%"+strings.ToLower(v)+"%")
								}

							case "startsWith":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or("LOWER("+k+")  LIKE ?", strings.ToLower(v)+"%")
								} else {
									query = query.Where("LOWER("+k+")  LIKE ?", strings.ToLower(v)+"%")
								}

							case "endsWith":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or("LOWER("+k+")  LIKE ?", "%"+strings.ToLower(v))
								} else {
									query = query.Where("LOWER("+k+")  LIKE ?", "%"+strings.ToLower(v))
								}

							case "greaterThan":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k+" > ?", v)
								} else {
									query = query.Where(k+" > ?", v)
								}

							case "greaterThanOrEqual":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k+" >= ?", v)
								} else {
									query = query.Where(k+" >= ?", v)
								}

							case "lessThan":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k+" < ?", v)
								} else {
									query = query.Where(k+" < ?", v)
								}

							case "lessThanOrEqual":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k+" <= ?", v)
								} else {
									query = query.Where(k+" <= ?", v)
								}

							case "isNull":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k + " IS NULL")
								} else {
									query = query.Where(k + " IS NULL")
								}

							case "notNull":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k + " IS NOT NULL")
								} else {
									query = query.Where(k + " IS NOT NULL")
								}

							case "whereIn":
								if filterIndex >= 1 && combine == "or" {
									query = query.Or(k+" IN (?)", strings.Split(v, ","))
								} else {
									query = query.Where(k+" IN (?)", strings.Split(v, ","))
								}

							default:
								return query, errors.New(filter["condition"] + ": is wrong condition")
							}
						}

					}

				}

			} else {
				return queryMain, errors.New("combine value must be OR, AND")
			}

		}

		return queryMain.Where(query), nil
	} else {
		return queryMain, nil
	}

}
func CheckColumns(column string, columns []string) error {
	for _, value := range columns {
		if value == column {
			return nil
		}
	}
	return errors.New(column + ": Column not found")
}

type CustomContext struct {
	*fiber.Ctx
	ctx context.Context
}

func Process() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.Context(), "FiberContextKey", c)
		c.SetUserContext(ctx)

		cc := &CustomContext{c, ctx}

		return cc.Ctx.Next()
	}
}
