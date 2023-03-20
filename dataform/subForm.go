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

				parentIdentityType, err := dataform.getFieldType(DBSchema.FieldName(parentIdentity))

				var parentId interface{}
				if parentIdentityType == "string" {
					preParentID, preParentIDErr := dataform.getStringField(DBSchema.FieldName(parentIdentity))
					if err == preParentIDErr {
						parentId = preParentID
					}

				} else {
					preParentID, preParentIDErr := dataform.getIntField(DBSchema.FieldName(parentIdentity))
					if err == preParentIDErr {
						parentId = preParentID
					}
				}
				Clear(subForm.Model)

				existingIDS := []interface{}{}

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

						var err error = nil

						createNewRow := true

						if subIdentityValue != nil {
							switch vtype := subIdentityValue.(type) {
							case string:
								if subIdentityValue.(string) != "" {
									createNewRow = false
								}
							case int:
								if subIdentityValue.(int) >= 1 {
									createNewRow = false
								}
							case float64:
								if subIdentityValue.(float64) >= 1 {
									createNewRow = false
								}
							case float32:
								if subIdentityValue.(float32) >= 1 {
									createNewRow = false
								}
							case int64:
								if subIdentityValue.(int64) >= 1 {
									createNewRow = false
								}
							case int32:
								if subIdentityValue.(int32) >= 1 {
									createNewRow = false
								}
							default:

								fmt.Println(vtype)
							}
						}
						if createNewRow {
							err = DB.DB.Create(subForm.Model).Error
						} else {
							dataform.setModelField(subForm.Model, DBSchema.FieldName(subForm.Identity), subIdentityValue)
							err = DB.DB.Save(subForm.Model).Error
						}

						currentID, currentIDErr := dataform.getModelFieldValue(subForm.Model, DBSchema.FieldName(subForm.Identity))

						fmt.Println(currentIDErr)
						existingIDS = append(existingIDS, currentID)

						if err == nil {
							//	CallTrigger("afterUpdate", subForm, subD, "")
							saveNestedSubItem(subForm, subD)
						}

					}
				}
				Clear(subForm.Model)
				if tableTypeColumn != "" && tableTypeValue != "" {
					DB.DB.Where(connectionField+" = ? AND "+tableTypeColumn+" = ? AND "+subIdentity+" NOT IN ?", parentId, tableTypeValue, existingIDS).Unscoped().Delete(subForm.Model)
				} else {

					DB.DB.Table(subForm.Table).Where(connectionField+" = ? AND "+subIdentity+" NOT IN ?", parentId, existingIDS).Unscoped().Delete(subForm.Model)
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
