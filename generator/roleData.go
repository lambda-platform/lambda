package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lambda-platform/lambda/DB"
	agentModels "github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	krudModels "github.com/lambda-platform/lambda/krud/models"
	"github.com/lambda-platform/lambda/models"
)

const roleDir = "lambda"

func RoleData() error {

	switch config.Config.Database.Connection {
	case "oracle":
		return generateRoleDataOracle()
	case "mssql":
		return generateRoleDataMSSQL()
	default:
		// Default: Postgres/MySQL
		return generateRoleDataDefault()
	}
}

func generateRoleDataOracle() error {
	// Krud + Role + Menu models for Oracle
	var kruds []krudModels.KrudOracle
	if err := DB.DB.Find(&kruds).Error; err != nil {
		return fmt.Errorf("load oracle kruds: %w", err)
	}

	var roles []agentModels.RoleOracle
	if err := DB.DB.Where("id >= ?", 1).Find(&roles).Error; err != nil {
		return fmt.Errorf("load oracle roles: %w", err)
	}

	roleData := make(map[int]map[string]interface{})

	for _, role := range roles {
		var perms models.Permissions
		if err := json.Unmarshal([]byte(role.Permissions), &perms); err != nil {
			// Skip broken record instead of panic
			continue
		}

		var menu models.VBSchemaOracle
		if err := DB.DB.Where("id = ?", perms.MenuID).First(&menu).Error; err != nil {
			// Хэрэв menu олдохгүй бол энэ role-г алгасна
			continue
		}

		var menuSchema interface{}
		if err := json.Unmarshal([]byte(menu.Schema), &menuSchema); err != nil {
			continue
		}

		roleData[role.ID] = map[string]interface{}{
			"permissions": perms,
			"menu":        menuSchema,
			"cruds":       kruds,
		}
	}

	return writeRoleFiles(roleData)
}

func generateRoleDataMSSQL() error {
	// Krud + Role + Menu models for MSSQL
	var kruds []krudModels.Krud
	if err := DB.DB.Find(&kruds).Error; err != nil {
		return fmt.Errorf("load mssql kruds: %w", err)
	}

	var roles []agentModels.Role
	if err := DB.DB.Where("id >= ?", 1).Find(&roles).Error; err != nil {
		return fmt.Errorf("load mssql roles: %w", err)
	}

	roleData := make(map[int]map[string]interface{})

	for _, role := range roles {
		var perms models.Permissions
		if err := json.Unmarshal([]byte(role.Permissions), &perms); err != nil {
			continue
		}

		var menu models.VBSchema
		if err := DB.DB.Where("id = ?", perms.MenuID).First(&menu).Error; err != nil {
			continue
		}

		var menuSchema interface{}
		if err := json.Unmarshal([]byte(menu.Schema), &menuSchema); err != nil {
			continue
		}

		roleData[role.ID] = map[string]interface{}{
			"permissions": perms,
			"menu":        menuSchema,
			"cruds":       kruds,
		}
	}

	return writeRoleFiles(roleData)
}

func generateRoleDataDefault() error {
	// Default connection: ихэвчлэн Postgres гэж төсөөлөөд generic Krud + generic VBSchema ашиглая.
	// Хэрэв танайд тусдаа Postgres model байвал энд солиорой.
	var kruds []krudModels.Krud
	if err := DB.DB.Find(&kruds).Error; err != nil {
		return fmt.Errorf("load default kruds: %w", err)
	}

	var roles []agentModels.Role
	if err := DB.DB.Where("id >= ?", 1).Find(&roles).Error; err != nil {
		return fmt.Errorf("load default roles: %w", err)
	}

	roleData := make(map[int]map[string]interface{})

	for _, role := range roles {
		var perms models.Permissions
		if err := json.Unmarshal([]byte(role.Permissions), &perms); err != nil {
			continue
		}

		// Хэрэв танай project-д VBSchemaPostgres гэдэг struct байгаа бол тэрийг ашигла.
		var menu models.VBSchema
		if err := DB.DB.Where("id = ?", perms.MenuID).First(&menu).Error; err != nil {
			continue
		}

		var menuSchema interface{}
		if err := json.Unmarshal([]byte(menu.Schema), &menuSchema); err != nil {
			continue
		}

		roleData[role.ID] = map[string]interface{}{
			"permissions": perms,
			"menu":        menuSchema,
			"cruds":       kruds,
		}
	}

	return writeRoleFiles(roleData)
}

func writeRoleFiles(roleData map[int]map[string]interface{}) error {
	for id, data := range roleData {
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal role %d: %w", id, err)
		}

		filename := filepath.Join(roleDir, "role_"+strconv.Itoa(id)+".json")
		if err := os.WriteFile(filename, b, 0600); err != nil {
			return fmt.Errorf("write %s: %w", filename, err)
		}
	}
	return nil
}
