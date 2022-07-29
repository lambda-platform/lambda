package dataform

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Exec(c *fiber.Ctx, schemaId string, action string, id string, GetMODEL func(schemaId string) Dataform) error {
	datform := GetMODEL(schemaId)

	switch action {
	case "store":
		return Store(c, datform, action, id)
	case "update":
		return Store(c, datform, action, id)
	case "edit":
		return Edit(c, datform, id)
	case "options":
		return Options(c)
	}

	return c.Status(http.StatusBadRequest).JSON(map[string]string{
		"status": "false",
	})

}
