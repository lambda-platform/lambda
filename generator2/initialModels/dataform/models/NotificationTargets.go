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

func NotificationTargetsDataform() dataform.Dataform {
	return dataform.Dataform{
		Name:     "Зорилтод мэдэгдэл",
		Identity: "id",
		Table:    "notification_targets",
		Model:    new(formModels.NotificationTargets),
		FieldTypes: map[string]string{
			"id":             "Text",
			"link":           "Text",
			"title":          "Text",
			"body":           "Textarea",
			"schema_id":      "Select",
			"target_role":    "Select",
			"condition":      "Text",
			"target_actions": "Select",
		},
		Formulas: []models.Formula{},
		ValidationRules: govalidator.MapData{

			"link":        []string{"required"},
			"title":       []string{"required"},
			"body":        []string{"required"},
			"schema_id":   []string{"required"},
			"target_role": []string{"required"},
			"condition":   []string{"required"}},
		ValidationMessages: govalidator.MapData{

			"link":        []string{"required:Талбарыг бөглөнө үү!"},
			"title":       []string{"required:Талбарыг бөглөнө үү!"},
			"body":        []string{"required:Талбарыг бөглөнө үү!"},
			"schema_id":   []string{"required:Талбарыг бөглөнө үү!"},
			"target_role": []string{"required:Талбарыг бөглөнө үү!"},
			"condition":   []string{"required:Талбарыг бөглөнө үү!"}},
		SubForms:         []map[string]interface{}{},
		AfterInsert:      nil,
		AfterUpdate:      nil,
		BeforeInsert:     nil,
		BeforeUpdate:     nil,
		TriggerNameSpace: "",
	}
}
