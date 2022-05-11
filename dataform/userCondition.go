package dataform

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	lbModel "github.com/lambda-platform/lambda/models"
	"net/http"
	"net/url"
	"reflect"
)

func SetCondition(condition string, c echo.Context, VBSchema lbModel.VBSchema) error {

	con, _ := url.ParseQuery(condition)
	var schema lbModel.SCHEMA

	json.Unmarshal([]byte(VBSchema.Schema), &schema)

	for uC, _ := range con {

		uString := reflect.ValueOf(uC).Interface().(string)

		var conditionData []map[string]string
		json.Unmarshal([]byte(uString), &conditionData)

		User := agentUtils.AuthUserObject(c)
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "true",
		"data":   VBSchema,
	})
}
