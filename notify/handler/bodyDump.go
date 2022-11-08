package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/dataform"
	"github.com/lambda-platform/lambda/datagrid"
	"strconv"
)

type crudResponse struct {
	Data map[string]interface{} `json:"data"`
}

func BodyDump(c *fiber.Ctx, GetGridMODEL func(schema_id string) datagrid.Datagrid, GetMODEL func(schema_id string) dataform.Dataform) error {
	if err := c.Next(); err != nil {
		return err
	}
	action := c.Params("action")
	if action == "store" || action == "update" || action == "delete" || action == "edit" {

		schemaId, _ := strconv.ParseInt(c.Params("schemaId"), 10, 64)
		user := agentUtils.AuthUserObject(c)

		var response crudResponse

		if err := json.Unmarshal(c.Response().Body(), &response); err != nil {
			panic(err)
		}

		if action == "store" || action == "update" || action == "delete" {
			go BuildNotification(response.Data, schemaId, action, user)
		}

	}

	return nil
}
