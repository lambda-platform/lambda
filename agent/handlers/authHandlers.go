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

var failedAttempts = make(map[string]int)     // Нэвтрэлтийн алдааны тоолуур
var lockoutUntil = make(map[string]time.Time) // Блок хийх хугацаа
func Login(c *fiber.Ctx) error {

	request := LoginRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "Нэвтрэх нэр болон нууц үг шаардлагатай",
			"status": false,
		})
	}
	// Блок хугацааг шалгах
	if lockoutTime, exists := lockoutUntil[request.Login]; exists && time.Now().Before(lockoutTime) {
		lockMinute := int(lockoutTime.Sub(time.Now()).Minutes())

		if lockMinute == 0 {
			lockMinute = lockMinute + 1
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": false,
			"error":  fmt.Sprintf("Таны бүртгэл хаагдсан байна. %v минутын дараа дахин оролдоно уу.", lockMinute),
		})
	}

	// Хэрэглэгчийг шалгах
	user, err := agentUtils.AuthUser(request.Login, "login")
	if err != nil {
		failedAttempts[request.Login]++
		lockAccountIfNeeded(request.Login)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "Нэвтрэх нэр эсвэл нууц үг буруу байна",
			"status": false,
		})
	}

	var roleID int64 = 0

	if agentUtils.IsSame(request.Password, user["password"].(string)) {
		// Амжилттай нэвтэрсэн тохиолдолд алдааны тоолуурыг цэвэрлэх
		delete(failedAttempts, request.Login)
		delete(lockoutUntil, request.Login)
		delete(user, "password")
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
		cookie.HTTPOnly = true
		if !config.Config.JWT.DisableCookieSecure {
			cookie.Secure = true
		}
		if config.LambdaConfig.CookieDomain != "" {
			cookie.Domain = config.LambdaConfig.CookieDomain
		}

		cookie.Expires = time.Now().Add(time.Minute * time.Duration(config.Config.JWT.Ttl))

		c.Cookie(cookie)

		OAuth := withOAuth(request.Login, c)

		return c.Status(fiber.StatusOK).JSON(models.LoginData{
			Token:  token,
			Path:   checkRole(roleID),
			Status: true,
			Data:   user,
			OAuth:  OAuth,
		})

	} else {
		failedAttempts[request.Login]++
		lockAccountIfNeeded(request.Login)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":  "Нэвтрэх нэр эсвэл нууц үг буруу байна",
			"status": false,
		})
	}

}

// Блоклох нөхцөлүүдийг хэрэгжүүлэх функц
func lockAccountIfNeeded(username string) {
	if failedAttempts[username] == 3 {
		lockoutUntil[username] = time.Now().Add(10 * time.Minute) // 10 минут түгжинэ
	} else if failedAttempts[username] == 6 {
		lockoutUntil[username] = time.Now().Add(1 * time.Hour) // 1 цаг түгжинэ
	} else if failedAttempts[username] >= 9 {
		lockoutUntil[username] = time.Now().Add(24 * time.Hour) // 1 өдөр түгжинэ
	}
}
func CheckAuth(c *fiber.Ctx) error {
	user, err := agentUtils.AuthUserObject(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
			Error:  err.Error(),
			Status: false,
		})
	}

	return c.JSON(fiber.Map{
		"authenticated": true,
		"user":          user,
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
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Set expiration in the past

	// Optional: Set MaxAge to -1 to explicitly delete the cookie
	cookie.MaxAge = -1

	if !config.Config.JWT.DisableCookieSecure {
		cookie.Secure = true
	}
	if config.LambdaConfig.CookieDomain != "" {
		cookie.Domain = config.LambdaConfig.CookieDomain
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
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
		if k != "role" && k != "password" && k != "bio" && k != "deleted_at" && k != "status" && k != "created_at" && k != "updated_at" && k != "gender" {
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
