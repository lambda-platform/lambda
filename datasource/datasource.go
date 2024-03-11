package datasource

import (
	"encoding/json"
	"fmt"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"regexp"
	"strings"
)

func DeleteView(viewName string) {
	DB.DB.Exec("DROP VIEW IF EXISTS " + viewName)
}

func CreateView(viewName string, schema string) error {

	schemaData := make(map[string]interface{})
	json.Unmarshal([]byte(schema), &schemaData)
	DeleteView(viewName)

	queryPre, firstSchema := addSchemaToTablesInQuery(schemaData["query"].(string))
	query := "CREATE VIEW " + firstSchema + viewName + " as " + queryPre

	err := DB.DB.Exec(query).Error
	return err
}

func addSchemaToTablesInQuery(query string) (string, string) {
	var firstSchema string

	if config.Config.Database.Connection == "postgres" {
		// Match FROM and JOIN clauses, capturing table references including potential schema prefix
		re := regexp.MustCompile(`(?i)(\bFROM\b|\bJOIN\b)\s+([a-zA-Z_][\w.]*)(\s+AS\s+\w+)?`)

		// This map is used to ensure we only process each unique table name once
		processedTables := make(map[string]bool)

		// Find all matches in the query
		matches := re.FindAllStringSubmatch(query, -1)

		for _, match := range matches {
			// Full match is match[0], the table reference (including potential schema) is match[2]
			tableRef := match[2]

			// Skip if we've already processed this table
			if _, processed := processedTables[tableRef]; processed {
				continue
			}
			processedTables[tableRef] = true

			// Extract table and schema names, if present
			parts := strings.SplitN(tableRef, ".", 2)
			var schemaName, tableName string

			if len(parts) == 2 {
				schemaName, tableName = parts[0], parts[1]
			} else {
				tableName = parts[0]
				// Assuming DB is your GORM database instance
				// Fetch schema name from the database if not specified
				DB.DB.Raw("SELECT table_schema FROM information_schema.tables WHERE table_name = ? LIMIT 1", tableName).Scan(&schemaName)
			}

			// Save the first encountered schema name
			if firstSchema == "" {
				firstSchema = schemaName + "."
			}

			// Correctly formatted schema-qualified table name
			correctTableName := fmt.Sprintf("%s.%s", schemaName, tableName)
			// Replace original table reference with the correctly formatted name
			query = strings.ReplaceAll(query, tableRef, correctTableName)
		}
	}

	return query, firstSchema
}
