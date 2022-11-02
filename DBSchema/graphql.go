package DBSchema

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"sort"
	"strings"
)

const (
	gqlNullInt    = "Int"
	gqlInt        = "Int!"
	gqlNullFloat  = "Float"
	gqlFloat      = "Float!"
	gqlNullString = "String"
	gqlString     = "String!"
	gqlNullTime   = "Time"
	gqlTime       = "Time!"
	dbNullDate    = "Date"
	dbDate        = "Date!"
)

func GenerateGrapql(columnTypes map[string]map[string]string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string, Subs []string, isInpute bool) ([]byte, error) {

	dbTypes := generateQraphqlTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)
	//if tableName == "aa_sudalsan_sain_turshilga" {
	//	fmt.Println(columnTypes)
	//	fmt.Println(dbTypes)
	//}
	subStchemas := ""

	for _, sub := range Subs {
		subStchemas = subStchemas + "\n    " + sub + ":[" + strcase.ToCamel(strings.ToLower(sub)) + "!]"
	}

	typeSchema := "type"
	if isInpute {
		typeSchema = "input"
		structName = structName + "Input"
	}
	src := fmt.Sprintf("%s %s %s %s %s \n} %s",
		typeSchema,
		structName,
		dbTypes,
		extraColumns, subStchemas, extraStucts)

	return []byte(src), nil
}
func GenerateGrapqlOrder(columnTypes map[string]map[string]string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string) ([]byte, error) {

	dbTypes := generateQraphqlTypesOrder(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	src := fmt.Sprintf("\n  \ninput %s %s %s \n} %s",
		structName,
		dbTypes,
		extraColumns, extraStucts)

	return []byte(src), nil
}

func generateQraphqlTypes(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {

	structure := " {"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//fmt.Println(key)
		columnType := obj[key]
		nullable := false
		if columnType["nullable"] == "YES" {
			nullable = true
		}

		primary := ""
		if columnType["primary"] == "PRI" {
			primary = ";primaryKey;autoIncrement"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGraphyType(columnType["value"], nullable, gureguTypes)

		if columnType["primary"] == "PRI" {
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
		if fieldName == "DeletedAt" || fieldName == "deleted_at" || fieldName == "DELETED_AT" {
			valueType = "GormDeletedAt"
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
func generateQraphqlTypesOrder(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {

	structure := " {"

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//fmt.Println(key)
		columnType := obj[key]

		primary := ""
		if columnType["primary"] == "PRI" {
			primary = ";primaryKey;autoIncrement"
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
func sqlTypeToGraphyType(columnType string, nullable bool, gureguTypes bool) string {
	switch columnType {
	case "tinyint", "int", "smallint", "mediumint", "int8", "int4", "int2", "year", "NUMBER", "INT", "INTEGER":
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
	case "char", "enum", "varchar", "nvarchar", "longtext", "mediumtext", "text", "ntext", "tinytext", "geometry", "uuid", "bpchar", "VARCHAR2", "CLOB", "LONG":
		if nullable {
			if gureguTypes {
				return gqlNullString
			}
			return gqlNullString
		}
		return gqlString
	case "time", "timestamp", "datetimeoffset", "timestamptz", "TIMESTAMP(6) WITH TIME ZONE", "TIMESTAMP", "TIMESTAMP(6)":
		if nullable && gureguTypes {
			return gqlNullTime
		}
		return gqlTime
	case "datetime", "date", "DATE":
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
