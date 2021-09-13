package DBSchema

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/models"
	"github.com/lambda-platform/lambda/config"
	"sort"
	"strings"
)

func GetColumnsFromSQLlTable(db *sql.DB, dbTable string, hiddenColumns []string) (*map[string]map[string]string, error) {

	// Store colum as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// Select columnd data from INFORMATION_SCHEMA

	var pkColumn models.PKColumn

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable+"' AND table_schema = '" + config.Config.Database.Database+"'"

	if config.Config.Database.Connection == "mssql"{

		DB.DB.Raw("SELECT COLUMN_NAME FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '"+dbTable+"' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable+"'"
	} else if config.Config.Database.Connection == "postgres"{


		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database,  dbTable)
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
		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}

		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns{
			if hiddenColumn == column{
				isHidden = true
			}
		}
		if isHidden == false{
			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
				if pkColumn.ColumnName == column{
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

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable+"' AND table_schema = '" + config.Config.Database.Database+"'"

	if config.Config.Database.Connection == "mssql"{

		DB.DB.Raw("SELECT COLUMN_NAME FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '"+dbTable+"' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable+"'"
	} else if config.Config.Database.Connection == "postgres"{


		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database,  dbTable)
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
		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}

		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns{
			if hiddenColumn == column{
				isHidden = true
			}
		}
		if isHidden == false{
			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
				if pkColumn.ColumnName == column{
					columnKey = "PRI"
				}
			}

			if(columns == ""){
				columns =columns + "\""+column+"\""
			} else {
				columns =  columns + ", "+ "\""+column+"\""
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

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable+"' AND table_schema = '" + config.Config.Database.Database+"'"

	if config.Config.Database.Connection == "mssql"{

		DB.DB.Raw("SELECT COLUMN_NAME FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '"+dbTable+"' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable+"'"
	} else if config.Config.Database.Connection == "postgres"{


		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database,  dbTable)
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
		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres"{
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}

		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns{
			if hiddenColumn == column{
				isHidden = true
			}
		}
		if isHidden == false{
			if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres"{
				if pkColumn.ColumnName == column{
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

	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + dbTable+"' AND table_schema = '" + config.Config.Database.Database+"'"

	if config.Config.Database.Connection == "mssql"{

		DB.DB.Raw("SELECT COLUMN_NAME FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '"+dbTable+"' AND CONSTRAINT_NAME LIKE 'PK%'").Scan(&pkColumn)

		columnDataTypeQuery = "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM "+config.Config.Database.Database+".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + dbTable+"'"
	}  else if config.Config.Database.Connection == "postgres"{


		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", dbTable, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)

		columnDataTypeQuery = fmt.Sprintf("SELECT  COLUMN_NAME, udt_name as DATA_TYPE, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database,  dbTable)
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
		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres" {
			rows.Scan(&column, &dataType, &nullable)
		} else {
			rows.Scan(&column, &columnKey, &dataType, &nullable)
		}


		if config.Config.Database.Connection == "mssql" || config.Config.Database.Connection == "postgres"{
			if pkColumn.ColumnName == column{
				columnKey = "PRI"
			}
		}

	//	if oneField == column {
			columnDataTypes[column] = map[string]string{"value": dataType, "nullable": nullable, "primary": columnKey}
	//	}


	}

	return &columnDataTypes, err
}

// Generate go struct entries for a map[string]interface{} structure
func generateMysqlTypes(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) (string, bool, bool) {

	structure := "struct {"
	time_found := false
	date_found := false
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//fmt.Println(key)
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}
		if mysqlType["value"] == "timestamp" || mysqlType["value"] == "timestamptz" || mysqlType["value"] == "datetime"  || mysqlType["value"] == "year"  || mysqlType["value"] == "time"{

			//if key == "created_at" ||  key == "updated_at" ||  key == "deleted_at"{
			time_found = true
			//}
			//else {
			//	mysqlType["value"] = "text"
			//}

		}
		if(mysqlType["value"] == "date"){
			date_found = true
		}
		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGoType(mysqlType["value"], nullable, gureguTypes)

		fieldName := FmtFieldName(StringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, ""))
		}
		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n%s %s",
				fieldName,
				valueType)
		}
	}

	return structure, time_found, date_found
}
func generateMysqlTypesNoTime(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) (string, bool) {

	structure := "struct {"
	time_found := false
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//fmt.Println(key)
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}
		if mysqlType["value"] == "timestamp" || mysqlType["value"] == "timestamptz" || mysqlType["value"] == "datetime" || mysqlType["value"] == "date"  || mysqlType["value"] == "year"  || mysqlType["value"] == "time"{

			mysqlType["value"] = "text"

		}

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGoType(mysqlType["value"], nullable, gureguTypes)

		fieldName := FmtFieldName(StringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, ""))
		}
		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n%s %s",
				fieldName,
				valueType)
		}
	}

	return structure, time_found
}
func generateQraphqlTypes(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) (string) {

	structure := " {"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//fmt.Println(key)
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGraphyType(mysqlType["value"], nullable, gureguTypes)

		if mysqlType["primary"] == "PRI" {
			valueType = "ID!"
			//primary = ""
		}

		fieldName := key
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, ""))
		}


		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n%s    %s: `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n    %s: %s",
				fieldName,
				valueType)
		}
	}

	return structure
}

func generateQraphqlTypesOrder(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) (string) {

	structure := " {"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//fmt.Println(key)
		mysqlType := obj[key]


		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
			//primary = ""
		}


		fieldName := key
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, ""))
		}


		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n%s    %s: order_by",
				fieldName,
				strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n    %s: order_by",
				fieldName)
		}
	}

	return structure
}

// sqlTypeToGoType converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func sqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint", "int8", "int4":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt
	case "bool":
		if nullable {
			if gureguTypes {
				return golangBool
			}
			return golangBool
		}
		return golangBool
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt64
	case "char", "enum", "varchar", "nvarchar", "longtext", "mediumtext", "text", "ntext",  "tinytext", "geometry":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case "datetime", "time", "timestamp", "datetimeoffset", "timestamptz":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return golangTime
	case "date":
		if nullable && gureguTypes {
			return dateNull
		}
		return date
	case "decimal", "double", "numeric":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "float", "float8", "float4", "real":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}
func sqlTypeToGraphyType(mysqlType string, nullable bool, gureguTypes bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint", "int8", "int4":
		if nullable {
			if gureguTypes {
				return gqlNullInt
			}
			return gqlNullInt
		}
		return gqlInt
	case "bigint":
		if nullable {
			if gureguTypes {
				return gqlNullInt
			}
			return gqlNullInt
		}
		return gqlInt
	case "char", "enum", "varchar", "nvarchar", "longtext", "mediumtext", "text", "ntext",  "tinytext", "geometry":
		if nullable {
			if gureguTypes {
				return gqlNullString
			}
			return gqlNullString
		}
		return gqlString
	case "datetime", "time", "timestamp", "datetimeoffset", "timestamptz":
		if nullable && gureguTypes {
			return gqlNullTime
		}
		return gqlTime
	case "date":
		if nullable && gureguTypes {
			return dbNullDate
		}
		return dbDate
	case "decimal", "double", "numeric":
		if nullable {
			if gureguTypes {
				return gqlNullFloat
			}
			return gqlNullFloat
		}
		return gqlFloat
	case "float", "float8", "float4", "real":
		if nullable {
			if gureguTypes {
				return gqlNullFloat
			}
			return gqlNullFloat
		}
		return gqlFloat
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return gqlString
	}
	return ""
}