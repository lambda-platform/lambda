package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	agentModels "github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	krudModels "github.com/lambda-platform/lambda/krud/models"
	"github.com/lambda-platform/lambda/models"
	"net/http"
)

func GetRolesMenus(c *fiber.Ctx) error {

	microserviceID := c.Params("microserviceID")
	if microserviceID != "" {
		roles := []agentModels.Role{}
		menus := []models.ProjectVBSchema{}
		kruds := []krudModels.ProjectCruds{}

		DB.DB.Where("id != 1 AND id != 3 AND id != 4 AND id != 2").Find(&roles)
		DB.DB.Find(&kruds)
		DB.DB.Where("type = 'menu'").Find(&menus)

		return c.JSON(map[string]interface{}{
			"status": true,
			"roles":  roles,
			"menus":  menus,
			"cruds":  kruds,
		})
	} else {
		if config.Config.Database.Connection == "oracle" {
			roles := []agentModels.RoleOracle{}
			menus := []models.VBSchemaOracle{}
			kruds := []krudModels.KrudOracle{}

			DB.DB.Where("ID != 1").Find(&roles)
			DB.DB.Find(&kruds)
			DB.DB.Where("TYPE = 'menu'").Find(&menus)

			return c.JSON(map[string]interface{}{
				"status": true,
				"roles":  roles,
				"menus":  menus,
				"cruds":  kruds,
			})
		} else {
			roles := []agentModels.Role{}
			menus := []models.VBSchema{}
			kruds := []krudModels.Krud{}

			DB.DB.Where("id != 1").Find(&roles)
			DB.DB.Find(&kruds)
			DB.DB.Where("type = 'menu'").Find(&menus)

			return c.JSON(map[string]interface{}{
				"status": true,
				"roles":  roles,
				"menus":  menus,
				"cruds":  kruds,
			})
		}
	}

}

type Role struct {
	ID          int                    `json:"id"`
	Permissions map[string]interface{} `json:"permissions"`
	Extra       map[string]interface{} `json:"extra"`
}
type RoleNew struct {
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

// TableName sets the insert table name for this struct type
func (v *Role) TableName() string {
	return "roles"
}

func SaveRole(c *fiber.Ctx) error {

	role := new(Role)
	if err := c.BodyParser(role); err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status": "false",
		})
	}

	Extra, _ := json.Marshal(role.Extra)
	Permissions, _ := json.Marshal(role.Permissions)

	if config.Config.Database.Connection == "oracle" {
		role_ := agentModels.RoleOracle{}

		DB.DB.Where("ID = ?", role.ID).First(&role_)

		role_.Extra = string(Extra)
		role_.Permissions = string(Permissions)

		err := DB.DB.Save(&role_).Error

		if err != nil {

			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		}
	} else {
		role_ := agentModels.Role{}

		DB.DB.Where("id = ?", role.ID).First(&role_)

		role_.Extra = string(Extra)
		role_.Permissions = string(Permissions)

		err := DB.DB.Save(&role_).Error

		if err != nil {

			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		}
	}
}

func CreateRole(c *fiber.Ctx) error {

	role_ := new(RoleNew)

	if err := c.BodyParser(role_); err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"errer":  err.Error(),
		})
	}

	if config.Config.Database.Connection == "oracle" {
		role := agentModels.RoleOracle{}
		role.Description = role_.Description
		role.DisplayName = role_.DisplayName
		role.Name = role_.Name

		err := DB.DB.Create(&role).Error
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		}
	} else {
		role := agentModels.Role{}
		role.Description = role_.Description
		role.DisplayName = role_.DisplayName
		role.Name = role_.Name

		err := DB.DB.Create(&role).Error
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		}
	}
}

func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	role_ := new(RoleNew)

	if err := c.BodyParser(role_); err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
	}

	role := agentModels.Role{}

	DB.DB.Where("id = ?", id).First(&role)
	role.Description = role_.Description
	role.DisplayName = role_.DisplayName
	role.Name = role_.Name

	err := DB.DB.Save(&role).Error
	if err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
	} else {
		return c.JSON(map[string]string{
			"status": "true",
		})
	}
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	role := new(agentModels.Role)

	err := DB.DB.Where("id = ?", id).Delete(&role).Error

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

func GetKrudFieldsConsole(c *fiber.Ctx) error {
	id := c.Params("id")
	krud := krudModels.ProjectCruds{}
	form := models.ProjectVBSchema{}
	grid := models.ProjectVBSchema{}

	DB.DB.Where("id = ?", id).Find(&krud)
	DB.DB.Where("id = ?", krud.Form).Find(&form)
	DB.DB.Where("id = ?", krud.Grid).Find(&grid)

	var schema models.SCHEMA
	var gridSchema models.SCHEMAGRID

	json.Unmarshal([]byte(form.Schema), &schema)
	json.Unmarshal([]byte(grid.Schema), &gridSchema)

	formFields := []string{}
	gridFields := []string{}

	for _, field := range schema.Schema {
		formFields = append(formFields, field.Model)
	}
	for _, field := range gridSchema.Schema {
		gridFields = append(gridFields, field.Model)
	}

	return c.JSON(map[string]interface{}{
		"status":           "true",
		"user_fields":      config.LambdaConfig.UserDataFields,
		"form_fields":      formFields,
		"grid_fields":      gridFields,
		"grid_fields_full": gridSchema.Schema,
	})

}
func GetKrudFields(c *fiber.Ctx) error {
	id := c.Params("id")

	if config.Config.Database.Connection == "oracle" {
		krud := krudModels.KrudOracle{}
		form := models.VBSchemaOracle{}
		grid := models.VBSchemaOracle{}

		DB.DB.Where("ID = ?", id).Find(&krud)
		DB.DB.Where("ID = ?", krud.Form).Find(&form)
		DB.DB.Where("ID = ?", krud.Grid).Find(&grid)

		var schema models.SCHEMA
		var gridSchema models.SCHEMAGRID

		json.Unmarshal([]byte(form.Schema), &schema)
		json.Unmarshal([]byte(grid.Schema), &gridSchema)

		formFields := []string{}
		gridFields := []string{}

		for _, field := range schema.Schema {
			formFields = append(formFields, field.Model)
		}
		for _, field := range gridSchema.Schema {
			gridFields = append(gridFields, field.Model)
		}

		return c.JSON(map[string]interface{}{
			"status":           "true",
			"user_fields":      config.LambdaConfig.UserDataFields,
			"form_fields":      formFields,
			"grid_fields":      gridFields,
			"grid_fields_full": gridSchema.Schema,
		})
	} else {
		krud := krudModels.Krud{}
		form := models.VBSchema{}
		grid := models.VBSchema{}

		DB.DB.Where("id = ?", id).Find(&krud)
		DB.DB.Where("id = ?", krud.Form).Find(&form)
		DB.DB.Where("id = ?", krud.Grid).Find(&grid)

		var schema models.SCHEMA
		var gridSchema models.SCHEMAGRID

		json.Unmarshal([]byte(form.Schema), &schema)
		json.Unmarshal([]byte(grid.Schema), &gridSchema)

		formFields := []string{}
		gridFields := []string{}

		for _, field := range schema.Schema {
			formFields = append(formFields, field.Model)
		}
		for _, field := range gridSchema.Schema {
			gridFields = append(gridFields, field.Model)
		}

		return c.JSON(map[string]interface{}{
			"status":           "true",
			"user_fields":      config.LambdaConfig.UserDataFields,
			"form_fields":      formFields,
			"grid_fields":      gridFields,
			"grid_fields_full": gridSchema.Schema,
		})
	}

}
