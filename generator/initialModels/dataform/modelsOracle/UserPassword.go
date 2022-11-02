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

func UsersDataform() dataform.Dataform {
	return dataform.Dataform{
		Name:     "Нууц үг солих",
		Identity: "id",
		Table:    "users",
		Model:    new(formModels.Users),
		FieldTypes: map[string]string{
			"id":              "",
			"login":           "",
			"password":        "Password",
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
			"password": []string{"required"}},
		ValidationMessages: govalidator.MapData{
			"password": []string{"required:Талбарыг бөглөнө үү!"}},
		SubForms:         []map[string]interface{}{},
		AfterInsert:      nil,
		AfterUpdate:      nil,
		BeforeInsert:     nil,
		BeforeUpdate:     nil,
		TriggerNameSpace: "",
	}
}
