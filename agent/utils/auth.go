package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/agent/agentMW"
	"github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
func AuthUser(c *fiber.Ctx) *models.User {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	Id := claims["id"]

	User := models.User{}

	DB.DB.Where("id = ?", Id).First(&User)

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
func AuthUserObject(c *fiber.Ctx) map[string]interface{} {

	if c.Locals("user") == nil {
		return map[string]interface{}{}
	} else {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		Id := claims["id"]

		query := ""
		if config.Config.SysAdmin.UUID {
			query = fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", Id.(string))

		} else {
			userQuery := "SELECT * FROM users  WHERE id = '%d'"

			if config.Config.Database.Connection == "oracle" {
				userQuery = "SELECT * FROM \"USERS\" WHERE \"ID\" = '%d'"
			}

			query = fmt.Sprintf(userQuery, int(Id.(float64)))

		}

		rows, _ := DB.DB.Raw(query).Rows()

		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		userData := map[string]interface{}{}
		result_id := 0
		for rows.Next() {
			for i, _ := range columns {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)

			for i, col := range columns {

				val := values[i]

				if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
					//if col == "id"{
					//	if config.Config.SysAdmin.UUID {
					//		b, ok := val.([]byte)
					//		if ok {
					//			stringValue := string(b)
					//			userData[col] = stringValue
					//		} else {
					//			userData[col] = val
					//		}
					//	} else {
					//		userData[col] = val
					//	}
					//} else {
					//	userData[col] = val
					//}

					b, ok := val.([]byte)
					if ok {
						v, err := strconv.ParseInt(string(b), 10, 64)
						if err != nil {
							stringValue := string(b)
							//	fmt.Println(stringValue)

							userData[col] = stringValue
						} else {
							userData[col] = v
						}

					} else {
						userData[col] = val
					}
				} else {
					if config.Config.Database.Connection == "oracle" {
						col = strings.ToLower(col)
					}
					b, ok := val.([]byte)
					if ok {
						v, err := strconv.ParseInt(string(b), 10, 64)
						if err != nil {
							stringValue := string(b)
							//	fmt.Println(stringValue)

							userData[col] = stringValue
						} else {
							userData[col] = v
						}

					} else {
						userData[col] = val
					}
				}

			}

			result_id++
		}

		delete(userData, "password")

		userData["role"] = agentMW.GetUserRole(claims)
		return userData
	}

}

func AuthUserObjectByLogin(login string) map[string]interface{} {
	userData := map[string]interface{}{}

	userQuery := "SELECT * FROM users WHERE login = '%s'"

	if config.Config.Database.Connection == "oracle" {
		userQuery = "SELECT * FROM \"USERS\" WHERE \"LOGIN\" = '%s'"
	}
	rows, errorDB := DB.DB.Raw(fmt.Sprintf(userQuery, login)).Rows()

	//fmt.Println(login)
	fmt.Println(errorDB)

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	result_id := 0
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		for i, col := range columns {

			val := values[i]

			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {

				if col == "id" {
					if config.Config.SysAdmin.UUID {
						b, ok := val.([]byte)
						if ok {
							stringValue := string(b)
							userData[col] = stringValue
						} else {
							userData[col] = val
						}
					} else {
						userData[col] = val
					}

				} else {
					userData[col] = val
				}

			} else {
				if config.Config.Database.Connection == "oracle" {
					col = strings.ToLower(col)
				}
				b, ok := val.([]byte)
				if ok {
					v, err := strconv.ParseInt(string(b), 10, 64)
					if err != nil {
						stringValue := string(b)
						//	fmt.Println(stringValue)

						userData[col] = stringValue
					} else {
						userData[col] = v
					}

				} else {
					userData[col] = val
				}
			}

		}

		result_id++
	}

	return userData
}
func AuthUserObjectByEmail(login string) map[string]interface{} {
	userQuery := "SELECT * FROM users WHERE email = '%s'"

	if config.Config.Database.Connection == "oracle" {
		userQuery = "SELECT * FROM \"USERS\" WHERE \"EMAIL\" = '%s'"

	}
	rows, _ := DB.DB.Raw(fmt.Sprintf(userQuery, login)).Rows()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	userData := map[string]interface{}{}
	result_id := 0
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		for i, col := range columns {

			val := values[i]

			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
				userData[col] = val

			} else {
				b, ok := val.([]byte)
				if ok {
					v, err := strconv.ParseInt(string(b), 10, 64)
					if err != nil {
						stringValue := string(b)
						//	fmt.Println(stringValue)

						userData[col] = stringValue
					} else {
						userData[col] = v
					}

				} else {
					userData[col] = val
				}
			}

		}

		result_id++
	}

	return userData
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
			user := AuthUser(c)

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
