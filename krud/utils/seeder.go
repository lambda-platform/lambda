package utils

import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	krudModels "github.com/lambda-platform/lambda/krud/models"
)

func AutoMigrateSeed() {

	if config.Config.Database.Connection == "mssql" {
		DB.DB.AutoMigrate(
			&krudModels.Krud{},
			&krudModels.KrudTemplate{},
		)
	} else {
		DB.DB.AutoMigrate(
			&krudModels.Krud{},
			&krudModels.KrudTemplate{},
		)
	}

	if config.Config.App.Seed == "true" {
		var vbs []krudModels.KrudTemplate
		DB.DB.Find(&vbs)

		if len(vbs) <= 0 {
			seedData()
		}
	}
}
func seedData() {
	/*KRUD TEMPLATES*/
	templates := []string{"canvas", "drawer", "window", "popup", "edit", "create"}

	for _, template := range templates {
		newTemplate := krudModels.KrudTemplate{
			TemplateName: template,
		}

		DB.DB.Create(&newTemplate)
	}

}
