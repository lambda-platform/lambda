package dataform

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/DBSchema"
	"github.com/lambda-platform/lambda/config"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Store(c *fiber.Ctx, dataform Dataform, action string, id string) error {

	/*FORM VALIDATION*/

	requestData, err := GetData(c, action, id, dataform)

	if err != nil {
		errData := map[string]interface{}{"error": err}
		errData["status"] = false
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
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
			err := map[string]interface{}{"error": e}
			err["status"] = false
			return c.Status(http.StatusBadRequest).JSON(err)
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

			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {

			saveNestedSubItem(dataform, *requestData)

			if dataform.TriggerNameSpace != "" && dataform.AfterUpdate != nil {
				dataform.AfterUpdate(dataform.Model)
			}

			return c.JSON(map[string]interface{}{
				"status": true,
				"data":   dataform.Model,
			})
		}
	} else {

		err := DB.DB.Create(dataform.Model).Error

		if err != nil {

			return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {

			saveNestedSubItem(dataform, *requestData)

			if dataform.TriggerNameSpace != "" && dataform.AfterInsert != nil {
				dataform.AfterInsert(dataform.Model)
			}

			idValue, errIDValue := dataform.getFieldValue(DBSchema.FieldName(dataform.Identity))

			if errIDValue != nil {

				return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
					"status": false,
					"error":  errIDValue.Error(),
				})
			}

			return c.JSON(map[string]interface{}{
				"status": true,
				"data":   dataform.Model,
				"id":     idValue,
			})
		}
	}

}
