package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils"
	"net/http"
	"strconv"

	agentModels "github.com/lambda-platform/lambda/agent/models"
)

func GetUsers(c echo.Context) error {
	role := c.QueryParam("role")
	sort := c.QueryParam("sort")
	direction := c.QueryParam("direction")

	query := DB.DB.Table("users").Order(sort + " " + direction)

	if role != "all" {
		query = query.Where("role = ?", role)
	}
	query = query.Where("deleted_at IS NULL")

	if config.Config.SysAdmin.UUID {
		users := []agentModels.UserWithoutPasswordUUID{}
		data := utils.Paging(&utils.Param{
			DB:    query,
			Page:  GetPage(c),
			Limit: 16,
		}, &users)
		return c.JSON(http.StatusOK, data)
	} else {
		users := []agentModels.UserWithoutPassword{}
		data := utils.Paging(&utils.Param{
			DB:    query,
			Page:  GetPage(c),
			Limit: 16,
		}, &users)
		return c.JSON(http.StatusOK, data)
	}

}

func SearchUsers(c echo.Context) error {
	q := c.Param("q")

	query := DB.DB.Table("users")

	if q != "=" {
		query = query.Where("deleted_at IS NULL")
		query = query.Where("login LIKE ?", "%"+q+"%")
		query = query.Or("first_name LIKE ?", "%"+q+"%")
		query = query.Or("last_name LIKE ?", "%"+q+"%")
		query = query.Or("register_number LIKE ?", "%"+q+"%")
		query = query.Or("phone LIKE ?", "%"+q+"%")
	}
	if config.Config.SysAdmin.UUID {
		users := []agentModels.UserWithoutPasswordUUID{}
		data := utils.Paging(&utils.Param{
			DB:    query,
			Page:  GetPage(c),
			Limit: 16,
		}, &users)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "true",
			"data":   data,
		})
	} else {
		users := []agentModels.UserWithoutPassword{}
		data := utils.Paging(&utils.Param{
			DB:    query,
			Page:  GetPage(c),
			Limit: 16,
		}, &users)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "true",
			"data":   data,
		})
	}
}
func GetDeletedUsers(c echo.Context) error {
	role := c.QueryParam("role")
	sort := c.QueryParam("sort")
	direction := c.QueryParam("direction")

	query := DB.DB.Table("users").Order(sort + " " + direction)

	if role != "all" {
		query = query.Where("role = ?", role)
	}
	query = query.Where("deleted_at IS NOT NULL")

	if config.Config.SysAdmin.UUID {
		users := []agentModels.UserWithoutPasswordUUID{}
		data := utils.Paging(&utils.Param{
			DB:    query,
			Page:  GetPage(c),
			Limit: 16,
		}, &users)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "true",
			"data":   data,
		})
	} else {
		users := []agentModels.UserWithoutPassword{}
		data := utils.Paging(&utils.Param{
			DB:    query,
			Page:  GetPage(c),
			Limit: 16,
		}, &users)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "true",
			"data":   data,
		})
	}

}

func GetRoles(c echo.Context) error {

	roles := []agentModels.Role{}
	DB.DB.Where("id != 1").Find(&roles)
	return c.JSON(http.StatusOK, roles)
}

func DeleteUser(c echo.Context) error {
	if config.Config.SysAdmin.UUID {
		id := c.Param("id")
		user := new(agentModels.UserWithoutPasswordUUID)
		err := DB.DB.Where("id = ?", id).Delete(&user).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "false",
			})
		} else {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "true",
			})
		}
	} else {
		id := c.Param("id")
		user := new(agentModels.UserWithoutPassword)
		err := DB.DB.Where("id = ?", id).Delete(&user).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "false",
			})
		} else {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "true",
			})
		}
	}

}
func GetPage(c echo.Context) int {
	page := c.QueryParam("page")

	var Page_ int = 1
	if page != "" {
		Page_, _ = strconv.Atoi(page)
	}

	return Page_
}
