package DBSchema

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/models"
)

func GetColumnsFromSQLlTable(db *sql.DB, dbTable string, hiddenColumns []string) (*map[string]map[string]string, error) {

	// Store colum as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// Select columnd data from INFORMATION_SCHEMA

	var pkColumn models.PKColumn

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable + "' AND table_schema = '" + config.Config.Database.Database + "'"

	if config.Config.Database.Connection == "mssql" {

		DB.DB.Raw("SELECT COLUMN_NAME FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '" + dbTable + "' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable + "'"
	} else if config.Config.Database.Connection == "postgres" {

		DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as \"COLUMN_NAME\" FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database, dbTable)
	} else if config.Config.Database.Connection == "oracle" {

		DB.DB.Raw(fmt.Sprintf("SELECT COLUMN_NAME FROM all_cons_columns WHERE constraint_name = (SELECT constraint_name FROM user_constraints WHERE table_name = '%s' AND CONSTRAINT_TYPE = '%s')", dbTable, "P")).Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, DATA_TYPE, (CASE WHEN NULLABLE = 'Y' THEN 'YES' ELSE 'NO' END) AS IS_NULLABLE FROM ALL_TAB_COLUMNS WHERE  OWNER = '%s' AND TABLE_NAME = '%s' ORDER  BY COLUMN_ID ASC", config.Config.Database.UserName, dbTable)
	}

	if Debug {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return nil, err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return nil, errors.New("No results returned for table")
	}

	for rows.Next() {
		var column string
		var columnKey string
		var dataType string
		var nullable string

		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}

		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns {
			if hiddenColumn == column {
				isHidden = true
			}
		}
		if isHidden == false {
			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {

				if pkColumn.ColumnName == column {
					columnKey = "PRI"
				}
			}

			columnDataTypes[column] = map[string]string{"value": dataType, "nullable": nullable, "primary": columnKey}
		}

	}

	return &columnDataTypes, err
}

func GetColumns(db *sql.DB, dbTable string, hiddenColumns []string) (string, error) {

	// Store colum as map of maps
	columns := ""
	// Select columnd data from INFORMATION_SCHEMA

	var pkColumn models.PKColumn

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable + "' AND table_schema = '" + config.Config.Database.Database + "'"

	if config.Config.Database.Connection == "mssql" {

		DB.DB.Raw("SELECT COLUMN_NAME FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '" + dbTable + "' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable + "'"
	} else if config.Config.Database.Connection == "postgres" {

		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database, dbTable)
	} else if config.Config.Database.Connection == "oracle" {

		DB.DB.Raw(fmt.Sprintf("SELECT COLUMN_NAME FROM all_cons_columns WHERE constraint_name = (SELECT constraint_name FROM user_constraints WHERE UPPER(table_name) = UPPER('%s') AND CONSTRAINT_TYPE = '%s')", dbTable, "P")).Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, DATA_TYPE, (CASE WHEN NULLABLE = 'Y' THEN 'YES' ELSE 'NO' END) AS IS_NULLABLE FROM ALL_TAB_COLUMNS WHERE  OWNER = '%s' AND TABLE_NAME = '%s' ORDER  BY COLUMN_ID ASC", config.Config.Database.UserName, dbTable)
	}

	if Debug {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return "", err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return "", errors.New("No results returned for table")
	}

	for rows.Next() {
		var column string
		var columnKey string
		var dataType string
		var nullable string
		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}

		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns {
			if hiddenColumn == column {
				isHidden = true
			}
		}
		if isHidden == false {
			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {
				if pkColumn.ColumnName == column {
					columnKey = "PRI"
				}
			}

			if columns == "" {
				columns = columns + "\"" + column + "\""
			} else {
				columns = columns + ", " + "\"" + column + "\""
			}

		}

	}

	return columns, err
}
func GetColumnsWithMeta(db *sql.DB, dbTable string, hiddenColumns []string) ([]map[string]string, error) {

	// Store colum as map of maps
	columns := []map[string]string{}
	// Select columnd data from INFORMATION_SCHEMA

	var pkColumn models.PKColumn

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable + "' AND table_schema = '" + config.Config.Database.Database + "'"

	if config.Config.Database.Connection == "mssql" {

		DB.DB.Raw("SELECT COLUMN_NAME FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '" + dbTable + "' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable + "'"
	} else if config.Config.Database.Connection == "postgres" {

		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database, dbTable)
	} else if config.Config.Database.Connection == "oracle" {

		DB.DB.Raw(fmt.Sprintf("SELECT COLUMN_NAME FROM all_cons_columns WHERE constraint_name = (SELECT constraint_name FROM user_constraints WHERE UPPER(table_name) = UPPER('%s') AND CONSTRAINT_TYPE = '%s')", dbTable, "P")).Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, DATA_TYPE, (CASE WHEN NULLABLE = 'Y' THEN 'YES' ELSE 'NO' END) AS IS_NULLABLE FROM ALL_TAB_COLUMNS WHERE  OWNER = '%s' AND TABLE_NAME = '%s' ORDER  BY COLUMN_ID ASC", config.Config.Database.UserName, dbTable)
	}

	if Debug {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return []map[string]string{}, err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return []map[string]string{}, errors.New("No results returned for table")
	}

	for rows.Next() {
		var column string
		var columnKey string
		var dataType string
		var nullable string
		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}

		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns {
			if hiddenColumn == column {
				isHidden = true
			}
		}
		if isHidden == false {
			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {
				if pkColumn.ColumnName == column {
					columnKey = "PRI"
				}
			}

			newColumn := map[string]string{}

			newColumn["column"] = column
			newColumn["nullable"] = nullable
			newColumn["dataType"] = dataType

			columns = append(columns, newColumn)

		}

	}

	return columns, err
}

