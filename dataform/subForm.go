package dataform

import (
	"encoding/json"
	"fmt"
	"github.com/lambda-platform/lambda/DB"
	"regexp"
	"strconv"

	"github.com/lambda-platform/lambda/DBSchema"
	"github.com/lambda-platform/lambda/config"
	"strings"
)

func saveNestedSubItem(dataform Dataform, data map[string]interface{}) {

	if len(dataform.SubForms) >= 1 {

		for _, Sub := range dataform.SubForms {
			table := Sub["table"].(string)
			parentIdentity := Sub["parentIdentity"].(string)
			subIdentity := Sub["subIdentity"].(string)
			connectionField := Sub["connection_field"].(string)
			tableTypeColumn := Sub["tableTypeColumn"].(string)
			tableTypeValue := Sub["tableTypeValue"].(string)
			subForm := Sub["subForm"].(Dataform)

			subData := data[table]

			if subData != nil {

				parentIdentityType := dataform.getFieldType(DBSchema.FieldName(parentIdentity))

				var parentId interface{}
				if parentIdentityType == "string" {
					parentId = dataform.getStringField(DBSchema.FieldName(parentIdentity))
				} else {
					parentId = dataform.getIntField(DBSchema.FieldName(parentIdentity))
				}

				if tableTypeColumn != "" && tableTypeValue != "" {
					DB.DB.Where(connectionField+" = ? AND "+tableTypeColumn+" = ?", parentId, tableTypeValue).Unscoped().Delete(subForm.Model)
				} else {
					DB.DB.Where(connectionField+" = ?", parentId).Unscoped().Delete(subForm.Model)
				}

				currentData := subData.([]interface{})

				if len(currentData) >= 1 {

					for _, sData := range currentData {

						subD := sData.(map[string]interface{})

						subIdentityValue := subD[subIdentity]

						subD[connectionField] = parentId
						if tableTypeColumn != "" && tableTypeValue != "" {
							if IsInt(tableTypeValue) {
								intVar, _ := strconv.Atoi(tableTypeValue)
								subD[tableTypeColumn] = intVar
							} else {

								subD[tableTypeColumn] = tableTypeValue
							}

						}

						if parentIdentityType != "string" {
							if subIdentityValue == nil || config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
								subD[subIdentity] = 0
							}
						}

						Clear(subForm.Model)
						saveData, _ := json.Marshal(subD)
						json.Unmarshal(saveData, subForm.Model)

						err := DB.DB.Create(subForm.Model).Error

						//err := DB.DB.Save(subForm).Error

						if err == nil {
							//	CallTrigger("afterUpdate", subForm, subD, "")

							saveNestedSubItem(subForm, subD)
						}

						/*creareNewRow := true


						switch vtype := subIdentityValue.(type) {
						case int:
							if(subIdentityValue.(int) >= 1){
								creareNewRow = false

							}
						case float64:
							if(subIdentityValue.(float64) >= 1){
						DB.Date		creareNewRow = false

							}
						case float32:
							if(subIdentityValue.(float32) >= 1){
								creareNewRow = false

							}
						case int64:

							if(subIdentityValue.(int64) >= 1){
								creareNewRow = false

							}
						default:

							fmt.Println(vtype)
						}


						if (!creareNewRow){


							err := DB.DB.Save(subForm).Error

							if err == nil {
								CallTrigger("afterUpdate", subForm, subD, "")

								saveNestedSubItem(subForm, subD)
							}

						} else {

							DB.DB.NewRecord(subForm)
							err := DB.DB.Create(subForm).Error

							if err == nil {

								CallTrigger("afterInsert", subForm, subD, "")
								saveNestedSubItem(subForm, subD)

							}

						}*/

					}
				}

			}

		}
	}

}

func IsInt(s string) bool {
	l := len(s)
	if strings.HasPrefix(s, "-") {
		l = l - 1
		s = s[1:]
	}

	reg := fmt.Sprintf("\\d{%d}", l)

	rs, err := regexp.MatchString(reg, s)

	if err != nil {
		return false
	}

	return rs
}
