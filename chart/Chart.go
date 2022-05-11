package chart

import (
	echo "github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/chart/handler"
)

func Set(e *echo.Echo) {
	ve := e.Group("/ve")
	ve.POST("/get-data-count", handler.CountData)
	ve.POST("/get-data-pie", handler.PieData)
	ve.POST("/get-data-table", handler.TableData)
	ve.POST("/get-data", handler.LineData)
	/* ROUTES */
}
