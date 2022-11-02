package utils

import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	krudModels "github.com/lambda-platform/lambda/krud/models"
)

func AutoMigrateSeed() {

	if config.Config.Database.Connection == "oracle" {
		DB.DB.AutoMigrate(
			&krudModels.KrudOracle{},
			&krudModels.KrudTemplateOracle{},
		)

		if config.Config.App.Seed == "true" {
			var vbs []krudModels.KrudTemplateOracle
			DB.DB.Find(&vbs)

			if len(vbs) <= 0 {
				seedData()
			}
		}
	} else {
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

}
func seedData() {
	/*KRUD TEMPLATES*/
	templates := []string{"canvas", "drawer", "window", "popup", "edit", "create"}

	for _, template := range templates {
		if config.Config.Database.Connection == "oracle" {
			newTemplate := krudModels.KrudTemplateOracle{
				TemplateName: template,
			}

			DB.DB.Create(&newTemplate)

		} else {
			newTemplate := krudModels.KrudTemplate{
				TemplateName: template,
			}

			DB.DB.Create(&newTemplate)
		}
	}

}
