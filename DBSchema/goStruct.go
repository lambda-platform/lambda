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

	switch {
	case TypeContains(columnType, TypeIntegers):
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt
	case TypeContains(columnType, TypeBool):
		if nullable {
			if gureguTypes {
				return golangBool
			}
			return golangBool
		}
		return golangBool
	case TypeContains(columnType, TypeBigIntegers):
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt64
	case TypeContains(columnType, TypeStrings):
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case TypeContains(columnType, TypeTimes):
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return golangTime
	case TypeContains(columnType, TypeDates):
		if nullable && gureguTypes {
			return dateNull
		}
		return date
	case TypeContains(columnType, TypeFloat64):
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case TypeContains(columnType, TypeFloat32):
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
	case TypeContains(columnType, TypeBinaries):
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

		scale := ""
		if columnType["scale"] != "" {
			scale = columnType["scale"]
			//primary = ""

		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGoType(columnType["value"], nullable, gureguTypes)

		fieldName := FmtFieldName(strings.ToLower(StringifyFirstChar(key)))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s%s\"", key, primary, scale))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, ""))
		}
		if fieldName == "DeletedAt" || fieldName == "DELETE_DAT" {
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
