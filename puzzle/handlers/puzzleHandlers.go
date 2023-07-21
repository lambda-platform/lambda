package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/DBSchema"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/datasource"
	"github.com/lambda-platform/lambda/models"
	"github.com/lambda-platform/lambda/utils"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type vb_schema struct {
	ID     int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

func Index(c *fiber.Ctx) error {

	dbSchema := models.DBSCHEMA{}

	if config.LambdaConfig.SchemaLoadMode == "auto" {

		dbSchema = DBSchema.GetDBSchema()

	} else {
		schemaFile, err := os.Open("lambda/db_schema.json")
		defer schemaFile.Close()
		if err != nil {
			fmt.Println("schema FILE NOT FOUND")
		}
		dbSchema = models.DBSCHEMA{}
		jsonParser := json.NewDecoder(schemaFile)
		jsonParser.Decode(&dbSchema)
	}

	var gridList []models.VBSchemaList
	var userRoles []models.UserRoles

	//DB.DB.Where("type = ?", "grid").Find(&gridList)
	DB.DB.Find(&userRoles)

	//gridList, err := models.VBSchemas(qm.Where("type = ?", "grid")).All(context.Background(), DB)
	//dieIF(err)

	User, err := agentUtils.AuthUserObject(c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  err.Error(),
			"status": false,
		})
	}
	csrfToken := ""
	return c.Render("puzzle", map[string]interface{}{
		"lambda_config":             config.LambdaConfig,
		"title":                     config.LambdaConfig.Title,
		"favicon":                   config.LambdaConfig.Favicon,
		"app_logo":                  config.LambdaConfig.Logo,
		"app_text":                  "СИСТЕМИЙН УДИРДЛАГА",
		"dbSchema":                  dbSchema,
		"gridList":                  gridList,
		"User":                      User,
		"user_fields":               config.LambdaConfig.UserDataFields,
		"user_roles":                userRoles,
		"data_form_custom_elements": config.LambdaConfig.DataFormCustomElements,
		"data_grid_custom_elements": config.LambdaConfig.DataGridCustomElements,
		"mix":                       utils.Mix,
		"csrfToken":                 csrfToken,
		"custom_icons":              config.LambdaConfig.CustomIcons,
	})

}

func GetTableSchema(c *fiber.Ctx) error {
	table := c.Params("table")
	tableMetas := DBSchema.TableMetas(table)
	return c.JSON(tableMetas)

}

