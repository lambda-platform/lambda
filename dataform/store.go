package dataform

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Store(c *fiber.Ctx, dataform Dataform, action string, id string) error {

	/*FORM VALIDATION*/
	requestData, err := GetData(c, action, id, dataform)
	fmt.Println("hi1")
	if err != nil {
		errData := map[string]interface{}{"error": err}
		errData["status"] = false
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	fmt.Println("hi2")

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
		err := DB.DB.Save(dataform.Model).Error
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

			return c.JSON(map[string]interface{}{
				"status": true,
				"data":   dataform.Model,
			})
		}
	}

}
