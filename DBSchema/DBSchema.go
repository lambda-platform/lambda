package DBSchema

import (
	"encoding/json"
	"fmt"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/models"
	"github.com/lambda-platform/lambda/utils"
	"os"
	"strings"
)

var Enums []models.PostgresEnum

func GetDBSchema() models.DBSCHEMA {
	tables := Tables()

	tableMetas := make(map[string][]models.TableMeta, 0)

	if config.Config.Database.Connection == "postgres" {
		DB.DB.Raw("SELECT pg_type.typname FROM pg_type JOIN pg_enum ON pg_enum.enumtypid = pg_type.oid  GROUP BY  pg_type.typname").Scan(&Enums)
	}

	for _, table := range tables["tables"] {
		fmt.Println(config.LambdaConfig.IgnoreTables)
		if !utils.StringInSlice(table, config.LambdaConfig.IgnoreTables) {
			tableMetas[table] = TableMetas(table)
		}
	}

	for _, table := range tables["views"] {
		if !utils.StringInSlice(table, config.LambdaConfig.IgnoreTables) {
			tableMetas[table] = TableMetas(table)
		}
	}

	vbSchemas := models.DBSCHEMA{
		TableList: tables["tables"],
		ViewList:  tables["views"],
		TableMeta: tableMetas,
	}

	file, _ := json.MarshalIndent(vbSchemas, "", " ")

	_ = os.WriteFile("lambda/db_schema.json", file, 0755)

	return vbSchemas
}
func GetDBSchemaWithTargets(tables map[string][]string) models.DBSCHEMA {

	tableMetas := make(map[string][]models.TableMeta, 0)

	if config.Config.Database.Connection == "postgres" {
		DB.DB.Raw("SELECT pg_type.typname FROM pg_type JOIN pg_enum ON pg_enum.enumtypid = pg_type.oid  GROUP BY  pg_type.typname").Scan(&Enums)
	}

	for _, table := range tables["tables"] {
		tableMetas[table] = TableMetas(table)
	}

	for _, table := range tables["views"] {
		tableMetas[table] = TableMetas(table)
	}

	vbSchemas := models.DBSCHEMA{
		TableList: tables["tables"],
		ViewList:  tables["views"],
		TableMeta: tableMetas,
	}

	file, _ := json.MarshalIndent(vbSchemas, "", " ")

	_ = os.WriteFile("lambda/db_schema.json", file, 0755)

	return vbSchemas
}

func Tables() map[string][]string {
	tables := make([]string, 0)
	views := make([]string, 0)

	//var dbTables []dbTable
	//DB.Raw("SHOW FULL TABLES").Scan(&dbTables)

	DB_ := DB.DBConnection()
	if config.Config.Database.Connection == "mssql" {
		rows, _ := DB_.Query("SELECT TABLE_NAME, TABLE_TYPE FROM INFORMATION_SCHEMA.TABLES ORDER BY TABLE_NAME")

		for rows.Next() {
			var TABLE_NAME, TABLE_TYPE string
			rows.Scan(&TABLE_NAME, &TABLE_TYPE)

			if TABLE_TYPE != "VIEW" {
				tables = append(tables, TABLE_NAME)
			} else {
				views = append(views, TABLE_NAME)
			}
		}
		result := map[string][]string{}

		result["tables"] = tables
		result["views"] = views

		return result
	} else if config.Config.Database.Connection == "postgres" {
		rows, _ := DB_.Query("SELECT tablename, concat('TABLE') as tabletype FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema' union SELECT table_name as tablename, concat('VIEW') as tabletype FROM information_schema.views where table_schema not in ('information_schema', 'pg_catalog') ORDER BY tablename")
		for rows.Next() {
			var tableName, tableType string
			rows.Scan(&tableName, &tableType)
			if tableType != "VIEW" {
				tables = append(tables, tableName)
			} else {
				views = append(views, tableName)
			}
		}
		result := map[string][]string{}

		result["tables"] = tables
		result["views"] = views

		return result
	} else if config.Config.Database.Connection == "oracle" {
		rows, _ := DB_.Query(fmt.Sprintf("SELECT TABLE_NAME  FROM All_TABLES WHERE OWNER = '%s' ORDER BY TABLE_NAME", config.Config.Database.UserName))

		for rows.Next() {
			var tableName string
			rows.Scan(&tableName)
			tables = append(tables, tableName)
		}

		rows, _ = DB_.Query(fmt.Sprintf("SELECT VIEW_NAME FROM All_VIEWS WHERE OWNER = '%s' ORDER BY VIEW_NAME", config.Config.Database.UserName))
		for rows.Next() {
			var tableType string
			rows.Scan(&tableType)
			views = append(views, tableType)
		}
		result := map[string][]string{}

		result["tables"] = tables
		result["views"] = views

		return result
	} else {
		rows, _ := DB_.Query("SHOW FULL TABLES")
		for rows.Next() {
			var tableName, tableType string
			rows.Scan(&tableName, &tableType)
			if tableType == "BASE TABLE" {
				tables = append(tables, tableName)
			} else {
				views = append(views, tableName)
			}
		}
		result := map[string][]string{}

		result["tables"] = tables
		result["views"] = views

		return result
	}

}

