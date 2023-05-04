package datagrid

import "github.com/gofiber/fiber/v2"

type ImportExcelRequest struct {
	ExcelFile string `json:"excelFile"`
	SchemaID  int    `json:"schemaID"`
}

func ImportExcel(c *fiber.Ctx, datagrid Datagrid) error {

	if datagrid.IsExcelUpload {
		if datagrid.ExcelUploadCustomNamespace != "" {
			return datagrid.ExcelUploadCustomTrigger(datagrid, c)
		} else {
			request := ImportExcelRequest{}
			if err := c.BodyParser(&request); err != nil {
				return err
			}
		}
	}

	return nil
}
