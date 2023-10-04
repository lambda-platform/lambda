package dataform

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/utils"
)

func Edit(c *fiber.Ctx, dataform Dataform, id string) error {
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
			for _, dataS := range dataSub {

				subIdentity := Sub["subIdentity"].(string)

				parentId := utils.GetString(dataS[subIdentity])

				for _, Sub2 := range subForm.SubForms {

					connectionField2 := Sub2["connection_field"].(string)
					tableTypeColumn2 := Sub2["tableTypeColumn"].(string)
					tableTypeValue2 := Sub2["tableTypeValue"].(string)
					subTable2 := Sub2["table"].(string)
					subForm2Array := Sub2["subFormArray"]

					dataS[subTable2] = []interface{}{}

					if tableTypeColumn2 != "" && tableTypeValue2 != "" {
						DB.DB.Where(connectionField+" = ? AND "+tableTypeColumn2+" = ?", id, tableTypeValue2).Find(subForm2Array)
					} else {
						DB.DB.Where(connectionField2+" = ?", parentId).Find(subForm2Array)
					}

					dataS[subTable2] = subForm2Array

				}
				dataSubNew := map[string]interface{}{}
				subDataNew, _ := json.Marshal(dataS)
				json.Unmarshal(subDataNew, &dataSubNew)

				dataWitSub = append(dataWitSub, dataSubNew)
			}

			data[subTable] = dataWitSub

		}

		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   data,
		})

	} else {
		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   dataform.Model,
		})
	}

}
