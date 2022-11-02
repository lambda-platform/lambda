package DBSchema

import (
	"fmt"
	"go/format"
	"sort"
	"strings"
)

const (
	golangByteArray = "[]byte"
	//gureguNullInt    = "null.Int"
	gureguNullInt   = "*int"
	sqlNullInt      = "sql.NullInt64"
	golangInt       = "int"
	golangBool      = "bool"
	golangInt64     = "int64"
	gureguNullFloat = "*float32"
	sqlNullFloat    = "sql.NullFloat64"
	golangFloat     = "float"
	golangFloat32   = "float32"
	golangFloat64   = "float64"
	//gureguNullString = "null.String"
	gureguNullString = "*string"
	sqlNullString    = "*string"
	//gureguNullTime   = "null.Time"
	gureguNullTime = "*time.Time"
	golangTime     = "time.Time"
	date           = "DB.Date"
	dateNull       = "*DB.Date"
)

func GenerateOnlyStruct(columnTypes map[string]map[string]string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string) ([]byte, error) {
	var dbTypes string

	dbTypes, _, _ = generateStructTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	src := fmt.Sprintf("\n  \ntype %s %s %s} %s",
		structName,
		dbTypes,
		extraColumns, extraStucts)
	if gormAnnotation == true {
		tableNameFunc := "" +
			"func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") TableName() string {\n" +
			"	return \"" + tableName + "\"" +
			"}"
		src = fmt.Sprintf("%s\n%s", src, tableNameFunc)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {

		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}
func GenerateWithImports(otherPackage string, columnTypes map[string]map[string]string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string, virtualColums string) ([]byte, error) {
	var dbTypes string

	dbTypes, _, _ = generateStructTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	importTime := "import (\n\"time\"\n\"github.com/lambda-platform/lambda/DB\"\n\"gorm.io/gorm\") \n var _ = time.Time{}  \n var _ = DB.Date{}  \nvar _ = gorm.DB{} \n "
	src := fmt.Sprintf("package %s\n %s\n %s\n \ntype %s %s %s %s} %s",
		pkgName,
		otherPackage,
		importTime,
		structName,
		dbTypes,
		extraColumns, virtualColums, extraStucts)
	if gormAnnotation == true {
		tableNameFunc := "//  TableName sets the insert table name for this struct type\n " +
			"func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") TableName() string {\n" +
			"	return \"" + tableName + "\"" +
			"}"
		src = fmt.Sprintf("%s\n%s", src, tableNameFunc)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}
func GenerateWithImportsNoTime(otherPackage string, columnTypes map[string]map[string]string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string) ([]byte, error) {
	var dbTypes string

	dbTypes, timeFound := generateStructTypesNoTime(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	var _ = timeFound
	importTime := ""
	//if time_found {
	//	importTime = "import \"time\" \n var _ = time.Time{}  \n"
	//}
	src := fmt.Sprintf("package %s\n %s\n %s\n \ntype %s %s %s} %s",
		pkgName,
		otherPackage,
		importTime,
		structName,
		dbTypes,
		extraColumns, extraStucts)
	if gormAnnotation == true {
		tableNameFunc := "//  TableName sets the insert table name for this struct type\n " +
			"func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") TableName() string {\n" +
			"	return \"" + tableName + "\"" +
			"}"
		src = fmt.Sprintf("%s\n%s", src, tableNameFunc)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

func sqlTypeToGoType(columnType string, nullable bool, gureguTypes bool) string {

	switch columnType {
	case "tinyint", "int", "smallint", "mediumint", "int8", "int4", "int2", "year", "NUMBER", "INT", "INTEGER":
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
	case "char", "enum", "varchar", "nvarchar", "longtext", "mediumtext", "text", "ntext", "tinytext", "geometry", "uuid", "bpchar", "VARCHAR2", "CLOB", "LONG":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case "time", "timestamp", "datetimeoffset", "timestamptz", "TIMESTAMP(6) WITH TIME ZONE", "TIMESTAMP", "TIMESTAMP(6)":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return golangTime
	case "datetime", "date", "DATE":
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

func generateStructTypes(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) (string, bool, bool) {

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
		columnType := obj[key]
		nullable := false
		if columnType["nullable"] == "YES" {
			nullable = true
		}
		if columnType["value"] == "timestamp" || columnType["value"] == "timestamptz" || columnType["value"] == "datetime" || columnType["value"] == "year" || columnType["value"] == "time" {

			//if key == "created_at" ||  key == "updated_at" ||  key == "deleted_at"{
			time_found = true
			//}
			//else {
			//	columnType["value"] = "text"
			//}

		}
		if columnType["value"] == "date" {
			date_found = true
		}
		primary := ""
		if columnType["primary"] == "PRI" {
			primary = ";primaryKey;autoIncrement"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGoType(columnType["value"], nullable, gureguTypes)

		fieldName := FmtFieldName(strings.ToLower(StringifyFirstChar(key)))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, ""))
		}
		if fieldName == "DeletedAt" || fieldName == "DELETEDAT" {
			valueType = "gorm.DeletedAt"
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
func generateStructTypesNoTime(obj map[string]map[string]string, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) (string, bool) {

	structure := "struct {"
	time_found := false
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
		if columnType["value"] == "timestamp" || columnType["value"] == "timestamptz" || columnType["value"] == "datetime" || columnType["value"] == "date" || columnType["value"] == "year" || columnType["value"] == "time" {

			columnType["value"] = "text"

		}

		primary := ""
		if columnType["primary"] == "PRI" {
			primary = ";primaryKey;autoIncrement"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGoType(columnType["value"], nullable, gureguTypes)

		fieldName := FmtFieldName(StringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, ""))
		}
		if fieldName == "DeletedAt" || fieldName == "DELETEDAT" {
			valueType = "gorm.DeletedAt"
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
