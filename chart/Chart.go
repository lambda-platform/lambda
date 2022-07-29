package chart

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/chart/handler"
)

func Set(e *fiber.App) {
	ve := e.Group("/ve")
	ve.Post("/get-data-count", handler.CountData)
	ve.Post("/get-data-pie", handler.PieData)
	ve.Post("/get-data-table", handler.TableData)
	ve.Post("/get-data", handler.LineData)
	/* ROUTES */
}
