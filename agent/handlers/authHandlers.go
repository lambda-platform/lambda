package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/agent/models"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"time"
)

type User struct {
	Login    string `json:"login" xml:"login" form:"login" query:"login"`
	Password string `json:"password" xml:"password" form:"password" query:"password"`
}
type UserData struct {
	Id    int64
	Login string
	Role  int64
}
type UserUUIDData struct {
	Id    string
	Login string
	Role  int64
}
type jwtClaims struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Role  int64  `json:"role"`
	jwt.StandardClaims
}
type jwtUUIDClaims struct {
	Id    string `json:"id"`
	Login string `json:"login"`
	Role  int64  `json:"role"`
	jwt.StandardClaims
}

func Login(c *fiber.Ctx) error {

	u := new(User)
	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
			Error:  "Username & password required",
			Status: false,
		})
	}

	foundUser := agentUtils.AuthUserObjectByLogin(u.Login)

	if len(foundUser) == 0 {

		return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
			Error:  "User not found",
			Status: false,
		})

	}

	//password, err := Hash(u.Password)
	//password_check1 := IsSame(password, foundUser.Password)
	if agentUtils.IsSame(u.Password, foundUser["password"].(string)) {

		var roleID int64 = 0
		var userID int64 = 0
		var userUUID string = ""

		if reflect.TypeOf(foundUser["id"]).String() == "string" {
			if config.Config.SysAdmin.UUID {
				userUUID = foundUser["id"].(string)
			} else {
				i, err := strconv.ParseInt(foundUser["id"].(string), 10, 64)
				if err != nil {
					panic(err)
				}
				userID = i
			}

		} else {
			userID = foundUser["id"].(int64)
		}
		if reflect.TypeOf(foundUser["role"]).String() == "string" {
			i, err := strconv.ParseInt(foundUser["role"].(string), 10, 64)
			if err != nil {
				panic(err)
			}
			roleID = i
		} else {
			roleID = foundUser["role"].(int64)
		}

		if config.Config.SysAdmin.UUID {
			// create jwt token
			token, err := createUUIDJwtToken(UserUUIDData{Id: userUUID, Login: foundUser["login"].(string), Role: roleID})
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
			cookie.Expires = time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl))
			//cookie.HttpOnly = true
			//cookie.Secure = true

			delete(foundUser, "password")

			foundUser["jwt"] = token

			c.Cookie(cookie)

			return c.Status(fiber.StatusOK).JSON(models.LoginData{
				Token:  token,
				Path:   checkRole(roleID),
				Status: true,
				Data:   foundUser,
			})
		} else {
			// create jwt token
			token, err := createJwtToken(UserData{Id: userID, Login: foundUser["login"].(string), Role: roleID})
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
			cookie.Expires = time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl))

			delete(foundUser, "password")

			foundUser["jwt"] = token

			c.Cookie(cookie)

			return c.Status(fiber.StatusOK).JSON(models.LoginData{
				Token:  token,
				Path:   checkRole(roleID),
				Status: true,
				Data:   foundUser,
			})
		}

	}

	return c.Status(fiber.StatusUnauthorized).JSON(models.Unauthorized{
		Error:  "Unauthorized",
		Status: false,
	})

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

func createJwtToken(user UserData) (string, error) {
	// Set custom claims
	claims := jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl)).Unix(),
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config.JWT.Secret))
	if err != nil {
		return "", err
	}
	return t, nil
}
func createUUIDJwtToken(user UserUUIDData) (string, error) {
	// Set custom claims
	claims := jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl)).Unix(),
	}
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
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &RoleData)

	permissonData, err := json.Marshal(RoleData["permissions"])
	if err != nil {
		return "/auth/login"
	}

	permissions := models.Permissions{}
	json.Unmarshal(permissonData, &permissions)
	return config.LambdaConfig.AppURL + permissions.DefaultMenu

}
