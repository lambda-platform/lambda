package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"time"
)

func Crud(GetMODEL func(schema_id string) dataform.Dataform) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaId := c.Params("schemaId")
		action := c.Params("action")
		id := c.Params("id")

		return dataform.Exec(c, schemaId, action, id, GetMODEL)
	}
}

func Today(c *fiber.Ctx) error {
	currentTime := time.Now()
	return c.JSON(map[string]string{
		"today": currentTime.Format("2006-01-02"),
	})
}
func CheckUnique(c *fiber.Ctx) error {
	return dataform.CheckUnique(c)
}
func Upload(c *fiber.Ctx) error {

	return dataform.Upload(c)

}
func CheckCurrentPassword(c *fiber.Ctx) error {

	return utils.CheckCurrentPassword(c)
}

func UpdateRow(GetGridMODEL func(schema_id string) datagrid.Datagrid) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaId := c.Params("schemaId")

		return datagrid.Exec(c, schemaId, "update-row", "", GetGridMODEL)
	}
}

func Delete(GetGridMODEL func(schema_id string) datagrid.Datagrid) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaId := c.Params("schemaId")
		id := c.Params("id")

		return datagrid.Exec(c, schemaId, "delete", id, GetGridMODEL)
	}
}
func ExportExcel(GetGridMODEL func(schema_id string) datagrid.Datagrid) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaId := c.Params("schemaId")
		return datagrid.Exec(c, schemaId, "excel", "", GetGridMODEL)
	}
}
func Print(GetGridMODEL func(schema_id string) datagrid.Datagrid) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaId := c.Params("schemaId")
		return datagrid.Exec(c, schemaId, "print", "", GetGridMODEL)
	}
}

func dieIF(err error) {
	if err != nil {
		panic(err)
	}
}
