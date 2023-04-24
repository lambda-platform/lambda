package DBSchema

import (
	"fmt"
	"github.com/iancoleman/strcase"
	generatorModels "github.com/lambda-platform/lambda/generator/models"
	"strings"
)

const (
	gqlNullInt    = "Int32"
	gqlInt        = "Int32!"
	gqlNullBigInt = "Int32"
	gqlBigInt     = "Int32!"
	gqlNullFloat  = "Float"
	gqlFloat      = "Float!"
	gqlNullString = "String"
	gqlString     = "String!"
	gqlNullTime   = "Time"
	gqlTime       = "Time!"
	dbNullDate    = "Date!"
	dbDate        = "Date!"
	gqlBinary     = "Byte!"
)

func GenerateGrapql(columnTypes []generatorModels.ColumnData, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string, Subs []string, isInpute bool) ([]byte, error) {

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
func GenerateGrapqlOrder(columnTypes []generatorModels.ColumnData, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string) ([]byte, error) {

	dbTypes := generateQraphqlTypesOrder(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	src := fmt.Sprintf("\n  \ninput %s %s %s \n} %s",
		structName,
		dbTypes,
		extraColumns, extraStucts)

	return []byte(src), nil
}

func generateQraphqlTypes(columnTypes []generatorModels.ColumnData, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {

	structure := " {"

	for _, columnType := range columnTypes {

		nullable := false
		if columnType.Nullable == "YES" {
			nullable = true
		}

		primary := ""
		if columnType.Primary == "PRI" {
			primary = ";primaryKey;autoIncrement"
			//primary = ""
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		// If the guregu (https://github.com/guregu/null) CLI option is passed use its types, otherwise use go's sql.NullX

		valueType = sqlTypeToGraphyType(columnType.DataType, nullable, gureguTypes)

		if columnType.Primary == "PRI" {
			valueType = "ID!"
			//primary = ""
		}

		fieldName := columnType.Name
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", columnType.Name, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", columnType.Name, ""))
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
func generateQraphqlTypesOrder(columnTypes []generatorModels.ColumnData, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {

	structure := " {"

	for _, columnType := range columnTypes {

		primary := ""
		if columnType.Primary == "PRI" {
			primary = ";primaryKey;autoIncrement"
			//primary = ""
		}

		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", columnType.Name, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", columnType.Name, ""))
		}

		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n%s    %s: order_by",
				columnType.Name,
				strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n    %s: order_by",
				columnType.Name)
		}
	}

	return structure
}
func sqlTypeToGraphyType(columnType string, nullable bool, gureguTypes bool) string {
	switch {
	case TypeContains(columnType, TypeIntegers):
		if nullable {
			if gureguTypes {
				return gqlNullInt
			}
			return gqlNullInt
		}
		return gqlInt
	case TypeContains(columnType, TypeBool):
		if nullable {
			if gureguTypes {
				return gqlNullInt
			}
			return gqlNullInt
		}
		return gqlInt
	case TypeContains(columnType, TypeBigIntegers):
		if nullable {
			if gureguTypes {
				return gqlNullBigInt
			}
			return gqlBigInt
		}
		return gqlBigInt
	case TypeContains(columnType, TypeStrings):
		if nullable {
			if gureguTypes {
				return gqlNullString
			}
			return gqlNullString
		}
		return gqlString
	case TypeContains(columnType, TypeTimes):
		if nullable && gureguTypes {
			return gqlNullTime
		}
		return gqlTime
	case TypeContains(columnType, TypeDates):
		if nullable && gureguTypes {
			return dbNullDate
		}
		return dbDate
	case TypeContains(columnType, TypeFloat64):
		if nullable {
			if gureguTypes {
				return gqlNullFloat
			}
			return gqlNullFloat
		}
		return gqlFloat
	case TypeContains(columnType, TypeFloat32):
		if nullable {
			if gureguTypes {
				return gqlNullFloat
			}
			return gqlNullFloat
		}
		return gqlFloat
	case TypeContains(columnType, TypeBinaries):
		return gqlBinary
	}
	return ""
}
