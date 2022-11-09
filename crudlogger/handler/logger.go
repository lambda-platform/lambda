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
	"reflect"
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

		RowId = GetID(response.ID)
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
func CrudLogHistory(c *fiber.Ctx) error {
	HistoryRequest := models.HistoryRequest{}

	if err := c.BodyParser(&HistoryRequest); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{
			"error": "Wrong request",
		})
	} else {
		if config.Config.Database.Connection == "oracle" {
			var crudLogs []models.CrudLogFullOracle

			DB.DB.Select("ID, LAST_NAME, FIRST_NAME, ACTION, CREATED_AT").Where("SCHEMA_ID = ? AND ROW_ID = ?", HistoryRequest.SchemaID, GetID(HistoryRequest.RowID)).Order("ID DESC").Find(&crudLogs)

			return c.JSON(crudLogs)
		} else {
			var crudLogs []models.CrudLogFull

			DB.DB.Select("id, last_name, first_name, action, created_at").Where("schema_id = ? AND row_id = ?", HistoryRequest.SchemaID, GetID(HistoryRequest.RowID)).Order("id DESC").Find(&crudLogs)

			return c.JSON(crudLogs)
		}
	}

}
func GetID(idPre interface{}) string {
	var id string

	roleDataType := reflect.TypeOf(idPre).String()

	if roleDataType == "float64" {
		id = strconv.Itoa(int(idPre.(float64)))
	} else if roleDataType == "float32" {
		id = strconv.Itoa(int(idPre.(float32)))
	} else if roleDataType == "int" {
		id = strconv.Itoa(int(idPre.(int)))
	} else if roleDataType == "int32" {
		id = strconv.Itoa(int(idPre.(int32)))
	} else if roleDataType == "int64" {
		id = strconv.Itoa(int(idPre.(int64)))
	} else if roleDataType == "string" {
		id = idPre.(string)
	}
	return id
}
