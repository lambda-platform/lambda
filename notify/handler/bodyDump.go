package handler

import (
	"github.com/gofiber/fiber/v2"
)

func BodyDump(c *fiber.Ctx, reqBody, resBody []byte) {

	//action := c.Param("action")
	//if(action == "store" || action == "update" || action == "delete" || action == "edit"){
	//
	//	schemaId, _ := strconv.ParseInt(c.Param("schemaId"), 10, 64)
	//	user := c.Get("user").(*jwt.Token)
	//	claims := user.Claims.(jwt.MapClaims)
	//	userID := claims["id"].(float64)
	//
	//	if(action == "store" || action == "update" || action == "delete"){
	//		BuildNotification(reqBody, schemaId, action, int64(userID))
	//	}
	//
	//}
	return

}
