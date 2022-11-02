package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/utils"
	"strconv"

	agentModels "github.com/lambda-platform/lambda/agent/models"
)

func GetUsers(c *fiber.Ctx) error {
	role := c.Query("role")
	sort := c.Query("sort")
	direction := c.Query("direction")

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
		return c.JSON(data)
	} else {
		if config.Config.Database.Connection == "oracle" {
			users := []agentModels.UserWithoutPasswordOracle{}
			data := utils.Paging(&utils.Param{
				DB:    query,
				Page:  GetPage(c),
				Limit: 16,
			}, &users)
			return c.JSON(data)
		} else {
			users := []agentModels.UserWithoutPassword{}
			data := utils.Paging(&utils.Param{
				DB:    query,
				Page:  GetPage(c),
				Limit: 16,
			}, &users)
			return c.JSON(data)
		}
	}

}

func SearchUsers(c *fiber.Ctx) error {
	q := c.Params("q")

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

		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   data,
		})
	} else {
		if config.Config.Database.Connection == "oracle" {
			users := []agentModels.UserWithoutPasswordOracle{}
			data := utils.Paging(&utils.Param{
				DB:    query,
				Page:  GetPage(c),
				Limit: 16,
			}, &users)

			return c.JSON(map[string]interface{}{
				"status": true,
				"data":   data,
			})
		} else {
			users := []agentModels.UserWithoutPassword{}
			data := utils.Paging(&utils.Param{
				DB:    query,
				Page:  GetPage(c),
				Limit: 16,
			}, &users)

			return c.JSON(map[string]interface{}{
				"status": true,
				"data":   data,
			})
		}
	}
}
func GetDeletedUsers(c *fiber.Ctx) error {
	role := c.Query("role")
	sort := c.Query("sort")
	direction := c.Query("direction")

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

		return c.JSON(map[string]interface{}{
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

		return c.JSON(map[string]interface{}{
			"status": "true",
			"data":   data,
		})
	}

}

func GetRoles(c *fiber.Ctx) error {

	roles := []agentModels.Role{}
	DB.DB.Where("id != 1").Find(&roles)
	return c.JSON(roles)
}

func DeleteUser(c *fiber.Ctx) error {
	if config.Config.SysAdmin.UUID {
		id := c.Params("id")
		user := new(agentModels.UserWithoutPasswordUUID)
		err := DB.DB.Where("id = ?", id).Delete(&user).Error
		if err != nil {
			return c.JSON(map[string]string{
				"status": "false",
			})
		} else {
			return c.JSON(map[string]string{
				"status": "true",
			})
		}
	} else {
		id := c.Params("id")
		user := new(agentModels.UserWithoutPassword)
		err := DB.DB.Where("id = ?", id).Delete(&user).Error
		if err != nil {
			return c.JSON(map[string]string{
				"status": "false",
			})
		} else {
			return c.JSON(map[string]string{
				"status": "true",
			})
		}
	}

}
func GetPage(c *fiber.Ctx) int {
	page := c.Query("page")

	var Page_ int = 1
	if page != "" {
		Page_, _ = strconv.Atoi(page)
	}

	return Page_
}