func TableMetas(tableName string) []models.TableMeta {

	tableMetas := make([]models.TableMeta, 0)

	if config.Config.Database.Connection == "mssql" {

		var pkColumn models.PKColumn
		DB.DB.Raw("SELECT COLUMN_NAME FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '" + tableName + "' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		var currentTableMetas []models.MSTableMata

		query := fmt.Sprintf(`
        SELECT 
            C.COLUMN_NAME, 
            C.DATA_TYPE, 
            DC.definition AS DEFAULT_VALUE
        FROM 
            %s.INFORMATION_SCHEMA.COLUMNS C
        LEFT JOIN 
            sys.columns SC ON SC.object_id = OBJECT_ID(C.TABLE_SCHEMA + '.' + C.TABLE_NAME) AND SC.name = C.COLUMN_NAME
        LEFT JOIN 
            sys.default_constraints DC ON DC.parent_object_id = SC.object_id AND DC.parent_column_id = SC.column_id
        WHERE 
            C.TABLE_NAME = '%s'`, config.Config.Database.Database, tableName)

		DB.DB.Raw(query).Scan(&currentTableMetas)
		for _, column := range currentTableMetas {
			key := ""
			extra := ""

			if column.ColumnName == pkColumn.ColumnName {
				key = "PRI"
				extra = "auto_increment"
			}

			dataType := column.DataType

			if column.DataType == "nvarchar" {
				dataType = "varchar"
			} else if column.DataType == "ntext" {
				dataType = "text"
			}

			dataType = IsSecureField(tableName, column.ColumnName, dataType)

			tableMetas = append(tableMetas, models.TableMeta{
				Model:        column.ColumnName,
				Title:        column.ColumnName,
				DbType:       dataType,
				Table:        tableName,
				Key:          key,
				Extra:        extra,
				DefaultValue: column.DefaultValue,
			})
		}

	} else if config.Config.Database.Connection == "postgres" {

		var currentTableMetas []models.PostgresTableMata

		DB.DB.Raw(fmt.Sprintf("SELECT column_name, udt_name, is_nullable, is_identity, column_default, numeric_scale, table_schema FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s' ORDER BY ORDINAL_POSITION", config.Config.Database.Database, tableName)).Scan(&currentTableMetas)

		for _, column := range currentTableMetas {

			key := ""
			extra := ""
			scale := ""

			if column.IsIdentity == "YES" {
				key = "PRI"
				extra = "auto_increment"
			} else if column.ColumnDefault != nil {
				if strings.Contains(*column.ColumnDefault, "nextval(") {
					key = "PRI"
					extra = "auto_increment"
				}
				if strings.Contains(*column.ColumnDefault, "SYS_GUID") {
					key = "PRI"
					extra = "auto_increment"
				}
				if strings.Contains(*column.ColumnDefault, "gen_random_uuid") {
					key = "PRI"
					extra = "auto_increment"
				}
			}
			if column.DataType == "geometry" {
				//extra = "Point"

				//scale = ";geotype:"+column.Type
			}
			for _, enum := range Enums {
				if enum.Typname == column.DataType {
					column.DataType = "varchar"
				}
			}

			if column.NumericScale != nil {
				if *column.NumericScale >= 1 {

					scale = ";scale:2"
				}

			}
			column.DataType = IsSecureField(tableName, column.ColumnName, column.DataType)

			tableMetas = append(tableMetas, models.TableMeta{
				Model:        column.ColumnName,
				Title:        column.ColumnName,
				DbType:       column.DataType,
				Table:        tableName,
				Key:          key,
				Extra:        extra,
				Scale:        scale,
				Nullable:     column.ISNullAble,
				DefaultValue: column.ColumnDefault,
				TableSchema:  column.TableSchema,
			})
		}

	} else if config.Config.Database.Connection == "oracle" {

		var currentTableMetas []models.OracleTableMata
		//fmt.Println(fmt.Sprintf("SELECT  COLUMN_NAME, DATA_TYPE, NULLABLE, IDENTITY_COLUMN, DATA_DEFAULT FROM ALL_TAB_COLUMNS WHERE  OWNER = '%s' AND TABLE_NAME = '%s' ORDER  BY COLUMN_ID ASC", config.Config.Database.UserName, tableName))
		DB.DB.Raw(fmt.Sprintf("SELECT  COLUMN_NAME, DATA_TYPE, NULLABLE, IDENTITY_COLUMN, DATA_DEFAULT FROM ALL_TAB_COLUMNS WHERE  OWNER = '%s' AND TABLE_NAME = '%s' ORDER  BY COLUMN_ID ASC", config.Config.Database.UserName, tableName)).Scan(&currentTableMetas)

		for _, column := range currentTableMetas {

			key := ""
			extra := ""
			Nullable := "YES"
			scale := ""

			if column.IdentityColumn == "YES" {
				key = "PRI"
				extra = "auto_increment"
			} else if column.DataDefault != nil {
				if strings.Contains(*column.DataDefault, "nextval") {
					key = "PRI"
					extra = "auto_increment"
				}
				if strings.Contains(*column.DataDefault, "SYS_GUID") {
					key = "PRI"
					extra = "auto_increment"
				}
			}
			if column.NullAble == "N" {
				Nullable = "NO"

			}

			dataType := column.DataType

			if column.DataType == "VARCHAR2" {
				dataType = "varchar"
			} else if column.DataType == "LONG" {
				dataType = "text"
			} else if column.DataType == "NUMBER" {

				if column.ColumnName == "ID" || strings.HasPrefix(column.ColumnName, "ID") || strings.HasSuffix(column.ColumnName, "ID") {
					dataType = "int"
				}

			}

			//if column.DataType == "BLOB" {
			//	//scale = ";serializer:gob"
			//}
			dataType = IsSecureField(tableName, column.ColumnName, dataType)
			tableMetas = append(tableMetas, models.TableMeta{
				Model:        column.ColumnName,
				Title:        column.ColumnName,
				DbType:       dataType,
				Table:        tableName,
				Key:          key,
				Extra:        extra,
				Nullable:     Nullable,
				Scale:        scale,
				DefaultValue: column.DataDefault,
			})
		}

	} else {

		currentTableMetas := []models.MySQLTableMata{}
		DB.DB.Raw(fmt.Sprintf("SELECT column_name as column_name, column_key as column_key, data_type as data_type, is_nullable as is_nullable, COLUMN_DEFAULT as default_value FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%s' AND table_schema = '%s' ORDER BY ORDINAL_POSITION", tableName, config.Config.Database.Database)).Scan(&currentTableMetas)

		for _, column := range currentTableMetas {

			key := ""
			extra := ""

			if column.ColumnKey == "PRI" {
				key = "PRI"
				extra = "auto_increment"
			}

			column.DataType = IsSecureField(tableName, column.ColumnName, column.DataType)
			tableMetas = append(tableMetas, models.TableMeta{
				Model:        column.ColumnName,
				Title:        column.ColumnName,
				DbType:       column.DataType,
				Table:        tableName,
				Key:          key,
				Extra:        extra,
				Nullable:     column.ISNullAble,
				DefaultValue: column.DefaultValue,
			})
		}

	}

	return tableMetas

}
func IsSecureField(table, column, dataType string) string {
	if len(config.LambdaConfig.SecureFields) >= 1 {
		for _, field := range config.LambdaConfig.SecureFields {
			if field.Table == table && field.Column == column {
				return "secure"
			}
		}
	}

	return dataType
}
func GenerateSchemaForCloud() models.DBSCHEMA {
	tables := TablesForCloud()

	table_metas := make(map[string][]models.TableMeta, 0)

	for _, table := range tables["tables"] {
		if table != "vb_schemas" && table != "vb_schemas_admin" && table != "krud" {
			table_metas_ := TableMetas(table)
			table_metas[table] = table_metas_
		}

	}

	for _, table := range tables["views"] {
		table_metas_ := TableMetas(table)
		table_metas[table] = table_metas_
	}

	vb_schemas := models.DBSCHEMA{
		tables["tables"],
		tables["views"],
		table_metas,
		0,
		"",
	}

	file, _ := json.MarshalIndent(vb_schemas, "", " ")
	_ = os.WriteFile("app/models/db_schema.json", file, 0755)

	return vb_schemas
}

func TablesForCloud() map[string][]string {
	tables := make([]string, 0)
	views := make([]string, 0)

	//var dbTables []dbTable
	//DB.Raw("SHOW FULL TABLES").Scan(&dbTables)

	DB_ := DB.DBConnection()
	if config.Config.Database.Connection == "mssql" {
		rows, _ := DB_.Query("SELECT TABLE_NAME, TABLE_TYPE FROM INFORMATION_SCHEMA.TABLES ORDER BY TABLE_NAME")

		for rows.Next() {
			var TABLE_NAME, TABLE_TYPE string
			rows.Scan(&TABLE_NAME, &TABLE_TYPE)

			if TABLE_TYPE != "VIEW" {
				if TABLE_NAME != "vb_schemas" && TABLE_NAME != "vb_schemas_admin" && TABLE_NAME != "krud" && TABLE_NAME != "password_resets" {
					tables = append(tables, TABLE_NAME)
				}
			} else {
				views = append(views, TABLE_NAME)
			}
		}
		result := map[string][]string{}

		result["tables"] = tables
		result["views"] = views

		return result
	} else if config.Config.Database.Connection == "postgres" {
		rows, _ := DB_.Query("SELECT tablename, concat('TABLE') as tabletype FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema' union SELECT table_name as tablename, concat('VIEW') as tabletype FROM information_schema.views where table_schema not in ('information_schema', 'pg_catalog')  ORDER BY tablename")
		for rows.Next() {
			var tableName, tableType string
			rows.Scan(&tableName, &tableType)
			if tableType != "VIEW" {
				if tableName != "vb_schemas" && tableName != "vb_schemas_admin" && tableName != "krud" && tableName != "password_resets" {
					tables = append(tables, tableName)
				}
			} else {
				views = append(views, tableName)
			}
		}
		result := map[string][]string{}

		result["tables"] = tables
		result["views"] = views

		return result
	} else {
		rows, _ := DB_.Query("SHOW FULL TABLES")
		for rows.Next() {
			var tableName, tableType string
			rows.Scan(&tableName, &tableType)
			if tableType == "BASE TABLE" {
				if tableName != "vb_schemas" && tableName != "vb_schemas_admin" && tableName != "krud" && tableName != "password_resets" {
					tables = append(tables, tableName)
				}
			} else {
				views = append(views, tableName)
			}
		}
		result := map[string][]string{}

		result["tables"] = tables
		result["views"] = views

		return result
	}

}
