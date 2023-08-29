package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"strconv"
	"strings"
)

type crudResponse struct {
	Data map[string]interface{} `json:"data"`
}

func BodyDump(c *fiber.Ctx, GetGridMODEL func(schema_id string) datagrid.Datagrid, GetMODEL func(schema_id string) dataform.Dataform) error {
	if err := c.Next(); err != nil {
		return err
	}

	if !strings.HasPrefix(c.Path(), "/lambda/krud-public") {
		action := c.Params("action")
		if action == "store" || action == "update" || action == "delete" || action == "edit" {

			schemaId, _ := strconv.ParseInt(c.Params("schemaId"), 10, 64)
			user, err := agentUtils.AuthUserObject(c)

			if err != nil {
				c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":  err.Error(),
					"from":   "Notification builder",
					"status": false,
				})
			}

			var response crudResponse

			if err := json.Unmarshal(c.Response().Body(), &response); err != nil {
				panic(err)
			}

			if action == "store" || action == "update" || action == "delete" {
				BuildNotification(response.Data, schemaId, action, user)
			}

		}
	}

	return nil
}
