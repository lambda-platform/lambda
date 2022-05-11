package handler

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/crudlogger/models"
	"strconv"
)

func CrudLogger(c echo.Context, reqBody, resBody []byte) {

	action := c.Param("action")
	if (c.Path() == "/lambda/krud/delete/:schemaId/:id") {
		action = "delete"
	}
	if (action == "store" || action == "update" || action == "delete" || action == "edit") {

		req := c.Request()
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := claims["id"].(float64)
		schemaId, _ := strconv.ParseInt(c.Param("schemaId"), 10, 64)
		RowId := c.Param("id")

		Log := models.CrudLog{
			UserId:    int64(userID),
			Ip:        c.RealIP(),
			UserAgent: req.UserAgent(),
			Action:    action,
			SchemaId:  schemaId,
			RowId:     RowId,
			Input:     string(resBody),
		}

		if (action == "store") {
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
