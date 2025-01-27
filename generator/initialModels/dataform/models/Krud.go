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

func KrudDataform() dataform.Dataform {
	return dataform.Dataform{
		Name:     "Маягт, хүснэгт нэгтгэл",
		Identity: "id",
		Table:    "krud",
		Model:    new(formModels.Krud),
		FieldTypes: map[string]string{
			"id":         "",
			"title":      "Text",
			"grid":       "Select",
			"form":       "Select",
			"actions":    "Textarea",
			"created_at": "",
			"updated_at": "",
			"template":   "Select",
		},
		Formulas: []models.Formula{},
		ValidationRules: govalidator.MapData{

			"title": []string{"required"}},
		ValidationMessages: govalidator.MapData{

			"title": []string{"required:Талбарыг бөглөнө үү!"}},
		SubForms:         []map[string]interface{}{},
		AfterInsert:      nil,
		AfterUpdate:      nil,
		BeforeInsert:     nil,
		BeforeUpdate:     nil,
		TriggerNameSpace: "",
		Relations: map[string]models.Relation{
			"grid": models.Relation{
				Table:              "vb_schemas",
				Key:                "id",
				Fields:             []string{"name"},
				FilterWithUser:     (*[]models.FilterWithUser)(nil),
				SortField:          "id",
				SortOrder:          "desc",
				ParentFieldOfForm:  "",
				ParentFieldOfTable: "",
				Filter:             "type = 'grid'",
			},
			"form": models.Relation{
				Table:              "vb_schemas",
				Key:                "id",
				Fields:             []string{"name"},
				FilterWithUser:     (*[]models.FilterWithUser)(nil),
				SortField:          "id",
				SortOrder:          "desc",
				ParentFieldOfForm:  "",
				ParentFieldOfTable: "",
				Filter:             "type = 'form'",
			},
			"krud_templates": models.Relation{
				Table:              "krud_templates",
				Key:                "template_name",
				Fields:             []string{"template_name"},
				FilterWithUser:     (*[]models.FilterWithUser)(nil),
				SortField:          "template_name",
				SortOrder:          "asc",
				ParentFieldOfForm:  "",
				ParentFieldOfTable: "",
				Filter:             "",
			},
		},
	}
}
