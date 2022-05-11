package dataform

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"net/http"
	"fmt"
)

type UniquePost struct {
	Table          string `json:"table"`
	IdentityColumn string `json:"identityColumn"`
	Identity       interface{} `json:"identity"`
	Field          string `json:"field"`
	Val            interface{}  `json:"val"`
}

func CheckUnique(c echo.Context) error {

	post := new(UniquePost)
	if err := c.Bind(post); err != nil {

		return c.JSON(http.StatusOK, map[string]interface{}{
			"msg": "false from json",
			"status": false,
		})
	}

	DB_ := DB.DBConnection()
	var count int

	IdentityValue := fmt.Sprintf("%v", post.Identity)

	if IdentityValue == "<nil>"{
		IdentityValue = ""
	}

	valueType := fmt.Sprintf("%T", post.Val)
	value := ""
	if valueType == "float64" {
		value = fmt.Sprintf("%d", int(post.Val.(float64)))
	} else if valueType == "float32" {
		value = fmt.Sprintf("%d", int(post.Val.(float32)))
	} else if valueType == "int" {
		value = fmt.Sprintf("%d", int(post.Val.(int)))
	}else if valueType == "int32" {
		value = fmt.Sprintf("%d", int(post.Val.(int32)))
	} else if valueType == "int64" {
		value = fmt.Sprintf("%d", int(post.Val.(int64)))
	} else {
		value = post.Val.(string)
	}

	if post.IdentityColumn != "" && IdentityValue != "" {

		err := DB_.QueryRow("SELECT COUNT(*) FROM " + post.Table + " WHERE " + post.IdentityColumn + " != '" + IdentityValue + "' AND " + post.Field + " = '" + value + "'").Scan(&count)

		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": false,
				"msg": err.Error(),
			})
		} else {

			if count > 0 {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": false,
					"msg":    "'" + fmt.Sprintf("%v", value) + "' утга бүртгэлтэй байна",
				})
			} else {

				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": true,
				})
			}
		}
	} else {
		err := DB_.QueryRow("SELECT COUNT(*) FROM " + post.Table + " WHERE " + post.Field + " = '" + value + "'").Scan(&count)

		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": false,
				"msg": err.Error(),
			})
		} else {

			if count > 0 {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": false,
					"msg":    "'" + fmt.Sprintf("%v", value) + "' утга бүртгэлтэй байна",
				})
			} else {

				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": true,
				})
			}
		}
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": false,
	})
}

