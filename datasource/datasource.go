package datasource

import (
	"encoding/json"
	"fmt"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"regexp"
)

func DeleteView(viewName string) {
	DB.DB.Exec("DROP VIEW IF EXISTS " + viewName)
}

func CreateView(viewName, schema string) error {
	schemaData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(schema), &schemaData); err != nil {
		return err
	}

	queryPre, schemas := addSchemaToTablesInQuery(schemaData["query"].(string))
	// Using the first schema found for the view might not be appropriate for all cases.
	// Consider a strategy for determining the schema of the view if multiple are involved.
	viewFullName := viewName
	if len(schemas) > 0 {
		// Prepending the first found schema to the view name. This may need adjustment.
		viewFullName = schemas[0] + "." + viewName
	}
	DeleteView(viewFullName)
	query := "CREATE VIEW " + viewFullName + " AS " + queryPre
	err := DB.DB.Exec(query).Error
	return err
}

func addSchemaToTablesInQuery(query string) (string, []string) {
	schemas := []string{}

	if config.Config.Database.Connection == "postgres" {
		re := regexp.MustCompile(`(?i)(FROM|JOIN)\s+(([\w]+)\.)?([\w]+)(\s+AS\s+[\w]+)?`)
		query = re.ReplaceAllStringFunc(query, func(match string) string {
			parts := re.FindStringSubmatch(match)
			schemaName, tableName := parts[3], parts[4]

			if schemaName == "" {
				// Fetch schema name from the database if not specified
				DB.DB.Raw("SELECT table_schema FROM information_schema.tables WHERE table_name = ? LIMIT 1", tableName).Scan(&schemaName)
			}

			// Add schema to list if not already included
			if schemaName != "" && !contains(schemas, schemaName) {
				schemas = append(schemas, schemaName)
			}

			// Rebuild the match with schema name, ensuring table names are schema-qualified
			if schemaName != "" {
				return fmt.Sprintf("%s %s.%s", parts[1], schemaName, tableName)
			}
			return match
		})
	}

	return query, schemas
}

// contains checks if a string is present in a slice.
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
