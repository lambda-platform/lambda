package dataform

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/DBSchema"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils"
	"github.com/thedevsaddam/govalidator"
)

func Store(c *fiber.Ctx, dataform Dataform, action string, id string) error {

	/*FORM VALIDATION*/

	requestData, err := GetData(c, action, id, dataform)

	if err != nil {
		errData := map[string]interface{}{"error": "Invalid request data."}
		errData["status"] = false
		if config.Config.Database.Debug {
			errData["details"] = err.Error()
		}
		return c.Status(http.StatusBadRequest).JSON(errData)
	}

	// Check password strength if password field exists
	passwordField := "password" // Change this if your password field has a different name
	if pwd, exists := (*requestData)[passwordField]; exists {

		switch v := pwd.(type) {
		case nil:
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  "Password cannot be nil",
			})
		case string:

			if err = validatePasswordStrength(v); err != nil {

				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"status": false,
					"error":  err.Error(),
				})
			}
		case *string:
			if v == nil {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"status": false,
					"error":  "Password pointer cannot be nil",
				})
			}

			// Use *v directly instead of *v.(string)
			if err = validatePasswordStrength(*v); err != nil {
				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"status": false,
					"error":  err.Error(),
				})
			}
		default:
			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  "Invalid password format",
			})
		}
	}

	if len(dataform.ValidationRules) >= 1 {
		opts := govalidator.Options{
			Data:     dataform.Model,              // request object
			Rules:    dataform.ValidationRules,    // rules map
			Messages: dataform.ValidationMessages, // custom message map (Optional)
			//RequiredDefault: false,     // all the field to be pass the rules
		}
		v := govalidator.New(opts)
		e := v.ValidateStruct()

		if len(e) >= 1 {
			errValidation := map[string]interface{}{"error": e}
			errValidation["status"] = false
			return c.Status(http.StatusBadRequest).JSON(errValidation)
		}
	}

	if id != "" {
		query := DB.DB
		if config.Config.Database.Connection != "mysql" {
			if dataform.Identity != "ID" && dataform.Identity != "id" {
				query = query.Where(dataform.Identity+" = ?", id)
			}
		}

		err := query.Save(dataform.Model).Error
		if err != nil {

			errResponse := map[string]interface{}{
				"status": false,
				"error":  "Хадгалах үед алдаа гарлаа.",
			}
			if config.Config.Database.Debug {
				errResponse["details"] = err.Error()
			}
			return c.Status(http.StatusBadRequest).JSON(errResponse)
		} else {

			saveNestedSubItem(dataform, *requestData)

			if dataform.TriggerNameSpace != "" && dataform.AfterUpdate != nil {
				dataform.AfterUpdate(dataform.Model)
			}

			if dataMap, ok := utils.StructToMap(dataform.Model); ok {
				delete(dataMap, "password")
				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   dataMap,
					"id":     id,
				})
			} else {
				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   dataform.Model,
					"id":     id,
				})
			}
		}
	} else {

		err := DB.DB.Create(dataform.Model).Error

		if err != nil {

			errResponse := map[string]interface{}{
				"status": false,
				"error":  "Failed to create the record.",
			}
			if config.Config.Database.Debug {
				errResponse["details"] = err.Error()
			}
			return c.Status(http.StatusBadRequest).JSON(errResponse)
		} else {

			saveNestedSubItem(dataform, *requestData)

			if dataform.TriggerNameSpace != "" && dataform.AfterInsert != nil {
				dataform.AfterInsert(dataform.Model)
			}

			idValue, errIDValue := dataform.getFieldValue(DBSchema.FieldName(dataform.Identity))

			if errIDValue != nil {

				errResponse := map[string]interface{}{
					"status": false,
					"error":  "Failed to retrieve the ID.",
				}
				if config.Config.Database.Debug {
					errResponse["details"] = errIDValue.Error()
				}
				return c.Status(http.StatusBadRequest).JSON(errResponse)
			}

			if dataMap, ok := utils.StructToMap(dataform.Model); ok {
				delete(dataMap, "password")

				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   dataMap,
					"id":     idValue,
				})
			} else {
				return c.JSON(map[string]interface{}{
					"status": true,
					"data":   dataform.Model,
					"id":     idValue,
				})
			}

		}
	}

}

// Function to validate password strength
func validatePasswordStrength(password string) error {
	// Define password strength regex
	var (
		lowercase = regexp.MustCompile(`[a-z]`)
		uppercase = regexp.MustCompile(`[A-Z]`)
		number    = regexp.MustCompile(`[0-9]`)
		special   = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`)
	)

	// Check password length
	if len(password) < 8 {
		return errors.New("Нууц үг хамгийн багадаа 8 тэмдэгтээс бүрдэх ёстой.")
	}

	// Check character requirements
	if !lowercase.MatchString(password) {
		return errors.New("Нууц үг дор хаяж нэг жижиг үсэг агуулсан байх ёстой.")
	}
	if !uppercase.MatchString(password) {
		return errors.New("Нууц үг дор хаяж нэг том үсэг агуулсан байх ёстой.")
	}
	if !number.MatchString(password) {
		return errors.New("Нууц үг дор хаяж нэг тоо агуулсан байх ёстой.")
	}
	if !special.MatchString(password) {
		return errors.New("Нууц үг дор хаяж нэг тусгай тэмдэгт (!@#$%^&*) агуулсан байх ёстой.")
	}

	return nil
}
