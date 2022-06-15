package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/exportImport/models"
	"io/ioutil"
	"strings"

	"net/http"
)

func Export(c echo.Context) error {
	idsPre := c.QueryParam("ids")
	isMicroservice := c.QueryParam("isMicroservice")
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
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = ioutil.WriteFile("lambda/crud-export.json", byteData, 0755)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.Attachment("lambda/crud-export.json", "crud-export.json")
}
