package datagrid

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"github.com/labstack/echo/v4"
)

func CallTrigger(action string, datagrid Datagrid, data []map[string]interface{}, id string, query *gorm.DB, c echo.Context) ([]map[string]interface{}, *gorm.DB, bool) {

	if len(datagrid.Triggers) <= 0 || datagrid.TriggerNameSpace == ""{
		return data, query, false
	}


	switch action {
	case "afterDelete":
		Method := datagrid.Triggers["afterDelete"].(string)
		Struct := datagrid.Triggers["afterDeleteStruct"]
		return execTrigger(action, Method, Struct, datagrid.DataModel, data, id, query, c)
	case "beforeFetch":
		Method := datagrid.Triggers["beforeFetch"].(string)
		Struct := datagrid.Triggers["beforeFetchStruct"]
		return execTrigger(action, Method, Struct, datagrid.DataModel, data, id, query, c)
	case "beforeDelete":
		Method := datagrid.Triggers["beforeDelete"].(string)
		Struct := datagrid.Triggers["beforeDeleteStruct"]
		return execTrigger(action, Method, Struct, datagrid.DataModel, data, id, query, c)
	case "beforePrint":
		Method := datagrid.Triggers["beforePrint"].(string)
		Struct := datagrid.Triggers["beforePrintStruct"]
		return execTrigger(action, Method, Struct, datagrid.DataModel, data, id, query, c)

	}

	return data, query, false

}

func execTrigger(action string, triggerMethod string, triggerStruct interface{}, Model interface{}, data []map[string]interface{}, id string, query *gorm.DB, c echo.Context, ) ([]map[string]interface{}, *gorm.DB, bool) {

	if triggerMethod != "" {
		triggerMethod_ := reflect.ValueOf(triggerStruct).MethodByName(triggerMethod)

		if action == "afterDelete" {
			if triggerMethod_.IsValid() {

				input := make([]reflect.Value, 3)
				input[0] = reflect.ValueOf(Model)
				input[1] = reflect.ValueOf(data)
				input[2] = reflect.ValueOf(id)
				triggerMethodRes := triggerMethod_.Call(input)

				return triggerMethodRes[0].Interface().([]map[string]interface{}), query, false
			}
			return data, query, false
		} else {
			if triggerMethod_.IsValid() {

				input := make([]reflect.Value, 5)
				input[0] = reflect.ValueOf(Model)
				input[1] = reflect.ValueOf(data)
				input[2] = reflect.ValueOf(id)
				input[3] = reflect.ValueOf(query)
				input[4] = reflect.ValueOf(c)
				triggerMethodRes := triggerMethod_.Call(input)

				return triggerMethodRes[0].Interface().([]map[string]interface{}), triggerMethodRes[1].Interface().(*gorm.DB), triggerMethodRes[2].Interface().(bool)
			}
			return data, query, false
		}

	} else {
		return data, query, false
	}

}