func GetOnlyOneField(db *sql.DB, dbTable string, oneField string) (*map[string]map[string]string, error) {

	// Store colum as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// Select columnd data from INFORMATION_SCHEMA
	var pkColumn models.PKColumn

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable + "' AND table_schema = '" + config.Config.Database.Database + "'"

	if config.Config.Database.Connection == "mssql" {

		DB.DB.Raw("SELECT COLUMN_NAME FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '" + dbTable + "' AND CONSTRAINT_NAME LIKE 'PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable + "'"
	} else if config.Config.Database.Connection == "postgres" {

		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database, dbTable)
	} else if config.Config.Database.Connection == "oracle" {

		DB.DB.Raw(fmt.Sprintf("SELECT COLUMN_NAME FROM all_cons_columns WHERE constraint_name = (SELECT constraint_name FROM user_constraints WHERE UPPER(table_name) = UPPER('%s') AND CONSTRAINT_TYPE = '%s')", dbTable, "P")).Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, DATA_TYPE, (CASE WHEN NULLABLE = 'Y' THEN 'YES' ELSE 'NO' END) AS IS_NULLABLE FROM ALL_TAB_COLUMNS WHERE  OWNER = '%s' AND TABLE_NAME = '%s' ORDER  BY COLUMN_ID ASC", config.Config.Database.UserName, dbTable)
	}

	if Debug {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery)

	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return nil, err
	}
	if rows != nil {
		defer rows.Close()
	} else {
		return nil, errors.New("No results returned for table")
	}

	for rows.Next() {
		var column string
		var columnKey string
		var dataType string
		var nullable string
		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}

		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" || config.Config.Database.Connection == "oracle" {
			if pkColumn.ColumnName == column {
				columnKey = "PRI"
			}
		}

		//	if oneField == column {
		columnDataTypes[column] = map[string]string{"value": dataType, "nullable": nullable, "primary": columnKey}
		//	}

	}

	return &columnDataTypes, err
}
