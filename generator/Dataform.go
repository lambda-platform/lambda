package generator

import (
	"encoding/json"
	"fmt"
	"github.com/lambda-platform/lambda/DBSchema"
	genertarModels "github.com/lambda-platform/lambda/generator/models"
	"github.com/lambda-platform/lambda/generator/utils"
	lambdaModels "github.com/lambda-platform/lambda/models"
	"strconv"
	"strings"
)

func WriteFormsModelData(dbSchema lambdaModels.DBSCHEMA, schemas []genertarModels.ProjectSchemas, copyClienModels bool) {
	WriteFormModel(dbSchema, schemas)
	WriteModelCaller(dbSchema, schemas, copyClienModels)
}

func WriteFormModel(dbSchema lambdaModels.DBSCHEMA, schemas []genertarModels.ProjectSchemas) {

	for _, vb := range schemas {
		var schema lambdaModels.SCHEMA

		json.Unmarshal([]byte(vb.Schema), &schema)

		modelAlias := GetModelAlias(schema.Model)
		modelAliasWithID := modelAlias + strconv.FormatInt(int64(vb.ID), 10)
		//DB_ := DB.DBConnection()

		hiddenColumns := []string{}

		for _, column := range schema.Schema {
			if (column.Hidden == true && column.Default == nil && column.Label == "") || (column.Hidden == true && column.Default == "" && column.Label == "") {
				hiddenColumns = append(hiddenColumns, column.Model)
			}
		}

		fmt.Println(modelAlias)

		columnDataTypes := GetColumnsFromTableMeta(dbSchema.TableMeta[schema.Model], hiddenColumns)

		gormStructs := ""
		for _, field := range schema.Schema {
			if field.FormType == "SubForm" {
				if field.SubType != "Form" {
					subAlis := GetModelAlias(field.Model)
					subForm := subAlis + modelAlias + strconv.FormatInt(int64(vb.ID), 10)
					subHiddenColumns := []string{}
					for _, sColumn := range field.Schema {
						if (sColumn.Hidden == true && sColumn.Default == nil && sColumn.Label == "") || (sColumn.Hidden == true && sColumn.Default == "" && sColumn.Label == "") {
							subHiddenColumns = append(subHiddenColumns, sColumn.Model)
						}
					}

					subColumnDataTypes := GetColumnsFromTableMeta(dbSchema.TableMeta[field.Model], subHiddenColumns)
					subStructs, _ := DBSchema.GenerateOnlyStruct(*subColumnDataTypes, field.Model, subForm, "", true, true, true, "", "")
					gormStructs = gormStructs + string(subStructs)
				}
			}
		}

		struc, _ := DBSchema.GenerateWithImports("", *columnDataTypes, schema.Model, modelAlias+strconv.FormatInt(int64(vb.ID), 10), "formModels", true, true, true, "", gormStructs, "")

		beforInsertTrigger := "nil"
		beforUpdateTrigger := "nil"
		afterInsertTrigger := "nil"
		afterUpdateTrigger := "nil"
		triggersNamespace := ""

		if schema.Triggers.Namespace != "" {
			if schema.Triggers.Insert.Before != "" {
				beforInsertTrigger = strings.ReplaceAll(schema.Triggers.Insert.Before, "@", ".")
			}
			if schema.Triggers.Update.Before != "" {
				beforUpdateTrigger = strings.ReplaceAll(schema.Triggers.Update.Before, "@", ".")
			}
			if schema.Triggers.Insert.After != "" {
				afterInsertTrigger = strings.ReplaceAll(schema.Triggers.Insert.After, "@", ".")
			}
			if schema.Triggers.Update.After != "" {
				afterUpdateTrigger = strings.ReplaceAll(schema.Triggers.Update.After, "@", ".")
			}

			triggersNamespace = "\"" + schema.Triggers.Namespace + "\""

		}

		formFields := createFieldTypes(schema)
		formulas := createFomulas(schema)
		rules, messages := createValidation(schema, *columnDataTypes)
		subForms, gridSubFroms := createSubForms(modelAliasWithID, schema)

		content := fmt.Sprintf(`package form

import (
    "github.com/lambda-platform/lambda/dataform"
    "github.com/lambda-platform/lambda/DB"
    "github.com/thedevsaddam/govalidator"
    "github.com/lambda-platform/lambda/models"
    "lambda/lambda/models/form/formModels"
    "time"
	%s
)

var _ = time.Time{}
var _ = DB.Date{}

func %sDataform() dataform.Dataform {
    return dataform.Dataform{
        Name:     "%s",
        Identity: "%s",
        Table:    "%s",
        Model:    new(formModels.%s),
        FieldTypes: %s,
        Formulas: %s,
        ValidationRules: %s,
        ValidationMessages: %s,
        SubForms:  %s,
        AfterInsert: %s,
        AfterUpdate: %s,
        BeforeInsert: %s,
        BeforeUpdate: %s,
        TriggerNameSpace: "%s",
    }
}
%s
`, triggersNamespace, modelAliasWithID, vb.Name, schema.Identity, schema.Model, modelAliasWithID, formFields, formulas, rules, messages, subForms, afterInsertTrigger, afterUpdateTrigger, beforInsertTrigger, beforUpdateTrigger, schema.Triggers.Namespace, gridSubFroms)
		utils.WriteFileFormat(string(struc), "lambda/models/form/formModels/"+modelAlias+strconv.FormatInt(int64(vb.ID), 10)+".go")
		utils.WriteFileFormat(content, "lambda/models/form/"+modelAlias+strconv.FormatInt(int64(vb.ID), 10)+".go")

	}

}
func createSubForms(modelAliasWithID string, schema lambdaModels.SCHEMA) (string, string) {
	gridSubFroms := ""
	subForms := `[]map[string]interface{}{`

	for _, field := range schema.Schema {
		if field.FormType == "SubForm" {
			if field.SubType == "Grid" {
				subAlis := GetModelAlias(field.Model)

				subForms = subForms + "\nmap[string]interface{}{"

				//subForm := subAlis+strconv.FormatUint(field.FormId, 10)
				subForm := subAlis + modelAliasWithID
				subForms = subForms + `
							"connection_field":"` + field.Parent + `",
							"tableTypeColumn":"` + field.TableTypeColumn + `",
							"tableTypeValue":"` + field.TableTypeValue + `",
							"table":"` + field.Model + `",
							"parentIdentity":"` + schema.Identity + `",
							"subIdentity":"` + field.Identity + `",
							"subForm":` + subForm + `Dataform(),
							"subFormArray":new([]formModels.` + subForm + `),
`

				subForms = subForms + `
},`

				content := fmt.Sprintf(`
func %sDataform() dataform.Dataform {
    return dataform.Dataform{
        Name:     "%s",
        Identity: "%s",
        Table:    "%s",
        Model:    new(formModels.%s),
    }
}
`, subForm, field.Name, field.Identity, field.Model, subForm)
				gridSubFroms = gridSubFroms + content
			} else {
				subAlis := GetModelAlias(field.Model)

				subForms = subForms + "\nmap[string]interface{}{"

				subForm := subAlis + strconv.FormatUint(field.FormId, 10)
				//subForm := subAlis+modelAlias+strconv.FormatUint(vb.ID, 10)
				subForms = subForms + `
							"connection_field":"` + field.Parent + `",
							"tableTypeColumn":"` + field.TableTypeColumn + `",
							"tableTypeValue":"` + field.TableTypeValue + `",
							"table":"` + field.Model + `",
							"parentIdentity":"` + schema.Identity + `",
							"subIdentity":"` + field.Identity + `",
							"subForm":` + subForm + `Dataform(),
							"subFormArray":new([]formModels.` + subForm + `),
`

				subForms = subForms + `
},`

			}

		}
	}
	subForms = subForms + `}`
	return subForms, gridSubFroms
}

