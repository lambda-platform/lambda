package dataform

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"strings"
)

type UniquePost struct {
	Table          string      `json:"table"`
	IdentityColumn string      `json:"identityColumn"`
	Identity       interface{} `json:"identity"`
	Field          string      `json:"field"`
	Val            interface{} `json:"val"`
}

func CheckUnique(c *fiber.Ctx) error {

	post := new(UniquePost)
	if err := c.BodyParser(post); err != nil {

		return c.JSON(map[string]interface{}{
			"msg":    "false from json",
			"status": false,
		})
	}

	var count int64

	IdentityValue := fmt.Sprintf("%v", post.Identity)

	if IdentityValue == "<nil>" {
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
	} else if valueType == "int32" {
		value = fmt.Sprintf("%d", int(post.Val.(int32)))
	} else if valueType == "int64" {
		value = fmt.Sprintf("%d", int(post.Val.(int64)))
	} else {
		value = post.Val.(string)
	}

	if post.IdentityColumn != "" && IdentityValue != "" {

		var err error
		if (post.Table == "users" || post.Table == "USERS") && (post.Field == "login" || post.Field == "LOGIN") {
			err = DB.DB.Table(post.Table).
				Where(post.IdentityColumn+" != ? AND LOWER("+post.Field+") = ?", IdentityValue, strings.ToLower(value)).
				Count(&count).Error

		} else {
			err = DB.DB.Table(post.Table).Where(post.IdentityColumn+" != ? AND "+post.Field+" = ?", IdentityValue, value).Count(&count).Error
		}

		if err != nil {
			return c.JSON(map[string]interface{}{
				"status": false,
				"msg":    err.Error(),
			})
		} else {

			if count > 0 {
				return c.JSON(map[string]interface{}{
					"status": false,
					"msg":    "'" + fmt.Sprintf("%v", value) + "' утга бүртгэлтэй байна",
				})
			} else {

				return c.JSON(map[string]interface{}{
					"status": true,
				})
			}
		}
	} else {

		var err error
		if (post.Table == "users" || post.Table == "USERS") && (post.Field == "login" || post.Field == "LOGIN") {
			err = DB.DB.Table(post.Table).
				Where("LOWER("+post.Field+") = ?", strings.ToLower(value)).
				Count(&count).Error
		} else {
			err = DB.DB.Table(post.Table).Where(post.Field+" = ?", value).Count(&count).Error
		}

		if err != nil {
			return c.JSON(map[string]interface{}{
				"status": false,
				"msg":    err.Error(),
			})
		} else {

			if count > 0 {
				return c.JSON(map[string]interface{}{
					"status": false,
					"msg":    "'" + fmt.Sprintf("%v", value) + "' утга бүртгэлтэй байна",
				})
			} else {

				return c.JSON(map[string]interface{}{
					"status": true,
				})
			}
		}
	}

	return c.JSON(map[string]interface{}{
		"status": false,
	})
}
