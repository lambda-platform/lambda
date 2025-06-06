package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func AuthUserOracle(c *fiber.Ctx) *models.USERSOracle {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	Id := claims["id"]

	User := models.USERSOracle{}

	DB.DB.Where("ID = ?", Id).First(&User)

	//User.Password = ""
	return &User
}

func AuthUserUUID(c *fiber.Ctx) *models.UserUUID {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	Id := claims["id"]

	User := models.UserUUID{}

	DB.DB.Where("id = ?", Id).First(&User)

	//User.Password = ""
	return &User
}
func AuthUserObject(c *fiber.Ctx) (map[string]interface{}, error) {

	if c.Locals("user") == nil {
		return map[string]interface{}{}, gorm.ErrRecordNotFound
	} else {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		if !config.Config.SysAdmin.UUID {

			claims["id"] = GetRole(claims["id"])

		}

		return claims, nil
	}
}

func AuthUser(value interface{}, uniqField string) (map[string]interface{}, error) {
	var userData map[string]interface{}

	table := "users"

	whereString := fmt.Sprintf("deleted_at IS NULL AND LOWER(%s) = ?", uniqField)

	if config.Config.Database.Connection == "oracle" {
		table = "USERS"
		uniqField = strings.ToUpper(uniqField)
		whereString = fmt.Sprintf("DELETED_AT IS NULL AND LOWER(%s) = ?", uniqField)
	}

	err := DB.DB.Table(table).Where(whereString, toLowerCase(value)).Find(&userData).Error

	if len(userData) >= 1 && err == nil {
		if config.Config.Database.Connection == "oracle" {
			userData = toLowerKeys(userData)
		}

		delete(userData, "updated_at")
		delete(userData, "created_at")
		delete(userData, "deleted_at")
		delete(userData, "bio")
		delete(userData, "status")
		delete(userData, "birthday")
		delete(userData, "register_number")
		delete(userData, "gender")

		return userData, err
	} else {
		return userData, gorm.ErrRecordNotFound
	}

}
func toLowerCase(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.ToLower(v)
	case float32:
		return strings.ToLower(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		return strings.ToLower(strconv.FormatFloat(v, 'f', -1, 64))
	case int:
		return strings.ToLower(strconv.Itoa(v))
	case int32:
		return strings.ToLower(strconv.Itoa(int(v)))
	case int64:
		return strings.ToLower(strconv.FormatInt(v, 10))
	default:
		return ""
	}
}

func toLowerKeys(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		switch v.(type) {
		case string:
			intV, parseError := strconv.ParseInt(v.(string), 10, 64)
			if parseError == nil {
				result[strings.ToLower(k)] = intV
			} else {
				result[strings.ToLower(k)] = v
			}
		default:
			result[strings.ToLower(k)] = v
		}

	}
	return result
}
func Hash(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashed), err
}
func IsSame(str string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str)) == nil
}

type passwordPost struct {
	Password string `json:"password"`
}

func AuthUserFromContext(c *fiber.Ctx) *models.User {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	Id := claims["id"]

	User := models.User{}

	DB.DB.Where("id = ?", Id).First(&User)

	//User.Password = ""
	return &User
}
func CheckCurrentPassword(c *fiber.Ctx) error {

	post := new(passwordPost)
	if err := c.BodyParser(post); err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"status": "false from json",
		})
	}

	if config.Config.SysAdmin.UUID {
		user := AuthUserUUID(c)

		if IsSame(post.Password, user.Password) {
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		} else {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"msg":    "Нууц үг буруу байна !!!",
			})

		}
	} else {
		if config.Config.Database.Connection == "oracle" {
			user := AuthUserOracle(c)

			if IsSame(post.Password, user.Password) {
				return c.JSON(map[string]interface{}{
					"status": true,
				})
			} else {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"status": false,
					"msg":    "Нууц үг буруу байна !!!",
				})

			}
		} else {
			user := AuthUserFromContext(c)

			if IsSame(post.Password, user.Password) {
				return c.JSON(map[string]interface{}{
					"status": true,
				})
			} else {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"status": false,
					"msg":    "Нууц үг буруу байна !!!",
				})

			}
		}

	}

}

func GetRole(role interface{}) int64 {

	roleDataType := reflect.TypeOf(role).String()
	var roleValue int64
	if roleDataType == "float64" {
		roleValue = int64(role.(float64))
	} else if roleDataType == "float32" {
		roleValue = int64(role.(float32))
	} else if roleDataType == "int" {
		roleValue = int64(role.(int))
	} else if roleDataType == "int32" {
		roleValue = int64(role.(int32))
	} else if roleDataType == "int64" {
		roleValue = role.(int64)
	} else if roleDataType == "string" {
		roleValue, _ = strconv.ParseInt(role.(string), 10, 64)
	} else {
		roleValue = int64(role.(int))
	}

	return roleValue
}