func GetVB(c *fiber.Ctx) error {

	type_ := c.Params("type")
	id := c.Params("id")
	condition := c.Params("condition")

	if id != "" {

		match, _ := regexp.MatchString("_", id)

		if match {
			if config.Config.Database.Connection == "oracle" {
				VBSchema := models.VBSchemaAdminOracle{}

				DB.DB.Where("ID = ?", id).First(&VBSchema)

				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   VBSchema,
				})
			} else {
				VBSchema := models.VBSchemaAdmin{}

				DB.DB.Where("id = ?", id).First(&VBSchema)

				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   VBSchema,
				})
			}
		} else {
			if config.Config.Database.Connection == "oracle" {
				VBSchema := models.VBSchemaOracle{}
				if (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "form") || (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "grid") {

					schemaFile, err := os.Open("lambda/schemas/" + type_ + "/" + id + ".json")

					if err == nil {
						defer schemaFile.Close()
						byteValue, _ := io.ReadAll(schemaFile)
						VBSchema.Schema = string(byteValue)
						id_, _ := strconv.ParseUint(id, 0, 64)
						VBSchema.ID = int(id_)
					}

				} else {

					DB.DB.Where("ID = ?", id).First(&VBSchema)
					if type_ == "form" {
						if condition != "" {
							if condition != "builder" {
								return dataform.SetConditionOracle(condition, c, VBSchema)
							}
						}
					}
				}

				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   VBSchema,
				})
			} else {
				VBSchema := models.VBSchema{}

				if (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "form") || (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "grid") || (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "chart") {

					schemaFile, err := os.Open("lambda/schemas/" + type_ + "/" + id + ".json")

					if err == nil {
						defer schemaFile.Close()
						byteValue, _ := io.ReadAll(schemaFile)
						VBSchema.Schema = string(byteValue)
						id_, _ := strconv.ParseUint(id, 0, 64)
						VBSchema.ID = id_

						// TO DO !!!
						//if type_ == "form" {
						//	if condition != "" {
						//		if condition != "builder" {
						//
						//			return dataform.SetCondition(condition, c, VBSchema)
						//		}
						//	}
						//}
					} else {
						fmt.Println(err.Error())
					}

				} else {

					DB.DB.Where("id = ?", id).First(&VBSchema)
					if type_ == "form" {
						if condition != "" {
							if condition != "builder" {
								return dataform.SetCondition(condition, c, VBSchema)
							}
						}
					}
				}

				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   VBSchema,
				})
			}

		}

	} else {

		if config.Config.Database.Connection == "oracle" {
			VBSchemas := []models.VBSchemaListOracle{}

			DB.DB.Select("ID, NAME, TYPE, CREATED_AT, UPDATED_AT").Where("TYPE = ?", type_).Order("ID ASC").Find(&VBSchemas)

			return c.JSON(map[string]interface{}{
				"status": true,
				"data":   VBSchemas,
			})
		} else {
			VBSchemas := []models.VBSchemaList{}

			DB.DB.Select("id, name, type, created_at, updated_at").Where("type = ?", type_).Order("id ASC").Find(&VBSchemas)

			return c.JSON(map[string]interface{}{
				"status": true,
				"data":   VBSchemas,
			})
		}

	}

	return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
		"status": false,
	})

}
func GetMenuVB(c *fiber.Ctx) error {

	type_ := "menu"
	id := c.Params("id")

	VBSchema := models.VBSchema{}
	if (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "form") || (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "grid") || (config.LambdaConfig.LambdaMainServicePath != "" && config.LambdaConfig.ProjectKey != "" && type_ == "menu") {

		schemaFile, err := os.Open("lambda/schemas/" + type_ + "/" + id + ".json")

		if err == nil {
			defer schemaFile.Close()
			byteValue, _ := io.ReadAll(schemaFile)
			VBSchema.Schema = string(byteValue)
			id_, _ := strconv.ParseUint(id, 0, 64)
			VBSchema.ID = id_
		}

	} else {
		VBSchema := models.VBSchema{}
		DB.DB.Where("id = ?", id).First(&VBSchema)

	}

	return c.JSON(map[string]interface{}{
		"status": true,
		"data":   VBSchema,
	})

}
func SaveVB(modelName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type_ := c.Params("type")
		id := c.Params("id")
		//condition := c.Params("condition")

		vbs := new(vb_schema)
		if err := c.BodyParser(vbs); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		}

		if id != "" {
			id_, _ := strconv.ParseUint(id, 0, 64)
			var err error
			if config.Config.Database.Connection == "oracle" {
				vb := models.VBSchemaOracle{}

				DB.DB.Where("id = ?", id_).First(&vb)

				vb.Name = vbs.Name
				vb.Schema = vbs.Schema
				//_, err := vb.Update(context.Background(), DB, boil.Infer())

				BeforeSave(id_, type_)

				err = DB.DB.Save(&vb).Error

				if err != nil {

					return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
						"status": false,
						"error":  err.Error(),
					})
				} else {

					error := AfterSaveOracle(vb, type_)

					if error != nil {
						return c.JSON(map[string]interface{}{
							"status": false,
							"error":  error.Error(),
						})
					} else {
						return c.JSON(map[string]interface{}{
							"status": true,
						})
					}
				}
			} else {
				vb := models.VBSchema{}

				DB.DB.Where("id = ?", id_).First(&vb)

				vb.Name = vbs.Name
				vb.Schema = vbs.Schema
				//_, err := vb.Update(context.Background(), DB, boil.Infer())

				BeforeSave(id_, type_)

				err = DB.DB.Save(&vb).Error

				if err != nil {

					return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
						"status": false,
						"error":  err.Error(),
					})
				} else {

					error := AfterSave(vb, type_)

					if error != nil {
						return c.JSON(map[string]interface{}{
							"status": false,
							"error":  error.Error(),
						})
					} else {
						return c.JSON(map[string]interface{}{
							"status": true,
						})
					}
				}
			}

			if type_ == "form" {
				//WriteModelData(id_)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModel(modelName)
				//WriteGridModelById(modelName, vb.ID)
			}

		} else {
			vb := models.VBSchema{
				Name:   vbs.Name,
				Schema: vbs.Schema,
				Type:   type_,
				ID:     0,
			}

			//err := vb.Insert(context.Background(), DB, boil.Infer())

			err := DB.DB.Create(&vb).Error

			if type_ == "form" {
				//WriteModelData(vb.ID)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModelById(modelName, vb.ID)
				//WriteGridModel(modelName)
			}

			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(map[string]string{
					"status": "false",
				})
			} else {
				error := AfterSave(vb, type_)

				if error != nil {
					return c.JSON(map[string]interface{}{
						"status": false,
						"error":  error.Error(),
					})
				} else {
					return c.JSON(map[string]interface{}{
						"status": true,
					})
				}

			}

		}

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
		})
	}
}
func DeleteVB(c *fiber.Ctx) error {

	type_ := c.Params("type")
	id := c.Params("id")
	//condition := c.Params("condition")

	vbs := new(vb_schema)
	id_, _ := strconv.ParseUint(id, 0, 64)

	BeforeDelete(id_, type_)

	err := DB.DB.Where("id = ?", id).Where("type = ?", type_).Delete(&vbs).Error

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status": "false",
		})
	} else {
		return c.JSON(map[string]string{
			"status": "true",
		})
	}

}

