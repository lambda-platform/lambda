package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/crudlogger/models"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"strconv"
	"strings"
)

func CrudLogger(UserAgent string, IP string, action string, resBody []byte, userID interface{}, schemaId int64, RowId string) {

	ID := userID.(int64)

	if action == "store" {
		var response models.CrudResponse
		if err := json.Unmarshal(resBody, &response); err != nil {
			panic(err)
		}
		RowId = strconv.Itoa(response.Data.ID)
	}
	if config.Config.Database.Connection == "oracle" {
		Log := models.CrudLogOracle{
			UserId:    ID,
			Ip:        IP,
			UserAgent: UserAgent,
			Action:    action,
			SchemaId:  schemaId,
			RowId:     RowId,
			Input:     string(resBody),
		}

		DB.DB.Create(&Log)
	} else {
		Log := models.CrudLog{
			UserId:    ID,
			Ip:        IP,
			UserAgent: UserAgent,
			Action:    action,
			SchemaId:  schemaId,
			RowId:     RowId,
			Input:     string(resBody),
		}

		DB.DB.Create(&Log)
	}

	return

}
func BodyDump(c *fiber.Ctx, GetGridMODEL func(schema_id string) datagrid.Datagrid, GetMODEL func(schema_id string) dataform.Dataform) error {
	if err := c.Next(); err != nil {
		return err
	}
	action := c.Params("action")
	if strings.Contains(c.Path(), "/lambda/krud/delete/") {
		action = "delete"
	}

	if action == "store" || action == "update" || action == "delete" || action == "edit" {
		if action == "edit" {
			action = "view"
		}
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		Id := claims["id"].(float64)
		schemaId, _ := strconv.ParseInt(c.Params("schemaId"), 10, 64)
		rowID := c.Params("id")
		CrudLogger(string(c.Context().UserAgent()), c.IP(), action, c.Response().Body(), int64(Id), schemaId, rowID)
	}

	return nil
}
