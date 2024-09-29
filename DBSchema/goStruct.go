package DBSchema

import (
	"fmt"
	"github.com/lambda-platform/lambda/config"
	generatorModels "github.com/lambda-platform/lambda/generator/models"
	"github.com/lambda-platform/lambda/utils"
	"go/format"
	"strings"
)

const (
	golangByteArray = "[]byte"
	//gureguNullInt    = "null.Int"
	golangNullInt   = "*int"
	golangInt       = "int"
	golangBool      = "bool"
	golangNullInt64 = "*int"
	golangInt64     = "int"
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
	dateNull       = "DB.Date"
	secureString   = "DB.SecureString"
)

func GenerateOnlyStruct(columnTypes []generatorModels.ColumnData, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string) ([]byte, error) {
	var dbTypes string

	dbTypes, _, _ = generateStructTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes, utils.StringInSlice(tableName, config.LambdaConfig.JsonLowerCaseTables))

	tableSchema := GetTableSchema(columnTypes)

	src := fmt.Sprintf("\n  \ntype %s %s %s} %s",
		structName,
		dbTypes,
		extraColumns, extraStucts)
	if gormAnnotation == true {
		tableNameFunc := "" +
			"func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") TableName() string {\n" +
			"	return \"" + tableSchema + tableName + "\"" +
			"}"
		src = fmt.Sprintf("%s\n%s", src, tableNameFunc)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {

		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}
func GenerateWithImports(otherPackage string, columnTypes []generatorModels.ColumnData, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string, virtualColums string) ([]byte, error) {
	var dbTypes string

	dbTypes, _, _ = generateStructTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes, utils.StringInSlice(tableName, config.LambdaConfig.JsonLowerCaseTables))

	tableSchema := GetTableSchema(columnTypes)

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
			"	return \"" + tableSchema + tableName + "\"" +
			"}"
		src = fmt.Sprintf("%s\n%s", src, tableNameFunc)
	}
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}
func GenerateWithImportsNoTime(otherPackage string, columnTypes []generatorModels.ColumnData, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string) ([]byte, error) {
	var dbTypes string

	dbTypes, timeFound := generateStructTypesNoTime(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes, utils.StringInSlice(tableName, config.LambdaConfig.JsonLowerCaseTables))
	tableSchema := GetTableSchema(columnTypes)
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
			"	return \"" + tableSchema + tableName + "\"" +
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
				return golangNullInt
			}
			return golangNullInt
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
				return golangNullInt64
			}
			return golangNullInt64
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

	case TypeContains(columnType, TypeGeo):
		return "string"

	case TypeContains(columnType, TypeSecure):
		return secureString
	}
	return ""
}

func generateStructTypes(columnTypes []generatorModels.ColumnData, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, jsonLowerCase bool) (string, bool, bool) {

	structure := "struct {"
	time_found := false
	date_found := false

	//sort.Strings(keys)

	for _, columnType := range columnTypes {

		nullable := false
		if columnType.Nullable == "YES" {
			nullable = true
		}
		if columnType.DataType == "timestamp" || columnType.DataType == "timestamptz" || columnType.DataType == "datetime" || columnType.DataType == "year" || columnType.DataType == "time" {

			//if key == "created_at" ||  key == "updated_at" ||  key == "deleted_at"{
			time_found = true
			//}
			//else {
			//	columnType.DataType = "text"
			//}

		}
		if columnType.DataType == "date" {
			date_found = true
		}
		primary := ""
		if columnType.Primary == "PRI" {
			primary = ";primaryKey;autoIncrement"
			//primary = ""
		}

		scale := ""
		if columnType.Scale != "" {
			scale = columnType.Scale
			//primary = ""

		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGoType(columnType.DataType, nullable, gureguTypes)

		fieldName := FmtFieldName(StringifyFirstChar(columnType.Name))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s%s\"", columnType.Name, primary, scale))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))

			if jsonLowerCase {
				annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", strings.ToLower(columnType.Name), ""))
			} else {
				annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", columnType.Name, ""))
			}

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
func generateStructTypesNoTime(columnTypes []generatorModels.ColumnData, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, jsonLowerCase bool) (string, bool) {

	structure := "struct {"
	time_found := false

	for _, columnType := range columnTypes {

		nullable := false
		if columnType.Nullable == "YES" {
			nullable = true
		}
		if columnType.DataType == "timestamp" || columnType.DataType == "timestamptz" || columnType.DataType == "datetime" || columnType.DataType == "date" || columnType.DataType == "year" || columnType.DataType == "time" {

			columnType.DataType = "text"

		}

		primary := ""
		if columnType.Primary == "PRI" {
			primary = ";primaryKey;autoIncrement"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGoType(columnType.DataType, nullable, gureguTypes)

		fieldName := FmtFieldName(StringifyFirstChar(columnType.Name))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", columnType.Name, primary))
		}
		if jsonAnnotation == true {
			if jsonLowerCase {
				annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", strings.ToLower(columnType.Name), ""))
			} else {
				annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", columnType.Name, ""))
			}

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

func GetTableSchema(columnTypes []generatorModels.ColumnData) string {
	schema := ""

	if len(columnTypes) >= 1 {
		if columnTypes[0].TableSchema != "" {
			schema = columnTypes[0].TableSchema + "."
		}
	}

	return schema
}
