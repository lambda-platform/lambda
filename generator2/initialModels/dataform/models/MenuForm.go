package form

import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/models"
	"github.com/thedevsaddam/govalidator"
	"lambda/lambda/models/form/formModels"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

func MenuFormDataform() dataform.Dataform {
	return dataform.Dataform{
		Name:     "Цэсний тохиргоо",
		Identity: "id",
		Table:    "vb_schemas",
		Model:    new(formModels.VbSchemas),
		FieldTypes: map[string]string{
			"id":         "Text",
			"name":       "Text",
			"schema":     "AdminMenu",
			"created_at": "Text",
			"type":       "Text",
			"updated_at": "Text",
		},
		Formulas: []models.Formula{},
		ValidationRules: govalidator.MapData{

			"name":   []string{"required"},
			"schema": []string{"required"}},
		ValidationMessages: govalidator.MapData{

			"name":   []string{"required:Талбарыг бөглөнө үү!"},
			"schema": []string{"required:Талбарыг бөглөнө үү!"}},
		SubForms:         []map[string]interface{}{},
		AfterInsert:      nil,
		AfterUpdate:      nil,
		BeforeInsert:     nil,
		BeforeUpdate:     nil,
		TriggerNameSpace: "",
	}
}
