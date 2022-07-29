package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/utils"

	"strconv"
)

func FetchData(c *fiber.Ctx, datagrid Datagrid) error {

	pageLimit := c.Query("paginate")
	page := c.Query("page")
	sort := c.Query("sort")
	order := c.Query("order")

	query := DB.DB.Table(datagrid.DataTable)

	query = query.Select(datagrid.ColumnList)

	if sort != "null" && sort != "undefined" {
		if order == "asc" || order == "desc" || order == "ASC" || order == "DESC" {
			query = query.Order(sort + " " + order)
		}
	}

	//DB.DB.LogMode(true)
	query, _ = Filter(c, datagrid, query)

	if len(datagrid.Condition) > 0 {
		query = query.Where(datagrid.Condition)
	}
	query = Search(c, datagrid.DataModel, query)
	var Page_ int = 1
	if page != "" {
		Page_, _ = strconv.Atoi(page)
	}
	Limit_, _ := strconv.Atoi(pageLimit)

	triggerData, query, skipPagination, returnData := ExecTrigger("beforeFetch", datagrid.Data, datagrid, query, c)

	if returnData {

		return c.JSON(triggerData)

	} else if skipPagination {
		data := query.Find(&datagrid.Data)
		res := utils.Paginator{
			Data:        &datagrid.Data,
			Total:       int64(int(data.RowsAffected)),
			CurrentPage: 1,
		}
		return c.JSON(res)

	} else {

		data := utils.Paging(&utils.Param{
			DB:    query,
			Page:  Page_,
			Limit: Limit_,
		}, datagrid.Data)

		if len(datagrid.Relations) >= 1 {
			data.Data = datagrid.FillVirtualColumns(datagrid.Data)
		}

		return c.JSON(data)
	}

}
