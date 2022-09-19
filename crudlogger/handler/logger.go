package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/crudlogger/models"
	"strconv"
)

func CrudLogger(c *fiber.Ctx, reqBody, resBody []byte) {

	action := c.Params("action")
	if c.Path() == "/lambda/krud/delete/:schemaId/:id" {
		action = "delete"
	}
	if action == "store" || action == "update" || action == "delete" || action == "edit" {

		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := claims["id"].(float64)
		schemaId, _ := strconv.ParseInt(c.Params("schemaId"), 10, 64)
		RowId := c.Params("id")

		Log := models.CrudLog{
			UserId:    int64(userID),
			Ip:        c.IP(),
			UserAgent: string(c.Context().UserAgent()),
			Action:    action,
			SchemaId:  schemaId,
			RowId:     RowId,
			Input:     string(resBody),
		}

		if action == "store" {
			var response models.CrudResponse
			if err := json.Unmarshal(resBody, &response); err != nil {
				panic(err)
			}
			Log.RowId = strconv.Itoa(response.Data.ID)
		}
		DB.DB.Create(&Log)

	}
	return

}
