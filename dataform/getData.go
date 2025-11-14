package dataform

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"github.com/PaesslerAG/gval"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DBSchema"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/utils"
	"github.com/pkg/errors"
)

func GetData(c *fiber.Ctx, action string, id string, dataform Dataform) (*map[string]interface{}, error) {

	requestData := new(map[string]interface{})

	if err := c.BodyParser(dataform.Model); err != nil {
		return requestData, err
	}

	bodyBytes := utils.GetBody(c)
	json.Unmarshal([]byte(bodyBytes), requestData)

	var optionsData = map[string][]map[string]interface{}{}

	for relationKey, relation := range dataform.Relations {
		data := OptionsData(relation, c)

		optionsData[relationKey] = data

	}

	if len(dataform.FieldTypes) >= 1 {

		for field, fieldType := range dataform.FieldTypes {

			if fieldType == "Select" {

				requestValue := (*requestData)[field]

				if requestValue != nil {
					// Loop through relations
					for relationKey, relation := range dataform.Relations {
						if relation.TargetField == field {

							allowed := false
							for _, opt := range optionsData[relationKey] {
								if optVal, ok := opt["value"]; ok {

									switch optVal.(type) {
									case string:
										if requestValue == optVal {
											allowed = true
											break
										}
									case float64:
										// JSON numbers are float64 by default
										if fmt.Sprintf("%v", requestValue) == fmt.Sprintf("%v", optVal) {
											allowed = true
											break
										}
									case float32:
										// JSON numbers are float64 by default
										if fmt.Sprintf("%v", requestValue) == fmt.Sprintf("%v", optVal) {
											allowed = true
											break
										}
									case int64:
										if fmt.Sprintf("%v", requestValue) == fmt.Sprintf("%v", optVal) {
											allowed = true
											break
										}
									case int32:
										if fmt.Sprintf("%v", requestValue) == fmt.Sprintf("%v", optVal) {
											allowed = true
											break
										}
									case int16:
										if fmt.Sprintf("%v", requestValue) == fmt.Sprintf("%v", optVal) {
											allowed = true
											break
										}
									case int:
										if fmt.Sprintf("%v", requestValue) == fmt.Sprintf("%v", optVal) {
											allowed = true
											break
										}
									}
								}
							}

							if !allowed {
								msg := map[string]interface{}{
									"status": false,
									"error":  fmt.Sprintf("Value %v not allowed for field %s", requestValue, field),
								}
								fmt.Println(fmt.Sprintf("Value %v not allowed for field %s", requestValue, field))
								return &msg, errors.New("Value not allowed")
							}
						}
					}
				}

			}
			if fieldType == "Password" {
				fieldName := DBSchema.FieldName(field)
				value, err := dataform.getStringField(fieldName)

				if err != nil {
					return requestData, err
				}

				if action == "store" {
					password, _ := agentUtils.Hash(value)

					dataform.setStringField(fieldName, password)
				} else {

					if value == "" {
						delete(dataform.ValidationRules, field)
					} else {
						password, _ := agentUtils.Hash(value)
						dataform.setStringField(fieldName, password)
					}

				}
			}

			if fieldType == "Date" || fieldType == "DateTime" {
				if reflect.TypeOf((*requestData)[field]) != nil {
					if (*requestData)[field].(string) != "" {
						delete(dataform.ValidationRules, field)
					}
				}

			}
		}
	}
	//
	if id != "" {
		if dataform.TriggerNameSpace != "" && dataform.BeforeUpdate != nil {
			dataform.BeforeUpdate(dataform.Model)
		}
	} else {
		if dataform.TriggerNameSpace != "" && dataform.BeforeInsert != nil {
			dataform.BeforeInsert(dataform.Model)
		}
	}

	if len(dataform.Formulas) >= 1 {

		for _, formula := range dataform.Formulas {

			for _, target := range formula.Targets {
				if target.Prop == "hidden" {

					var re0 = regexp.MustCompile("'{" + formula.Model + "}'")
					template := re0.ReplaceAllString(formula.Template, DBSchema.FieldName(formula.Model))

					var re3 = regexp.MustCompile(`'`)
					template = re3.ReplaceAllString(template, `"`)

					value, _ := gval.Evaluate(template, dataform.Model)

					if value == true {

						delete(dataform.ValidationRules, target.Field)
					}

				}
			}

		}
	}

	return requestData, nil
}