func createFieldTypes(schema lambdaModels.SCHEMA) string {
	formFields := `map[string]string{
			`
	for i := range schema.Schema {
		formFields = formFields + `"` + schema.Schema[i].Model + `":"` + schema.Schema[i].FormType + `",
`
	}
	formFields = formFields + `
			}`
	return formFields
}

func createFomulas(schema lambdaModels.SCHEMA) string {
	formulas := `[]models.Formula{
			`
	if len(schema.Formula) >= 1 {

		for _, formula := range schema.Formula {
			targets := ""

			for _, target := range formula.Targets {
				targets = targets + fmt.Sprintf(`models.Target{
                        Field: "%s",
                        Prop: "%s",
                    },`, target.Field, target.Prop)
			}
			formulas = formulas + fmt.Sprintf(`models.Formula{
                Targets: []models.Target{
                   %s
                },
                Template: "%s",
                Model:"%s",
                Form: "%s",
            },`, targets, formula.Template, formula.Model, formula.Form)
		}

	}
	formulas = formulas + `
			}`
	return formulas
}

func createValidation(schema lambdaModels.SCHEMA, columnDataTypes map[string]map[string]string) (string, string) {

	rules := `govalidator.MapData{
		`
	messages := `govalidator.MapData{
		`

	for _, field := range schema.Schema {

		if len(field.Rules) >= 1 && schema.Identity != field.Model && field.Model != "created_at" && field.Model != "updated_at" && field.Model != "deleted_at" {
			fieldRules := ""
			fieldMessages := ""
			for _, rule := range field.Rules {

				if rule.Type == "required" {

					if columnDataTypes[field.Model]["nullable"] != "YES" {
						fieldRules = fieldRules + "\"" + rule.Type + "\","
					}

				} else {
					//if rule.Type != "unique" && rule.Type != "englishAlphabet" && rule.Type != "lambda-account" {
					if rule.Type != "unique" && rule.Type != "mongolianCyrillic" && rule.Type != "mongolianMobileNumber" && rule.Type != "englishAlphabet" && rule.Type != "lambda-account" {
						if rule.Type == "number" {
							fieldRules = fieldRules + "\"" + "numeric" + "\","
						} else {
							fieldRules = fieldRules + "\"" + rule.Type + "\","
						}

					}
					fieldMessages = fieldMessages + "\"" + rule.Type + ":" + rule.Msg + "\","
				}

			}

			rules = rules + "\n\"" + field.Model + "\": []string{" + fieldRules + "},"
			messages = messages + "\n\"" + field.Model + "\": []string{" + fieldMessages + "},"

		}
	}

	rules = rules + `}`
	messages = messages + `}`

	return rules, messages
}

