package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"strconv"
)

// Handler to fetch logs with filters and pagination
func GetLogs(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	perPageStr := c.Query("per_page", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	baseQuery := `
		SELECT 
			id,
			user_id,
			ip,
			user_agent,
			action,
			schema_id,
			row_id,
			created_at,
			input,
			name AS schema_name,
			db_table_name,
			role,
			login,
			email,
			first_name,
			last_name,
			geo_type AS log_type
		FROM view_crud_log
		WHERE 1=1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM view_crud_log
		WHERE 1=1
	`

	var params []interface{}
	var countParams []interface{}

	if category := c.Query("category"); category != "" {
		baseQuery += ` AND name = ?`
		countQuery += ` AND name = ?`
		params = append(params, category)
		countParams = append(countParams, category)
	}

	if logType := c.Query("type"); logType != "" {
		baseQuery += ` AND geo_type = ?`
		countQuery += ` AND geo_type = ?`
		params = append(params, logType)
		countParams = append(countParams, logType)
	}

	if action := c.Query("action"); action != "" {
		baseQuery += ` AND action = ?`
		countQuery += ` AND action = ?`
		params = append(params, action)
		countParams = append(countParams, action)
	}

	// Optional: Add date filters if needed (example for created_at range)
	if fromDate := c.Query("from_date"); fromDate != "" {
		baseQuery += ` AND created_at >= (?::date)`
		countQuery += ` AND created_at >= (?::date)`
		params = append(params, fromDate)
		countParams = append(countParams, fromDate)
	}
	if toDate := c.Query("to_date"); toDate != "" {
		baseQuery += ` AND created_at < (?::date + interval '1 day')`
		countQuery += ` AND created_at < (?::date + interval '1 day')`
		params = append(params, toDate)
		countParams = append(countParams, toDate)
	}

	baseQuery += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	params = append(params, perPage, offset)

	var logs []map[string]interface{}
	if err := DB.DB.Raw(baseQuery, params...).Scan(&logs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": "Мэдээлэл татахад алдаа гарлаа: " + err.Error()})
	}

	var total int
	if err := DB.DB.Raw(countQuery, countParams...).Scan(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": "Нийт тоо тооцоход алдаа гарлаа: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":     logs,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// Handler to get unique categories (based on vb_schemas.name) that have associated logs
func GetLogCategories(c *fiber.Ctx) error {
	var categories []string
	if err := DB.DB.Raw(`
		SELECT DISTINCT vb_schemas.name 
		FROM vb_schemas 
		INNER JOIN crud_log ON vb_schemas.id = crud_log.schema_id
		WHERE vb_schemas.type = ANY (ARRAY['form'::text, 'grid'::text])
		ORDER BY vb_schemas.name
	`).Scan(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": "Ангилал татахад алдаа гарлаа: " + err.Error()})
	}
	return c.JSON(categories)
}

// Handler to get unique actions (for filter dropdown)
func GetLogActions(c *fiber.Ctx) error {
	var actions []string
	if err := DB.DB.Raw(`
		SELECT DISTINCT action 
		FROM crud_log 
		ORDER BY action
	`).Scan(&actions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": "Үйлдэл татахад алдаа гарлаа: " + err.Error()})
	}
	return c.JSON(actions)
}

// Chart: Log by Category
func GetLogChartCategory(c *fiber.Ctx) error {
	baseQuery := `
		SELECT name AS category, COUNT(*) AS count
		FROM view_crud_log
		WHERE 1=1
	`

	var params []interface{}

	if category := c.Query("category"); category != "" {
		baseQuery += ` AND name = ?`
		params = append(params, category)
	}

	if logType := c.Query("type"); logType != "" {
		baseQuery += ` AND geo_type = ?`
		params = append(params, logType)
	}

	if action := c.Query("action"); action != "" {
		baseQuery += ` AND action = ?`
		params = append(params, action)
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		baseQuery += ` AND created_at >= (?::date)`
		params = append(params, fromDate)
	}
	if toDate := c.Query("to_date"); toDate != "" {
		baseQuery += ` AND created_at < (?::date + interval '1 day')`
		params = append(params, toDate)
	}

	baseQuery += ` GROUP BY name ORDER BY count DESC`

	var data []map[string]interface{}
	if err := DB.DB.Raw(baseQuery, params...).Scan(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": err.Error()})
	}
	return c.JSON(data)
}

// Chart: Log by Type
func GetLogChartType(c *fiber.Ctx) error {
	baseQuery := `
		SELECT geo_type AS type, COUNT(*) AS count
		FROM view_crud_log
		WHERE 1=1
	`

	var params []interface{}

	if category := c.Query("category"); category != "" {
		baseQuery += ` AND name = ?`
		params = append(params, category)
	}

	if logType := c.Query("type"); logType != "" {
		baseQuery += ` AND geo_type = ?`
		params = append(params, logType)
	}

	if action := c.Query("action"); action != "" {
		baseQuery += ` AND action = ?`
		params = append(params, action)
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		baseQuery += ` AND created_at >= (?::date)`
		params = append(params, fromDate)
	}
	if toDate := c.Query("to_date"); toDate != "" {
		baseQuery += ` AND created_at < (?::date + interval '1 day')`
		params = append(params, toDate)
	}

	baseQuery += ` GROUP BY geo_type ORDER BY count DESC`

	var data []map[string]interface{}
	if err := DB.DB.Raw(baseQuery, params...).Scan(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": err.Error()})
	}
	return c.JSON(data)
}

// Chart: Log by User
func GetLogChartUser(c *fiber.Ctx) error {
	baseQuery := `
		SELECT CONCAT(first_name, ' ', last_name, ' (', login, ')') AS user, COUNT(*) AS count
		FROM view_crud_log
		WHERE 1=1
	`

	var params []interface{}

	if category := c.Query("category"); category != "" {
		baseQuery += ` AND name = ?`
		params = append(params, category)
	}

	if logType := c.Query("type"); logType != "" {
		baseQuery += ` AND geo_type = ?`
		params = append(params, logType)
	}

	if action := c.Query("action"); action != "" {
		baseQuery += ` AND action = ?`
		params = append(params, action)
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		baseQuery += ` AND created_at >= (?::date)`
		params = append(params, fromDate)
	}
	if toDate := c.Query("to_date"); toDate != "" {
		baseQuery += ` AND created_at < (?::date + interval '1 day')`
		params = append(params, toDate)
	}

	baseQuery += ` GROUP BY user_id, first_name, last_name, login ORDER BY count DESC LIMIT 10`

	var data []map[string]interface{}
	if err := DB.DB.Raw(baseQuery, params...).Scan(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": err.Error()})
	}
	return c.JSON(data)
}

// Chart: Log by Action
func GetLogChartAction(c *fiber.Ctx) error {
	baseQuery := `
		SELECT action, COUNT(*) AS count
		FROM view_crud_log
		WHERE 1=1
	`

	var params []interface{}

	if category := c.Query("category"); category != "" {
		baseQuery += ` AND name = ?`
		params = append(params, category)
	}

	if logType := c.Query("type"); logType != "" {
		baseQuery += ` AND geo_type = ?`
		params = append(params, logType)
	}

	if action := c.Query("action"); action != "" {
		baseQuery += ` AND action = ?`
		params = append(params, action)
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		baseQuery += ` AND created_at >= (?::date)`
		params = append(params, fromDate)
	}
	if toDate := c.Query("to_date"); toDate != "" {
		baseQuery += ` AND created_at < (?::date + interval '1 day')`
		params = append(params, toDate)
	}

	baseQuery += ` GROUP BY action ORDER BY count DESC`

	var data []map[string]interface{}
	if err := DB.DB.Raw(baseQuery, params...).Scan(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": err.Error()})
	}
	return c.JSON(data)
}

// Chart: Monthly Log Count (Bar)
func GetLogChartMonthly(c *fiber.Ctx) error {
	baseQuery := `
		SELECT TO_CHAR(date_trunc('month', created_at), 'YYYY-MM') AS month, COUNT(*) AS count
		FROM view_crud_log
		WHERE 1=1
	`

	var params []interface{}

	if category := c.Query("category"); category != "" {
		baseQuery += ` AND name = ?`
		params = append(params, category)
	}

	if logType := c.Query("type"); logType != "" {
		baseQuery += ` AND geo_type = ?`
		params = append(params, logType)
	}

	if action := c.Query("action"); action != "" {
		baseQuery += ` AND action = ?`
		params = append(params, action)
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		baseQuery += ` AND created_at >= (?::date)`
		params = append(params, fromDate)
	}
	if toDate := c.Query("to_date"); toDate != "" {
		baseQuery += ` AND created_at < (?::date + interval '1 day')`
		params = append(params, toDate)
	}

	baseQuery += ` GROUP BY date_trunc('month', created_at) ORDER BY month`

	var data []map[string]interface{}
	if err := DB.DB.Raw(baseQuery, params...).Scan(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": err.Error()})
	}
	return c.JSON(data)
}

// Creative Chart: Log by Hour (Hourly distribution, e.g., peak hours)
func GetLogChartHourly(c *fiber.Ctx) error {
	baseQuery := `
		SELECT EXTRACT(HOUR FROM created_at) AS hour, COUNT(*) AS count
		FROM view_crud_log
		WHERE 1=1
	`

	var params []interface{}

	if category := c.Query("category"); category != "" {
		baseQuery += ` AND name = ?`
		params = append(params, category)
	}

	if logType := c.Query("type"); logType != "" {
		baseQuery += ` AND geo_type = ?`
		params = append(params, logType)
	}

	if action := c.Query("action"); action != "" {
		baseQuery += ` AND action = ?`
		params = append(params, action)
	}

	if fromDate := c.Query("from_date"); fromDate != "" {
		baseQuery += ` AND created_at >= (?::date)`
		params = append(params, fromDate)
	}
	if toDate := c.Query("to_date"); toDate != "" {
		baseQuery += ` AND created_at < (?::date + interval '1 day')`
		params = append(params, toDate)
	}

	baseQuery += ` GROUP BY EXTRACT(HOUR FROM created_at) ORDER BY hour`

	var data []map[string]interface{}
	if err := DB.DB.Raw(baseQuery, params...).Scan(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"алдаа": err.Error()})
	}
	return c.JSON(data)
}