func GetProjectVBs(c *fiber.Ctx) error {

	type_ := c.Params("type")
	id := c.Params("id")
	VBSchemas := []models.ProjectVBSchema{}

	if id != "" {
		VBSchema := models.VBSchema{}

		DB.DB.Table("project_schemas").Where("id = ?", id).First(&VBSchema)

		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   VBSchema,
		})

	} else {
		DB.DB.Table("project_schemas").Select("id, name, type, created_at, updated_at, projects_id").Where("type = ?", type_).Order("id ASC").Find(&VBSchemas)

		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   VBSchemas,
		})
	}

}
func GetProjectVB(c *fiber.Ctx) error {

	pid := c.Params("pid")
	type_ := c.Params("type")
	id := c.Params("id")
	condition := c.Params("condition")

	if id != "" {

		VBSchema := models.VBSchema{}

		DB.DB.Table("project_schemas").Where("id = ? AND projects_id = ?", id, pid).First(&VBSchema)

		if type_ == "form" {

			if condition != "" {
				if condition != "builder" {
					return dataform.SetCondition(condition, c, VBSchema)
				}
			}
		}

		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   VBSchema,
		})

	} else {
		VBSchemas := []models.VBSchemaList{}

		DB.DB.Table("project_schemas").Select("id, name, type, created_at, updated_at").Where("type = ? AND projects_id = ?", type_, pid).Order("id ASC").Find(&VBSchemas)

		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   VBSchemas,
		})
	}

	return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
		"status": false,
	})

}