func WriteModelCaller(dbSchema lambdaModels.DBSCHEMA, forms []genertarModels.ProjectSchemas, copyClienModels bool) {
	//return new(models.Naiz)

	content := ""
	content = content + "package caller\n"

	content = content + "import \"lambda/lambda/models/form\"\n"
	content = content + "import \"github.com/lambda-platform/lambda/dataform\"\n"

	content = content + "func GetMODEL(schema_id string) (dataform.Dataform) {\n\nswitch schema_id {\n"

	if copyClienModels {
		content = content + `
 case "crud_form":
return form.KrudDataform()

 case "notification_target_form":
return form.NotificationTargetsDataform()

 case "menu_form":
return form.MenuFormDataform()

 case "user_form":
return form.UserFormDataform()

 case "user_profile":
return form.UserProfile()

 case "user_password":
return form.UsersDataform()

`
	}

	for _, vb := range forms {
		var schema lambdaModels.SCHEMA

		json.Unmarshal([]byte(vb.Schema), &schema)

		modelAlias := GetModelAlias(schema.Model)

		content = content + "\n case \"" + strconv.FormatInt(int64(vb.ID), 10) + "\": \nreturn form." + modelAlias + strconv.FormatInt(int64(vb.ID), 10) + "Dataform()\n"

	}

	content = content + "\n} \nreturn dataform.Dataform{} \n\n}"

	utils.WriteFileFormat(content, "lambda/models/form/caller/modelCaller.go")
}
