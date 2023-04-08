package generator

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/lambda-platform/lambda/DBSchema"
)

func GetStruct(table string) {

	if table != "" {

		tableMeta := DBSchema.TableMetas(table)

		columnDataTypes := GetColumnsFromTableMeta(tableMeta, []string{})

		tableStruct, _ := DBSchema.GenerateOnlyStruct(columnDataTypes, table, strcase.ToCamel(table), "models", true, true, true, "", "")
		fmt.Println(string(tableStruct))
	}

}
