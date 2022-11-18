package krudMW

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/agent/agentMW"
	agentModels "github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	"net/http"
)

type PermissionData struct {
	C      bool   `json:"c"`
	D      bool   `json:"d"`
	MenuID string `json:"menu_id"`
	R      bool   `json:"r"`
	Show   bool   `json:"show"`
	Title  string `json:"title"`
	U      bool   `json:"u"`
}
type Permissions struct {
	DefaultMenu string `json:"default_menu"`
	Extra       struct {
		Chart       bool `json:"chart"`
		Datasourcce bool `json:"datasourcce"`
		Datasource  bool `json:"datasource"`
		Moqup       bool `json:"moqup"`
	} `json:"extra"`
	MenuID      int                       `json:"menu_id"`
	Permissions map[string]PermissionData `json:"permissions"`
}

func PermissionEdit(c *fiber.Ctx) error {

	page_id := c.Query("page_id")
	action := c.Params("action")

	if page_id != "" {

		editPermission := GetPermission(c)

		if action == "edit" {
			if editPermission.R {
				return c.Next()
			} else {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  "Засах эрх олгогдоогүй байна",
					"status": false,
				})
			}
		}
		if action == "update" {
			if editPermission.U {
				return c.Next()
			} else {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  "Засах эрх олгогдоогүй байна",
					"status": false,
				})
			}
		}

	}
	return c.Next()

}
func PermissionCreate(c *fiber.Ctx) error {

	page_id := c.Query("page_id")
	if page_id != "" {
		editPermission := GetPermission(c)
		if editPermission.C {
			return c.Next()
		} else {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"error":  "Нэмэх эрх олгогдоогүй байна",
				"status": false,
			})
		}
	}
	return c.Next()

}
func PermissionDelete(c *fiber.Ctx) error {

	page_id := c.Query("page_id")
	if page_id != "" {
		editPermission := GetPermission(c)

		if editPermission.D {
			return c.Next()
		} else {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"error":  "Устгах эрх олгогдоогүй байна",
				"status": false,
			})
		}
	}
	return c.Next()

}

func GetPermission(c *fiber.Ctx) PermissionData {

	page_id := c.Query("page_id")
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := agentMW.GetUserRole(claims)
	if config.Config.Database.Connection == "oracle" {
		Role := agentModels.RoleOracle{}
		DB.DB.Where("ID = ?", role).Find(&Role)
		Permissions_ := Permissions{}
		json.Unmarshal([]byte(Role.Permissions), &Permissions_)

		return Permissions_.Permissions[page_id]
	} else {
		Role := agentModels.Role{}
		DB.DB.Where("id = ?", role).Find(&Role)
		Permissions_ := Permissions{}
		json.Unmarshal([]byte(Role.Permissions), &Permissions_)
		return Permissions_.Permissions[page_id]
	}

}
