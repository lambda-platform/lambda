package dataform

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
	"fmt"
)

func Edit(c echo.Context, dataform Dataform, id string) error {
	DB.DB.Where(dataform.Identity+" = ?", id).Find(dataform.Model)
	if len(dataform.SubForms) >= 1 {



		data := make(map[string]interface{})
		dataPre, _ := json.Marshal(dataform.Model)
		json.Unmarshal(dataPre, &data)

		for _, Sub := range dataform.SubForms {

			connectionField := Sub["connection_field"].(string)
			tableTypeColumn := Sub["tableTypeColumn"].(string)
			tableTypeValue := Sub["tableTypeValue"].(string)
			subTable := Sub["table"].(string)
			subForm := Sub["subForm"].(Dataform)
			subFormArray := Sub["subFormArray"]

			if tableTypeColumn != "" && tableTypeValue != "" {
				DB.DB.Where(connectionField+" = ? AND "+tableTypeColumn+" = ?", id, tableTypeValue).Find(subFormArray)
			} else {
				DB.DB.Where(connectionField+" = ?", id).Find(subFormArray)
			}

			dataSub := []map[string]interface{}{}
			subData, _ := json.Marshal(subFormArray)
			json.Unmarshal(subData, &dataSub)

			dataWitSub := []map[string]interface{}{}
			for _, sData := range dataSub {

				subIdentity := Sub["subIdentity"].(string)

				parentId := fmt.Sprintf("%g", sData[subIdentity])


				for _, Sub2 := range subForm.SubForms {

					connectionField2 := Sub2["connection_field"].(string)
					tableTypeColumn := Sub2["tableTypeColumn"].(string)
					tableTypeValue := Sub2["tableTypeValue"].(string)
					subTable2 := Sub2["table"].(string)
					subForm2Array := Sub2["subFormArray"]
					if tableTypeColumn != "" && tableTypeValue != "" {
						DB.DB.Where(connectionField+" = ? AND "+tableTypeColumn+" = ?", id, tableTypeValue).Find(Sub2["subForm"])
					} else {
						DB.DB.Where(connectionField2+" = ?", parentId).Find(subForm2Array)
					}

					sData[subTable2] = subForm2Array

				}

				dataWitSub = append(dataWitSub, sData)
			}

			data[subTable] = dataWitSub

		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "true",
			"data":   data,
		})

	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "true",
			"data":   dataform.Model,
		})
	}


}
