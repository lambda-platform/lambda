package dataform

import (
	"encoding/json"
	"github.com/PaesslerAG/gval"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DBSchema"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/utils"
	"reflect"
	"regexp"
)

func GetData(c *fiber.Ctx, action string, id string, dataform Dataform) (*map[string]interface{}, error) {
	
	requestData := new(map[string]interface{})

	bodyBytes := utils.GetBody(c)
	json.Unmarshal([]byte(bodyBytes), requestData)

	if err := c.BodyParser(dataform.Model); err != nil {
		return requestData, err
	}

	if len(dataform.FieldTypes) >= 1 {

		for field, fieldType := range dataform.FieldTypes {

			if fieldType == "Password" {
				fieldName := DBSchema.FieldName(field)
				value := dataform.getStringField(fieldName)

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
