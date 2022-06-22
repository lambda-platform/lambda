package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func AuthUser(c echo.Context) *models.User {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	Id := claims["id"]

	User := models.User{}

	DB.DB.Where("id = ?", Id).First(&User)

	//User.Password = ""
	return &User
}
func AuthUserUUID(c echo.Context) *models.UserUUID {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	Id := claims["id"]

	User := models.UserUUID{}

	DB.DB.Where("id = ?", Id).First(&User)

	//User.Password = ""
	return &User
}
func AuthUserObject(c echo.Context) map[string]interface{} {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	Id := claims["id"]

	query := ""
	if config.Config.SysAdmin.UUID {
		query = fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", Id.(string))

	} else {
		query = fmt.Sprintf("SELECT * FROM users WHERE id = %d", int(Id.(float64)))
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

	return userData
}

func AuthUserObjectByLogin(login string) map[string]interface{} {
	userData := map[string]interface{}{}

	rows, errorDB := DB.DB.Raw(fmt.Sprintf("SELECT * FROM users WHERE login = '%s'", login)).Rows()

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

	rows, _ := DB.DB.Raw(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", login)).Rows()

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

func CheckCurrentPassword(c echo.Context) error {

	post := new(passwordPost)
	if err := c.Bind(post); err != nil {

		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "false from json",
		})
	}

	if config.Config.SysAdmin.UUID {
		user := AuthUserUUID(c)

		if IsSame(post.Password, user.Password) {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "true",
			})
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "false",
				"msg":    "Нууц үг буруу байна !!!",
			})

		}
	} else {
		user := AuthUser(c)

		if IsSame(post.Password, user.Password) {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "true",
			})
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "false",
				"msg":    "Нууц үг буруу байна !!!",
			})

		}
	}

}
