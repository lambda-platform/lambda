package handler

import (
	"github.com/golang-jwt/jwt"
	"strconv"
	"github.com/labstack/echo/v4"
)

func BodyDump(c echo.Context, reqBody, resBody []byte) {

	action := c.Param("action")
	if(action == "store" || action == "update" || action == "delete" || action == "edit"){

		schemaId, _ := strconv.ParseInt(c.Param("schemaId"), 10, 64)
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := claims["id"].(float64)

		if(action == "store" || action == "update" || action == "delete"){
			BuildNotification(reqBody, schemaId, action, int64(userID))
		}

	}
	return

}
