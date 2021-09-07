package DBSchema

import (
	"encoding/json"
	"github.com/lambda-platform/lambda/models"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"io/ioutil"
	"fmt"
)

func GetDBSchema() models.DBSCHEMA {
	tables := Tables()

	table_metas := make(map[string][]models.TableMeta, 0)

	for _, table := range tables["tables"] {
		table_metas_ := TableMetas(table)
		table_metas[table] = table_metas_
	}

	for _, table := range tables["views"] {
		table_metas_ := TableMetas(table)
		table_metas[table] = table_metas_
	}

	vb_schemas := models.DBSCHEMA{
		tables["tables"],
		tables["views"],
		table_metas,
	}

	file, _ := json.MarshalIndent(vb_schemas, "", " ")

	_ = ioutil.WriteFile("models/db_schema.json", file, 0755)

	return vb_schemas
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
	} else if config.Config.Database.Connection == "postgres"  {
		rows, _ := DB_.Query("SELECT tablename, concat('TABLE') as tabletype FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema' union SELECT table_name as tablename, concat('VIEW') as tabletype FROM information_schema.views where table_schema not in ('information_schema', 'pg_catalog')")
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
	table_metas := make([]models.TableMeta, 0)
	DB_ := DB.DBConnection()

	if config.Config.Database.Connection == "mssql" {

		var pkColumn models.PKColumn
		DB.DB.Raw("SELECT COLUMN_NAME FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_NAME LIKE '" + tableName + "' AND CONSTRAINT_NAME LIKE '%PK%'").Scan(&pkColumn)

		table_metas_ms := []models.MSTableMata{}
		DB.DB.Raw("SELECT * FROM " + config.Config.Database.Database + ".INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '" + tableName + "'").Scan(&table_metas_ms)

		for _, column := range table_metas_ms {
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

			table_metas = append(table_metas, models.TableMeta{
				Model:  column.ColumnName,
				Title:  column.ColumnName,
				DbType: dataType,
				Table:  tableName,
				Key:    key,
				Extra:  extra,
			})
		}

	}else if config.Config.Database.Connection == "postgres" {
		pkColumn := ""
		rowPK := DB.DB.Raw(fmt.Sprintf("SELECT k.COLUMN_NAME as pkColumn FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", tableName, config.Config.Database.Database, "'%_pkey'")).Row()
		rowPK.Scan(&pkColumn)
		//	fmt.Println(fmt.Sprintf("SELECT k.column_name FROM information_schema.key_column_usage k   WHERE k.table_name = '%s' AND k.table_catalog ='%s'AND k.constraint_name LIKE %s", tableName, config.Config.Database.Database, "'%_pkey'"))

		Enums := []models.PostgresEnum{}
//
		DB.DB.Raw("SELECT pg_type.typname FROM pg_type JOIN pg_enum ON pg_enum.enumtypid = pg_type.oid  GROUP BY  pg_type.typname").Scan(&Enums)

		rows, _ := DB.DB.Raw(fmt.Sprintf("SELECT udt_name as DATA_TYPE, COLUMN_NAME, IS_NULLABLE FROM information_schema.columns WHERE udt_catalog = '%s' AND table_name   = '%s'", config.Config.Database.Database,  tableName)).Rows()

		defer rows.Close()
		for rows.Next() {
			var dataType string
			var nullable string
			columnName := ""
			rows.Scan(&dataType, &columnName, &nullable)

			key := ""
			extra := ""




			if columnName == pkColumn {
				key = "PRI"
				extra = "auto_increment"
			}


			if dataType == "varchar" {
				dataType = "varchar"
			} else if dataType == "ntext" {
				dataType = "text"
			} else if dataType == "int8" {
				dataType = "int"
			} else if dataType == "int4" {
				dataType = "int"
			} else if dataType == "timestamptz" {
				dataType = "timestamp"
			}

			for _, enum :=range Enums{
				if(enum.Typname == dataType) {
					dataType = "varchar"
				}
			}

			table_metas = append(table_metas, models.TableMeta{
				Model:  columnName,
				Title:  columnName,
				DbType: dataType,
				Table:  tableName,
				Key:    key,
				Extra:  extra,
				Nullable:  nullable,
			})
		}


	} else {
		columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '" + tableName+"' AND table_schema = '" + config.Config.Database.Database+"'"

		columns, db_error := DB_.Query(columnDataTypeQuery)

		if db_error == nil {
			for columns.Next() {
				var Field, Type, Null, Key, Extra string
				columns.Scan(&Field, &Key, &Type, &Null)

				table_metas = append(table_metas, models.TableMeta{
					Model:  Field,
					Title:  Field,
					DbType: Type,
					Table:  tableName,
					Key:    Key,
					Extra:  Extra,
					Nullable:  Null,
				})
			}
		}
	}

	return table_metas

}
func GenerateSchemaForCloud() models.DBSCHEMA {
	tables := TablesForCloud()

	table_metas := make(map[string][]models.TableMeta, 0)

	for _, table := range tables["tables"] {
		if(table != "vb_schemas" && table != "vb_schemas_admin" && table != "krud"){
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
	}

	file, _ := json.MarshalIndent(vb_schemas, "", " ")
	_ = ioutil.WriteFile("app/models/db_schema.json", file, 0755)

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
				if(TABLE_NAME != "vb_schemas" && TABLE_NAME != "vb_schemas_admin" && TABLE_NAME != "krud" && TABLE_NAME != "password_resets") {
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
	} else if config.Config.Database.Connection == "postgres"  {
		rows, _ := DB_.Query("SELECT tablename, concat('TABLE') as tabletype FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema' union SELECT table_name as tablename, concat('VIEW') as tabletype FROM information_schema.views where table_schema not in ('information_schema', 'pg_catalog')")
		for rows.Next() {
			var tableName, tableType string
			rows.Scan(&tableName, &tableType)
			if tableType != "VIEW" {
				if(tableName != "vb_schemas" && tableName != "vb_schemas_admin" && tableName != "krud" && tableName != "password_resets") {
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
				if(tableName != "vb_schemas" && tableName != "vb_schemas_admin" && tableName != "krud" && tableName != "password_resets") {
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