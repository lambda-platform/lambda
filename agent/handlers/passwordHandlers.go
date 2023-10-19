package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/agent/models"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils/mailer"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type mailPost struct {
	Email string `json:"email" `
	Lang  string `json:"lang"`
}
type passwordResetPost struct {
	Email           string `json:"email" `
	Code            string `json:"code"`
	Lang            string `json:"lang"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func SendForgotMail(c *fiber.Ctx) error {

	data := new(mailPost)
	if err := c.BodyParser(data); err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
			"msg":    "Post data error ",
		})
	}

	if data.Lang == "" {
		data.Lang = "mn_MN"
	}

	StaticWords := reflect.ValueOf(config.LambdaConfig.StaticWords[data.Lang]).Interface().(map[string]interface{})
	if data.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  StaticWords["emailRequired"],
			"msg":    StaticWords["emailRequired"],
		})
	}

	user, err := agentUtils.AuthUser(data.Email, "email")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"error":  StaticWords["userNotFound"],
			"msg":    StaticWords["userNotFound"],
			"status": false,
		})
	}

	pReset := models.PasswordReset{}
	DB.DB.Where("email = ?", data.Email).Delete(pReset)

	permittedChars := strings.Join(shuffle([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}), "")
	//permittedChars := strings.Join(shuffle([]string{"0","1","2","3","4","5","6","7","8","9","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}),"")
	tokenPre := string([]rune(permittedChars)[0:6])
	token, _ := agentUtils.Hash(tokenPre)

	pReset.Email = data.Email
	pReset.Token = token
	pReset.CreatedAt = time.Now()

	DB.DB.Create(&pReset)

	mail := mailer.NewRequest([]string{data.Email}, StaticWords["passwordResetCode"].(string), config.Config.Mail.FromAddress)
	mailSent := mail.Send("views/forgot.html", map[string]string{
		"keyword":           tokenPre,
		"passwordReset":     StaticWords["passwordReset"].(string),
		"passwordResetCode": StaticWords["passwordResetCode"].(string),
		"title":             config.LambdaConfig.Title,
		"noReply":           StaticWords["noReply"].(string),
		"copyright":         config.LambdaConfig.Copyright,
	})

	if mailSent {

		delete(user, "password")

		return c.JSON(map[string]interface{}{
			"msg":    StaticWords["passwordResetCodeSent"],
			"status": true,
			"data":   user,
		})
	} else {
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  "Mail error please contact system administrator",
			"msg":    "Mail error please contact system administrator",
		})
	}

}
func PasswordReset(c *fiber.Ctx) error {

	data := new(passwordResetPost)
	if err := c.BodyParser(data); err != nil {

		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
			"msg":    "Post data error ",
		})
	}

	if data.Lang == "" {
		data.Lang = "mn_MN"
	}

	StaticWords := reflect.ValueOf(config.LambdaConfig.StaticWords[data.Lang]).Interface().(map[string]interface{})
	if data.Email == "" || data.Code == "" || data.Password == "" || data.PasswordConfirm == "" {
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  StaticWords["passwordResetCodeRequired"],
			"msg":    StaticWords["passwordResetCodeRequired"],
		})
	}

	if config.Config.SysAdmin.UUID {
		foundUser := models.UserUUID{}
		pReset := models.PasswordReset{}
		PasswordResetTimeOut := config.LambdaConfig.PasswordResetTimeOut

		errU := DB.DB.Where("email = ?", data.Email).First(&foundUser).Error
		errR := DB.DB.Where("email = ?", data.Email).First(&pReset).Error

		if errU != nil || foundUser.Login == "" {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"error":  StaticWords["userNotFound"],
				"msg":    StaticWords["userNotFound"],
				"status": false,
			})
		}
		if errR != nil || pReset.Email == "" {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"error":  StaticWords["userNotFound"],
				"msg":    StaticWords["userNotFound"],
				"status": false,
			})
		}

		now := time.Now()
		diff := now.Sub(pReset.CreatedAt)

		mins := int(diff.Minutes() * 1)

		if PasswordResetTimeOut <= 1 {
			PasswordResetTimeOut = 3
		}
		if PasswordResetTimeOut >= mins && pReset.Wrong <= 2 {

			if agentUtils.IsSame(data.Code, pReset.Token) {

				if data.Password == data.PasswordConfirm {

					newPassword, _ := agentUtils.Hash(data.Password)

					foundUser.Password = newPassword
					err := DB.DB.Save(foundUser).Error
					if err != nil {

						return c.JSON(map[string]interface{}{
							"status": false,
						})
					} else {
						DB.DB.Where("email = ?", data.Email).Delete(pReset)
						return c.JSON(map[string]interface{}{
							"msg":    StaticWords["passwordResetSuccess"],
							"status": true,
						})
					}

				} else {

					return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
						"error":  StaticWords["passwordConfirmError"],
						"msg":    StaticWords["passwordConfirmError"],
						"status": false,
					})
				}

			} else {
				pReset.Wrong = pReset.Wrong + 1
				DB.DB.Save(pReset)
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  StaticWords["passwordResetCodeIncorrect"],
					"msg":    StaticWords["passwordResetCodeIncorrect"],
					"status": false,
				})
			}
		} else {
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"error":  StaticWords["passwordResetCodeTimeout"],
				"msg":    StaticWords["passwordResetCodeTimeout"],
				"status": false,
			})
		}
	} else {
		if config.Config.Database.Connection == "oracle" {
			foundUser := models.USERSOracle{}
			pReset := models.PASSWORDRESETSOracle{}
			PasswordResetTimeOut := config.LambdaConfig.PasswordResetTimeOut

			errU := DB.DB.Where("EMAIL = ?", data.Email).Find(&foundUser).Error
			errR := DB.DB.Where("EMAIL = ?", data.Email).Find(&pReset).Error

			if errU != nil || foundUser.Login == "" {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  StaticWords["userNotFound"],
					"msg":    StaticWords["userNotFound"],
					"status": false,
				})
			}
			if errR != nil || pReset.Email == "" {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  StaticWords["userNotFound"],
					"msg":    StaticWords["userNotFound"],
					"status": false,
				})
			}

			now := time.Now()
			diff := now.Sub(pReset.CreatedAt)

			if PasswordResetTimeOut <= 1 {
				PasswordResetTimeOut = 3
			}
			mins := int(diff.Minutes())

			if PasswordResetTimeOut >= mins && pReset.Wrong <= 21 {

				if agentUtils.IsSame(data.Code, pReset.Token) {

					if data.Password == data.PasswordConfirm {

						newPassword, _ := agentUtils.Hash(data.Password)

						foundUser.Password = newPassword
						err := DB.DB.Save(foundUser).Error
						if err != nil {

							return c.JSON(map[string]interface{}{
								"status": false,
							})
						} else {
							DB.DB.Where("email = ?", data.Email).Delete(pReset)
							return c.JSON(map[string]interface{}{
								"msg":    StaticWords["passwordResetSuccess"],
								"status": true,
							})
						}

					} else {

						return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
							"error":  StaticWords["passwordConfirmError"],
							"msg":    StaticWords["passwordConfirmError"],
							"status": false,
						})
					}

				} else {
					pReset.Wrong = pReset.Wrong + 1
					DB.DB.Save(pReset)
					return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
						"error":  StaticWords["passwordResetCodeIncorrect"],
						"msg":    StaticWords["passwordResetCodeIncorrect"],
						"status": false,
					})
				}
			} else {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  StaticWords["passwordResetCodeTimeout"],
					"msg":    StaticWords["passwordResetCodeTimeout"],
					"status": false,
				})
			}
		} else {
			foundUser := models.User{}
			pReset := models.PasswordReset{}
			PasswordResetTimeOut := config.LambdaConfig.PasswordResetTimeOut

			errU := DB.DB.Where("email = ?", data.Email).Find(&foundUser).Error
			errR := DB.DB.Where("email = ?", data.Email).Find(&pReset).Error

			if errU != nil || foundUser.Login == "" {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  StaticWords["userNotFound"],
					"msg":    StaticWords["userNotFound"],
					"status": false,
				})
			}
			if errR != nil || pReset.Email == "" {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  StaticWords["userNotFound"],
					"msg":    StaticWords["userNotFound"],
					"status": false,
				})
			}

			now := time.Now()
			diff := now.Sub(pReset.CreatedAt)

			if PasswordResetTimeOut <= 1 {
				PasswordResetTimeOut = 3
			}
			mins := int(diff.Minutes())

			if PasswordResetTimeOut >= mins && pReset.Wrong <= 2 {

				if agentUtils.IsSame(data.Code, pReset.Token) {

					if data.Password == data.PasswordConfirm {

						newPassword, _ := agentUtils.Hash(data.Password)

						foundUser.Password = newPassword
						err := DB.DB.Save(foundUser).Error
						if err != nil {

							return c.JSON(map[string]interface{}{
								"status": false,
							})
						} else {
							DB.DB.Where("email = ?", data.Email).Delete(pReset)
							return c.JSON(map[string]interface{}{
								"msg":    StaticWords["passwordResetSuccess"],
								"status": true,
							})
						}

					} else {

						return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
							"error":  StaticWords["passwordConfirmError"],
							"msg":    StaticWords["passwordConfirmError"],
							"status": false,
						})
					}

				} else {
					pReset.Wrong = pReset.Wrong + 1
					DB.DB.Save(pReset)
					return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
						"error":  StaticWords["passwordResetCodeIncorrect"],
						"msg":    StaticWords["passwordResetCodeIncorrect"],
						"status": false,
					})
				}
			} else {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"error":  StaticWords["passwordResetCodeTimeout"],
					"msg":    StaticWords["passwordResetCodeTimeout"],
					"status": false,
				})
			}
		}
	}

	return c.JSON(map[string]interface{}{
		"status": false,
	})
}
func shuffle(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}
