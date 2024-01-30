package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/agent/models"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	krudModels "github.com/lambda-platform/lambda/krud/models"
	puzzleModels "github.com/lambda-platform/lambda/models"
	"github.com/lambda-platform/lambda/utils"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type LoginRequest struct {
	Login    string `json:"login" xml:"login" form:"login" query:"login"`
	Password string `json:"password" xml:"password" form:"password" query:"password"`
}

type Permissions struct {
	DefaultMenu string      `json:"default_menu"`
	Extra       interface{} `json:"extra"`
	MenuID      int         `json:"menu_id"`
	Permissions interface{} `json:"permissions"`
}

func Login(c *fiber.Ctx) error {

	request := LoginRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
			Error:  "Username & password required",
			Status: false,
		})
	}

	user, err := agentUtils.AuthUser(request.Login, "login")
	var roleID int64 = 0
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
			Error:  err.Error(),
			Status: false,
		})
	}

	if agentUtils.IsSame(request.Password, user["password"].(string)) {

		roleID = agentUtils.GetRole(user["role"])
		// create jwt token
		token, err := CreateJwtToken(user, roleID)
		if err != nil {
			//log.Println("Error Creating JWT token", err)
			return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
				Error:  "Unauthorized",
				Status: false,
			})
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "token"
		cookie.Path = "/"
		cookie.Value = token
		if config.Config.App.CookieSecure {
			cookie.Secure = true
		}

		cookie.Expires = time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl))

		delete(user, "password")

		c.Cookie(cookie)

		OAuth := withOAuth(request.Login, c)

		return c.Status(fiber.StatusOK).JSON(models.LoginData{
			Token:  token,
			Path:   checkRole(roleID),
			Status: true,
			Data:   user,
			OAuth:  OAuth,
		})

	}

	return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
		Error:  "Unauthorized",
		Status: false,
	})

}

func withOAuth(username string, c *fiber.Ctx) bool {
	OAuth := false

	// Get value from ReturnUri cookie
	value := c.Cookies("ReturnUri")

	// Check if value exists and is not empty
	if value != "" {
		// Set the LoggedInUserID cookie
		c.Cookie(&fiber.Cookie{
			Name:     "LoggedInUserID",
			Value:    username,
			Expires:  time.Now().Add(2 * time.Minute),
			HTTPOnly: true,
		})
		OAuth = true
	}

	return OAuth
}
func GetPermissions(c *fiber.Ctx) error {

	user, err := agentUtils.AuthUserObject(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
			Error:  err.Error(),
			Status: false,
		})
	}

	roleID := agentUtils.GetRole(user["role"])

	permissionData := PermissionData(roleID)
	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"status":     true,
		"permission": permissionData,
	})

}
func PermissionData(roleID int64) map[string]interface{} {

	if config.LambdaConfig.ProjectKey != "" && config.LambdaConfig.LambdaMainServicePath != "" {
		RoleData := map[string]interface{}{}

		jsonFile, err := os.Open("lambda/role_" + fmt.Sprintf("%v", roleID) + ".json")
		defer jsonFile.Close()
		if err != nil {

		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &RoleData)

		return map[string]interface{}{
			"menu":                 RoleData["menu"],
			"kruds":                RoleData["cruds"],
			"permissions":          RoleData["permissions"],
			"subCrudFormGrid":      RoleData["subCrudFormGrid"],
			"subCrudSection":       RoleData["subCrudSection"],
			"subCruds":             RoleData["subCruds"],
			"microserviceSettings": RoleData["microserviceSettings"],
		}
	} else {
		if config.Config.Database.Connection == "oracle" {
			Role := models.RoleOracle{}
			DB.DB.Where("ID = ?", roleID).Find(&Role)

			Permissions_ := Permissions{}
			json.Unmarshal([]byte(Role.Permissions), &Permissions_)

			Menu := puzzleModels.VBSchemaOracle{}
			DB.DB.Where("ID = ?", Permissions_.MenuID).Find(&Menu)

			MenuSchema := new(interface{})
			json.Unmarshal([]byte(Menu.Schema), &MenuSchema)
			Kruds := []krudModels.KrudOracle{}
			DB.DB.Find(&Kruds)
			return map[string]interface{}{
				"menu":        MenuSchema,
				"kruds":       Kruds,
				"permissions": Permissions_,
			}
		} else {
			Role := models.Role{}
			DB.DB.Where("id = ?", roleID).Find(&Role)

			Permissions_ := Permissions{}
			json.Unmarshal([]byte(Role.Permissions), &Permissions_)

			Menu := puzzleModels.VBSchema{}
			DB.DB.Where("id = ?", Permissions_.MenuID).Find(&Menu)

			MenuSchema := new(interface{})
			json.Unmarshal([]byte(Menu.Schema), &MenuSchema)
			Kruds := []krudModels.Krud{}
			DB.DB.Where("deleted_at IS NULL").Find(&Kruds)
			return map[string]interface{}{
				"menu":        MenuSchema,
				"kruds":       Kruds,
				"permissions": Permissions_,
			}
		}
	}

}
func Logout(c *fiber.Ctx) error {

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Path = "/"
	cookie.Value = ""
	cookie.Expires = time.Now()

	c.Cookie(cookie)
	return c.JSON(map[string]string{
		"status": "true",
		"data":   "",
		"path":   "auth/login",
		"token":  "",
	})

}

func LoginPage(c *fiber.Ctx) error {
	//csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	csrfToken := ""
	return c.Render("login", map[string]interface{}{
		"title":         config.LambdaConfig.Title,
		"favicon":       config.LambdaConfig.Favicon,
		"lambda_config": config.LambdaConfig,
		"mix":           utils.Mix,
		"csrfToken":     csrfToken,
	})
}
func LambdaConfig(c *fiber.Ctx) error {

	return c.JSON(config.LambdaConfig)
}

func CreateJwtToken(user map[string]interface{}, role int64) (string, error) {
	// Set custom claims
	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl)).Unix(),
	}

	for k, v := range user {
		if k != "role" && k != "password" && k != "bio" && k != "deleted_at" && k != "status" && k != "created_at" && k != "updated_at" && k != "avatar" && k != "gender" {
			claims[k] = v
		}
	}
	//for i := 0; i < len(config.LambdaConfig.UserDataFields); i++ {
	//	userField := config.LambdaConfig.UserDataFields[i]
	//
	//	if userField != "id" && userField != "login" && userField != "role" {
	//		claims[userField] = user[userField]
	//	}
	//}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config.JWT.Secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func checkRole(role int64) string {
	for _, r := range config.LambdaConfig.RoleRedirects {
		if role == r.RoleID {
			return r.URL
		}
	}

	RoleData := map[string]interface{}{}

	jsonFile, err := os.Open("lambda/role_" + fmt.Sprintf("%v", role) + ".json")

	defer jsonFile.Close()
	if err != nil {
		foundRole := models.Role{}
		DB.DB.Where("id = ?", role).First(&foundRole)
		if foundRole.Permissions != "" {

			Permissions := models.Permissions{}
			json.Unmarshal([]byte(foundRole.Permissions), &Permissions)
			if Permissions.DefaultMenu != "" {
				return config.LambdaConfig.AppURL + Permissions.DefaultMenu
			}
		}
	}
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &RoleData)

	permissonData, err := json.Marshal(RoleData["permissions"])
	if err != nil {
		return "/auth/login"
	}

	permissions := models.Permissions{}
	json.Unmarshal(permissonData, &permissions)
	return config.LambdaConfig.AppURL + permissions.DefaultMenu

}
