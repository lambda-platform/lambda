package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
		return claims, nil
	}
}

func AuthUser(value interface{}, uniqField string) (map[string]interface{}, error) {
	var userData map[string]interface{}

	table := "users"

	whereString := fmt.Sprintf("deleted_at IS NULL AND %s = ?", uniqField)

	if config.Config.Database.Connection == "oracle" {
		table = "USERS"
		uniqField = strings.ToUpper(uniqField)
		whereString = fmt.Sprintf("DELETED_AT IS NULL AND %s = ?", uniqField)
	}

	err := DB.DB.Table(table).Where(whereString, value).Find(&userData).Error

	if len(userData) >= 1 && err == nil {
		if config.Config.Database.Connection == "oracle" {
			userData = toLowerKeys(userData)
		}
		return userData, err
	} else {
		return userData, gorm.ErrRecordNotFound
	}

}
func toLowerKeys(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[strings.ToLower(k)] = v
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
