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

func UserProfile() dataform.Dataform {
	return dataform.Dataform{
		Name:     "UserProfile",
		Identity: "id",
		Table:    "users",
		Model:    new(formModels.UserProfile),
		FieldTypes: map[string]string{
			"id":              "",
			"login":           "",
			"status":          "",
			"role":            "",
			"uuid":            "",
			"first_name":      "",
			"last_name":       "",
			"bio":             "",
			"birthday":        "",
			"gender":          "",
			"register_number": "",
			"email":           "",
			"phone":           "",
			"avatar":          "",
			"created_at":      "",
			"updated_at":      "",
			"deleted_at":      "",
		},
		Formulas: []models.Formula{},
		ValidationRules: govalidator.MapData{
			"email":           []string{"required"},
			"register_number": []string{"required"},
			"first_name":      []string{"required"},
			"last_name":       []string{"required"},
			"birthday":        []string{"required"},
			"gender":          []string{"required"},
		},
		ValidationMessages: govalidator.MapData{
			"email":           []string{"required:Талбарыг бөглөнө үү!", "unique:Давхацсан утга оруулсан байна!"},
			"register_number": []string{"required:Талбарыг бөглөнө үү!", "unique:Давхацсан утга оруулсан байна!"},
			"first_name":      []string{"required:Талбарыг бөглөнө үү!"},
			"last_name":       []string{"required:Талбарыг бөглөнө үү!"},
			"birthday":        []string{"required:Талбарыг бөглөнө үү!"},
			"gender":          []string{"required:Талбарыг бөглөнө үү!"},
		},
		SubForms:         []map[string]interface{}{},
		AfterInsert:      nil,
		AfterUpdate:      nil,
		BeforeInsert:     nil,
		BeforeUpdate:     nil,
		TriggerNameSpace: "",
	}
}
