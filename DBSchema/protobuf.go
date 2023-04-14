package DBSchema

import (
	"fmt"
	"github.com/iancoleman/strcase"
	generatorModels "github.com/lambda-platform/lambda/generator/models"
)

const (
	protoByteArray = "bytes"
	protoInt       = "int32"
	protoBigInt    = "int64"
	protoBool      = "bool"
	protoString    = "string"
	protoFloat32   = "float"
	protoFloat64   = "double"
	protoDate      = "string"
	protoTime      = "string"
)

func GenerateProtobuf(columnTypes []generatorModels.ColumnData, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool, extraColumns string, extraStucts string, Subs []string, isInpute bool) ([]byte, error) {

	dbTypes := generateProtobufTypes(columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	subStchemas := ""

	for _, sub := range Subs {
		subStchemas = subStchemas + "\n    " + sub + ":[" + strcase.ToCamel(sub) + "!]"
	}

	typeSchema := "message"
	if isInpute {
		typeSchema = "input"
		structName = structName + "Input"
	}

	src := fmt.Sprintf("syntax = \"proto3\";\n\npackage %s;\noption go_package = \"./;%s\";\n%s %s %s %s %s \n} %s",
		tableName,
		tableName,
		typeSchema,
		structName,
		dbTypes,
		extraColumns, subStchemas, extraStucts)

	return []byte(src), nil
}

func generateProtobufTypes(columnTypes []generatorModels.ColumnData, depth int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {

	structure := " {"

	for index, columnType := range columnTypes {
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

		valueType = sqlTypeToProtobufType(columnType.DataType, nullable, gureguTypes)

		fieldName := FmtFieldName(StringifyFirstChar(columnType.Name))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", columnType.Name, primary))
		}
		if jsonAnnotation == true {
			//annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", key, primary))
			annotations = append(annotations, fmt.Sprintf("json:\"%s%s\"", columnType.Name, ""))
		}

		if len(annotations) > 0 {
			structure += fmt.Sprintf("\n\t%s %s = %d;",

				valueType,
				fieldName,
				index+1)

			//strings.Join(annotations, " "))

		} else {
			structure += fmt.Sprintf("\n\t%s %s = %d;",

				valueType,
				fieldName,
				index+1)
		}
	}

	return structure
}

func sqlTypeToProtobufType(columnType string, nullable bool, gureguTypes bool) string {
	switch {
	case TypeContains(columnType, TypeIntegers):
		if nullable {
			if gureguTypes {
				return protoInt
			}
			return protoInt
		}
		return protoInt
	case TypeContains(columnType, TypeBool):
		if nullable {
			if gureguTypes {
				return protoBool
			}
			return protoBool
		}
		return protoBool
	case TypeContains(columnType, TypeBigIntegers):
		if nullable {
			if gureguTypes {
				return protoBigInt
			}
			return protoBigInt
		}
		return protoBigInt
	case TypeContains(columnType, TypeStrings):
		if nullable {
			if gureguTypes {
				return protoString
			}
			return protoString
		}
		return protoString
	case TypeContains(columnType, TypeTimes):
		if nullable && gureguTypes {
			return protoTime
		}
		return protoTime
	case TypeContains(columnType, TypeDates):
		if nullable && gureguTypes {
			return protoDate
		}
		return protoDate
	case TypeContains(columnType, TypeFloat64):
		if nullable {
			if gureguTypes {
				return protoFloat64
			}
			return protoFloat64
		}
		return protoFloat64
	case TypeContains(columnType, TypeFloat32):
		if nullable {
			if gureguTypes {
				return protoFloat32
			}
			return protoFloat32
		}
		return protoFloat32
	case TypeContains(columnType, TypeBinaries):
		return protoByteArray
	}
	return ""
}
