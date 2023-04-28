package dataform

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	lbModel "github.com/lambda-platform/lambda/models"
	"net/url"
	"reflect"
)

func SetCondition(condition string, c *fiber.Ctx, VBSchema lbModel.VBSchema) error {

	con, _ := url.ParseQuery(condition)
	var schema lbModel.SCHEMA

	json.Unmarshal([]byte(VBSchema.Schema), &schema)

	for uC, _ := range con {

		uString := reflect.ValueOf(uC).Interface().(string)

		var conditionData []map[string]string
		json.Unmarshal([]byte(uString), &conditionData)

		User, err := agentUtils.AuthUserObject(c)

		if err != nil {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  err.Error(),
				"status": false,
			})
		}
		for _, userCondition := range conditionData {

			for i := range schema.Schema {

				if schema.Schema[i].Model == userCondition["form_field"] {
					schema.Schema[i].Disabled = true
					schema.Schema[i].Default = User[userCondition["user_field"]]
				}

			}
		}

	}
	schemaString, _ := json.Marshal(schema)
	VBSchema.Schema = string(schemaString)

	return c.JSON(map[string]interface{}{
		"status": true,
		"data":   VBSchema,
	})
}

func SetConditionOracle(condition string, c *fiber.Ctx, VBSchema lbModel.VBSchemaOracle) error {

	con, _ := url.ParseQuery(condition)
	var schema lbModel.SCHEMA

	json.Unmarshal([]byte(VBSchema.Schema), &schema)

	for uC, _ := range con {

		uString := reflect.ValueOf(uC).Interface().(string)

		var conditionData []map[string]string
		json.Unmarshal([]byte(uString), &conditionData)

		User, err := agentUtils.AuthUserObject(c)

		if err != nil {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  err.Error(),
				"status": false,
			})
		}
		for _, userCondition := range conditionData {

			for i := range schema.Schema {

				if schema.Schema[i].Model == userCondition["form_field"] {
					schema.Schema[i].Disabled = true
					schema.Schema[i].Default = User[userCondition["user_field"]]
				}

			}
		}

	}
	schemaString, _ := json.Marshal(schema)
	VBSchema.Schema = string(schemaString)

	return c.JSON(map[string]interface{}{
		"status": true,
		"data":   VBSchema,
	})
}
