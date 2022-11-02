package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"

	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/exportImport/models"
	"net/http"
	"os"
)

func Import(c *fiber.Ctx) error {

	data := models.LambdaExportData{}

	isMicroservice := c.Query("isMicroservice")
	schemaTable := "vb_schemas"
	krudTable := "krud"
	if isMicroservice == "true" {
		schemaTable = "project_schemas"
		krudTable = "project_cruds"
	}

	file := c.Params("file")

	jsonFile, err := os.Open("lambda/" + file)
	defer jsonFile.Close()
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
	}
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)

	for index, _ := range data.Kruds {

		savedKrud := models.Krud{}
		DB.DB.Table(krudTable).Where("title = ?", &data.Kruds[index].Title).Find(&savedKrud)

		if savedKrud.ID < 1 {

			if data.Kruds[index].Form >= 1 {

				data.Kruds[index].FormSchema.ID = 0

				DB.DB.Table(schemaTable).Create(&data.Kruds[index].FormSchema)
				data.Kruds[index].Form = data.Kruds[index].FormSchema.ID
			}
			if data.Kruds[index].Grid >= 1 {
				data.Kruds[index].GridSchema.ID = 0

				DB.DB.Table(schemaTable).Create(&data.Kruds[index].GridSchema)
				data.Kruds[index].Grid = data.Kruds[index].GridSchema.ID
			}
			data.Kruds[index].ID = 0

			DB.DB.Table(krudTable).Create(&data.Kruds[index])

		}

	}

	return c.JSON(map[string]interface{}{
		"status":          true,
		"converted-cruds": len(data.Kruds),
	})
}
