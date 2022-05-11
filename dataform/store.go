package dataform

import (

	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Store(c echo.Context, dataform Dataform, action string, id string) error {

	/*FORM VALIDATION*/
	requestData, err := GetData(c, action, id, dataform)


	if err !=nil{
		errData := map[string]interface{}{"error": err}
		errData["status"] = false
		return c.JSON(http.StatusBadRequest, err)
	}


	if len(dataform.ValidationRules) >= 1 {
		opts := govalidator.Options{
			Data:     dataform.Model,    // request object
			Rules:    dataform.ValidationRules,    // rules map
			Messages: dataform.ValidationMessages, // custom message map (Optional)
			//RequiredDefault: false,     // all the field to be pass the rules
		}
		v := govalidator.New(opts)
		e := v.ValidateStruct()



		if len(e) >= 1 {
			err := map[string]interface{}{"error": e}
			err["status"] = false
			return c.JSON(http.StatusBadRequest, err)
		}
	}


	if id != "" {
		err := DB.DB.Save(dataform.Model).Error
		if err != nil {

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {

			saveNestedSubItem(dataform, *requestData)

			if dataform.TriggerNameSpace != "" && dataform.AfterUpdate != nil {
				dataform.AfterUpdate(dataform.Model)
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": true,
				"data":   dataform.Model,
			})
		}
	} else {



		err := DB.DB.Create(dataform.Model).Error

		if err != nil {

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": false,
				"error":  err.Error(),
			})
		} else {

			saveNestedSubItem(dataform, *requestData)


			if dataform.TriggerNameSpace != "" && dataform.AfterInsert != nil {
				dataform.AfterInsert(dataform.Model)
			}
	

			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": true,
				"data":   dataform.Model,
			})
		}
	}

}




