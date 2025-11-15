package krudMW

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/lambda-platform/lambda/agent/agentMW"

	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
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

// Crud represents a single CRUD entry
type Crud struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Template  string `json:"template"`
	Grid      int    `json:"grid"`
	Form      int    `json:"form"`
	Actions   string `json:"actions"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Condition represents a condition object in formCondition or gridCondition
type Condition struct {
	FormField string `json:"form_field"`
	UserField string `json:"user_field"`
	GridField string `json:"grid_field,omitempty"` // Optional, only in gridCondition
}

// PermissionObj represents a permission object
type PermissionObj struct {
	C                      bool        `json:"c"`
	D                      bool        `json:"d"`
	R                      bool        `json:"r"`
	U                      bool        `json:"u"`
	Show                   bool        `json:"show"`
	Title                  string      `json:"title"`
	MenuID                 string      `json:"menu_id"`
	GridDeleteConditionJS  interface{} `json:"gridDeleteConditionJS"`
	GridDeleteConditionSQL interface{} `json:"gridDeleteConditionSQL"`
	GridEditConditionJS    interface{} `json:"gridEditConditionJS"`
	GridEditConditionSQL   interface{} `json:"gridEditConditionSQL"`
	FormCondition          []Condition `json:"formCondition,omitempty"`
	GridCondition          []Condition `json:"gridCondition,omitempty"`
}

// Extra represents the extra permissions object
type Extra struct {
	Approve            bool `json:"approve"`
	Chart              bool `json:"chart"`
	Datasource         bool `json:"datasource"`
	Excelupload        bool `json:"excelupload"`
	Hascustomcreatebtn bool `json:"hascustomcreatebtn"`
	Moqup              bool `json:"moqup"`
	Userlist           bool `json:"userlist"`
}

// Permissions represents the permissions section
type Permissions struct {
	DefaultMenu string                   `json:"default_menu"`
	Extra       Extra                    `json:"extra"`
	MenuID      int                      `json:"menu_id"`
	Permissions map[string]PermissionObj `json:"permissions"`
}

// MenuItem represents a menu item (recursive via Children)
type MenuItem struct {
	Children []MenuItem  `json:"children"`
	Icon     *string     `json:"icon,omitempty"`
	ID       string      `json:"id"`
	Key      *string     `json:"key,omitempty"`
	LinkTo   string      `json:"link_to"`
	SVG      string      `json:"svg"`
	Title    *string     `json:"title,omitempty"`
	URL      interface{} `json:"url"` // Can be string, int, or null
}

// Response represents the top-level JSON structure
type RoleData struct {
	Cruds       []Crud      `json:"cruds"`
	Menu        []MenuItem  `json:"menu"`
	Permissions Permissions `json:"permissions"`
}

func PermissionEdit(GetPermissionHandler func(c *fiber.Ctx, vbType string) PermissionObj, ignoreList []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		pageID := c.Query("page_id")
		id := c.Params("id")
		action := c.Params("action")
		schemaId := c.Params("schemaId")

		if isIgnore(schemaId, ignoreList) {
			c.Next()
		}

		profileSchemaId := os.Getenv("PROFILE_FORM_ID")
		changePasswordSchemaId := os.Getenv("CHANGE_PASSWORD_FORM_ID")

		// Хэрэглэгч өөрийн profile / password-оо л засна
		if schemaId == "user_profile" || schemaId == "user_password" || (schemaId == profileSchemaId && profileSchemaId != "") || (schemaId == changePasswordSchemaId && changePasswordSchemaId != "") {
			userID, err := agentUtils.AuthUserIDString(c)

			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":  err.Error(),
					"status": false,
				})
			}

			if userID != id {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error":  "Засах эрхгүй байна ",
					"status": false,
				})
			}

			// Нууц үг өөрчлөх үед current_password шалгана
			if action == "update" {
				var requestData map[string]interface{}
				if err := c.BodyParser(&requestData); err != nil {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"error":  "Invalid request body",
					})
				}

				requestID, _ := agentUtils.ToStringID(requestData["id"])
				if requestID != "" && requestID != userID {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{
						"error":  "Засах эрхгүй байна",
						"status": false,
					})
				}

				if schemaId == "user_password" || (schemaId == changePasswordSchemaId && changePasswordSchemaId != "") {
					currentPassword := ""
					if val, ok := requestData["current_password"]; ok {
						switch v := val.(type) {
						case string:
							currentPassword = v
						case []byte:
							currentPassword = string(v)
						case nil:
							return c.Status(http.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"error":  "Password cannot be nil",
							})
						default:
							currentPassword = fmt.Sprintf("%v", v)
						}
					}

					if currentPassword == "" {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{
							"status": false,
							"error":  "Password field required",
						})
					}

					var user struct{ Password string }

					if config.Config.SysAdmin.UUID {
						u := agentUtils.AuthUserUUID(c)
						user.Password = u.Password
					} else if config.Config.Database.Connection == "oracle" {
						u := agentUtils.AuthUserOracle(c)
						user.Password = u.Password
					} else {
						u := agentUtils.AuthUserFromContext(c)
						user.Password = u.Password
					}

					if !agentUtils.IsSame(currentPassword, user.Password) {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{
							"status": false,
							"msg":    "Нууц үг буруу байна !!!",
						})
					}
				}

				return c.Next()
			}

			if action == "edit" {
				return c.Next()
			}
		}

		// page_id дээр permission шалгах (form context)
		if pageID != "" {

			if GetPermissionHandler != nil {
				perm := GetPermissionHandler(c, "form")
				if (action == "edit" && perm.R) || (action == "update" && perm.U) {
					return c.Next()
				}
			} else {
				perm := GetPermission(c, "form")
				if (action == "edit" && perm.R) || (action == "update" && perm.U) {
					return c.Next()
				}
			}
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":  "Засах эрхгүй байна",
			"status": false,
		})
	}

}

func PermissionCreate(GetPermissionHandler func(c *fiber.Ctx, vbType string) PermissionObj, ignoreList []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		pageID := c.Query("page_id")
		action := c.Params("action")

		schemaId := c.Params("schemaId")
		if isIgnore(schemaId, ignoreList) {
			c.Next()
		}

		if action == "options" {
			return c.Next()
		}

		if pageID != "" {
			if GetPermissionHandler != nil {
				perm := GetPermissionHandler(c, "form")
				if perm.C {
					return c.Next()
				}
			} else {
				perm := GetPermission(c, "form")
				if perm.C {
					return c.Next()
				}
			}

		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":  "Нэмэх эрх олгогдоогүй байна",
			"status": false,
		})
	}
}

func PermissionDelete(GetPermissionHandler func(c *fiber.Ctx, vbType string) PermissionObj, ignoreList []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		pageID := c.Query("page_id")
		action := c.Params("filter-options")

		schemaId := c.Params("schemaId")
		if isIgnore(schemaId, ignoreList) {
			c.Next()
		}
		if action == "filter-options" {
			return c.Next()
		}

		if pageID != "" {
			if GetPermissionHandler != nil {
				perm := GetPermissionHandler(c, "grid")
				if perm.D {
					return c.Next()
				}
			} else {
				perm := GetPermission(c, "grid")
				if perm.D {
					return c.Next()
				}
			}
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":  "Устгах эрх олгогдоогүй байна",
			"status": false,
		})
	}
}

func PermissionRead(GetPermissionHandler func(c *fiber.Ctx, vbType string) PermissionObj, ignoreList []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		pageID := c.Query("page_id")
		schemaId := c.Params("schemaId")
		if isIgnore(schemaId, ignoreList) {
			c.Next()
		}
		if pageID != "" {
			if GetPermissionHandler != nil {
				perm := GetPermissionHandler(c, "grid")
				if perm.R || perm.Show {
					return c.Next()
				}
			} else {
				perm := GetPermission(c, "grid")
				if perm.R || perm.Show {
					return c.Next()
				}
			}
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":  "Унших эрхгүй байна",
			"status": false,
		})
	}
}

// ---------------- GetPermission core logic ----------------

func GetPermission(c *fiber.Ctx, vbType string) PermissionObj {
	pageID := c.Query("page_id")
	schemaIdSTr := c.Params("schemaId")
	schemaID, err := strconv.Atoi(schemaIdSTr)

	//fmt.Println(schemaID)
	//fmt.Println(pageID)
	//fmt.Println(vbType)

	if err != nil {
		return PermissionObj{}
	}
	// JWT-с role авах
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return PermissionObj{}
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return PermissionObj{}
	}
	role := agentMW.GetUserRole(claims)

	//fmt.Println(role)
	// role JSON файл унших
	jsonFile, err := os.Open("lambda/role_" + fmt.Sprintf("%v", role) + ".json")
	if err != nil {
		if config.Config.Database.Debug {
			fmt.Println("failed to open role file:", err)
		}
		return PermissionObj{}
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		if config.Config.Database.Debug {
			fmt.Println("failed to read role file:", err)
		}
		return PermissionObj{}
	}

	var roleData RoleData
	if err := json.Unmarshal(byteValue, &roleData); err != nil {
		if config.Config.Database.Debug {
			fmt.Println("failed to unmarshal role file:", err)
		}
		return PermissionObj{}
	}

	perms := roleData.Permissions.Permissions

	// 1. page_id дээрх base permission хайна
	if basePerm, ok := perms[pageID]; ok {
		// 1.1 menu_id-гаас menu олно
		menuID := basePerm.MenuID
		if menuID == "" {
			menuID = pageID
		}

		if menu := FindMenuByID(roleData.Menu, menuID); menu != nil {
			// 1.2 menu.URL-с CRUD олно
			if crud := FindCrudByMenuURL(roleData.Cruds, menu.URL); crud != nil {

				realPermission := PermissionObj{}

				if vbType == "form" {
					if crud.Form == schemaID {
						realPermission.C = basePerm.C
						realPermission.R = basePerm.R
						realPermission.U = basePerm.U
					}
				}
				if vbType == "grid" {
					if crud.Grid == schemaID {
						realPermission.R = basePerm.R
						realPermission.D = basePerm.D
					}
				}

				return realPermission

			}
		}

	}

	// Олдсонгүй
	return PermissionObj{}
}

// ---------------- Helper функцууд ----------------

func FindMenuByID(menu []MenuItem, id string) *MenuItem {
	for i := range menu {
		if menu[i].ID == id {
			return &menu[i]
		}
		if len(menu[i].Children) > 0 {
			if found := FindMenuByID(menu[i].Children, id); found != nil {
				return found
			}
		}
	}
	return nil
}

// menu.URL-ээс crud хайх
func FindCrudByMenuURL(cruds []Crud, url interface{}) *Crud {
	if url == nil {
		return nil
	}

	var raw string
	switch v := url.(type) {
	case string:
		raw = v
	case float64:
		raw = fmt.Sprintf("%.0f", v)
	case int:
		raw = fmt.Sprintf("%d", v)
	default:
		raw = fmt.Sprintf("%v", v)
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	// Хэрэв бүхэл тоо бол шууд CRUD ID гэж үзнэ
	if id, err := strconv.Atoi(raw); err == nil {
		for i := range cruds {
			if cruds[i].ID == id {
				return &cruds[i]
			}
		}
	}

	// URL path бол хамгийн сүүлийн segment-ийг ID гэж үзнэ (/vb/12 гэх мэт)
	parts := strings.Split(strings.Trim(raw, "/"), "/")
	last := parts[len(parts)-1]
	if id, err := strconv.Atoi(last); err == nil {
		for i := range cruds {
			if cruds[i].ID == id {
				return &cruds[i]
			}
		}
	}

	return nil
}
func isIgnore(item string, ignoreList []string) bool {
	for _, a := range ignoreList {
		if a == item {
			return true
		}
	}
	return false
}