func SaveProjectVB(modelName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid := c.Params("pid")
		ProjectID, _ := strconv.Atoi(pid)
		type_ := c.Params("type")
		id := c.Params("id")
		//condition := c.Params("condition")

		vbs := new(vb_schema)
		if err := c.BodyParser(vbs); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		}

		if id != "" {
			id_, _ := strconv.ParseUint(id, 0, 64)

			vb := models.ProjectVBSchema{}

			DB.DB.Where("id = ?", id_).First(&vb)

			vb.Name = vbs.Name
			vb.ProjectID = ProjectID
			vb.Schema = vbs.Schema
			//_, err := vb.Update(context.Background(), DB, boil.Infer())

			BeforeSave(id_, type_)

			err := DB.DB.Save(&vb).Error

			if type_ == "form" {
				//WriteModelData(id_)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModel(modelName)
				//WriteGridModelById(modelName, vb.ID)
			}

			if err != nil {

				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"status": false,
					"error":  err.Error(),
				})
			} else {

				//error := AfterSave(vb, type_)
				//
				//if(error != nil){
				//	return c.JSON(map[string]interface{}{
				//		"status": false,
				//		"error":error.Error(),
				//	})
				//} else {
				return c.JSON(map[string]interface{}{
					"status": true,
				})
				//}
			}

		} else {
			vb := models.ProjectVBSchema{
				Name:      vbs.Name,
				Schema:    vbs.Schema,
				Type:      type_,
				ProjectID: ProjectID,
				ID:        0,
			}

			//err := vb.Insert(context.Background(), DB, boil.Infer())

			err := DB.DB.Create(&vb).Error

			if type_ == "form" {
				//WriteModelData(vb.ID)
				//WriteModelData(modelName)
				//WriteModelDataById(modelName, vb.ID)
			} else if type_ == "grid" {
				//WriteGridModelById(modelName, vb.ID)
				//WriteGridModel(modelName)
			}

			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(map[string]string{
					"status": "false",
				})
			} else {
				//error := AfterSave(vb, type_)
				//
				//if(error != nil){
				//	return c.JSON(map[string]interface{}{
				//		"status": false,
				//		"error":error.Error(),
				//	})
				//} else {
				return c.JSON(map[string]interface{}{
					"status": true,
				})
				//}

			}

		}

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
		})
	}
}
func DeleteProjectVB(c *fiber.Ctx) error {

	pid := c.Params("pid")
	type_ := c.Params("type")
	id := c.Params("id")
	//condition := c.Params("condition")

	vbs := new(vb_schema)
	//id_, _ := strconv.ParseUint(id, 0, 64)
	//
	//BeforeDelete(id_, type_)

	err := DB.DB.Table("project_schemas").Where("id = ? AND projects_id = ? AND type = ?", id, pid, type_).Delete(&vbs).Error

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status": "false",
		})
	} else {
		return c.JSON(map[string]string{
			"status": "true",
		})
	}

}

func BeforeDelete(id uint64, type_ string) {

	if type_ == "datasource" {
		vb := models.VBSchema{}

		DB.DB.Where("id = ?", id).First(&vb)

		datasource.DeleteView("ds_" + vb.Name)
	}

}
func BeforeSave(id uint64, type_ string) {

	if type_ == "datasource" {
		vb := models.VBSchema{}

		DB.DB.Where("id = ?", id).First(&vb)

		datasource.DeleteView("ds_" + vb.Name)
	}

}
func AfterSave(vb models.VBSchema, type_ string) error {

	if type_ == "datasource" {
		return datasource.CreateView(vb.Name, vb.Schema)
	}

	return nil

}
func AfterSaveOracle(vb models.VBSchemaOracle, type_ string) error {

	if type_ == "datasource" {
		return datasource.CreateView(vb.Name, vb.Schema)
	}

	return nil

}

/*GRID*/

func GridVB(GetGridMODEL func(schema_id string) datagrid.Datagrid) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaId := c.Params("schemaId")
		action := c.Params("action")
		id := c.Params("id")

		return datagrid.Exec(c, schemaId, action, id, GetGridMODEL)
	}
}

/*FROM*/

func GetOptions(c *fiber.Ctx) error {

	r := new(dataform.Relations)
	if err := c.BodyParser(r); err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
	}
	var optionsData = map[string][]map[string]interface{}{}

	for table, relation := range r.Relations {
		data := dataform.OptionsData(relation, c)

		optionsData[table] = data

	}
	return c.JSON(optionsData)

}
