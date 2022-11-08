package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/exportImport/models"
	"os"
	"strings"
)

func Export(c *fiber.Ctx) error {
	idsPre := c.Query("ids")
	isMicroservice := c.Query("isMicroservice")
	schemaTable := "vb_schemas"
	krudTable := "krud"
	if isMicroservice == "true" {
		schemaTable = "project_schemas"
		krudTable = "project_cruds"
	}
	ids := strings.Split(idsPre, ",")

	exportData := models.LambdaExportData{}
	DB.DB.Table(krudTable).Where("id IN (?)", ids).Find(&exportData.Kruds)

	for index, _ := range exportData.Kruds {
		if exportData.Kruds[index].Form >= 1 {
			DB.DB.Table(schemaTable).Where("id = ?", exportData.Kruds[index].Form).Find(&exportData.Kruds[index].FormSchema)
		}
		if exportData.Kruds[index].Grid >= 1 {
			DB.DB.Table(schemaTable).Where("id = ?", exportData.Kruds[index].Grid).Find(&exportData.Kruds[index].GridSchema)
		}

	}

	byteData, err := json.Marshal(exportData)
	if err != nil {
		return c.JSON(map[string]string{
			"error": err.Error(),
		})
	}

	err = os.WriteFile("lambda/schemas/crud-export.json", byteData, 0755)
	if err != nil {
		return c.JSON(map[string]string{
			"error": err.Error(),
		})
	} else {
		return c.Download("lambda/schemas/crud-export.json")
	}

}
