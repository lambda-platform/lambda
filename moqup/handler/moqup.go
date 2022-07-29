package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils"
)

func Moqup(c *fiber.Ctx) error {
	id := c.Params("id")
	//csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	csrfToken := ""
	return c.Render("moqup.html", map[string]interface{}{
		"title":     config.LambdaConfig.Title,
		"favicon":   config.LambdaConfig.Favicon,
		"id":        id,
		"csrfToken": csrfToken,
		"mix":       utils.Mix,
	})
}
