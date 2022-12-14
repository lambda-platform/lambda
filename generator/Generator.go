package generator

import (
	"github.com/iancoleman/strcase"
	"github.com/lambda-platform/lambda/config"
	"strings"

	"github.com/lambda-platform/lambda/DBSchema"
	genertarModels "github.com/lambda-platform/lambda/generator/models"
	"github.com/lambda-platform/lambda/generator/utils"
	lambdaModels "github.com/lambda-platform/lambda/models"
	"github.com/otiai10/copy"
	"os"
	//"strconv"
	"fmt"
)

func ModelInit(dbSchema lambdaModels.DBSCHEMA, formSchemas []genertarModels.ProjectSchemas, gridSchemas []genertarModels.ProjectSchemas, copyClienModels bool, user_with_uuid string) {

	//dir := projectPath
	//dir := "schemas/" + projectPath

	AbsolutePath := utils.AbsolutePath()

	formPatch := "lambda/models/form/"
	gridPatch := "lambda/models/grid/"

	if _, err := os.Stat(formPatch); os.IsNotExist(err) {

		os.MkdirAll("lambda/models", 0755)

		os.MkdirAll("lambda/schemas", 0755)
		os.MkdirAll("lambda/microservices", 0755)
		os.MkdirAll("lambda/schemas/form", 0755)
		os.MkdirAll("lambda/schemas/menu", 0755)
		os.MkdirAll("lambda/schemas/grid", 0755)
		os.MkdirAll(formPatch, 0755)
		os.MkdirAll("lambda/models/form/formModels/", 0755)
		os.MkdirAll("lambda/models/form/caller/", 0755)

		os.MkdirAll(gridPatch, 0755)
		os.MkdirAll("lambda/models/grid/caller", 0755)

	} else {

		os.RemoveAll("lambda")
		os.MkdirAll("lambda/microservices", 0755)
		os.MkdirAll("lambda/models", 0755)
		os.MkdirAll("lambda/schemas", 0755)
		os.MkdirAll("lambda/schemas/form", 0755)
		os.MkdirAll("lambda/schemas/menu", 0755)
		os.MkdirAll("lambda/schemas/grid", 0755)
		os.MkdirAll(formPatch, 0755)
		os.MkdirAll("lambda/models/form/formModels/", 0755)
		os.MkdirAll("lambda/models/form/caller/", 0755)

		os.MkdirAll(gridPatch, 0755)
		os.MkdirAll("lambda/models/grid/caller", 0755)
	}

	WriteGridsModel(dbSchema, gridSchemas, copyClienModels)
	WriteFormsModelData(dbSchema, formSchemas, copyClienModels)

	if copyClienModels {

		if config.Config.Database.Connection == "oracle" {
			copy.Copy(AbsolutePath+"initialModels/dataform/modelsOracle/", "lambda/models/form/")
			copy.Copy(AbsolutePath+"initialModels/datagrid/modelsOracle/", "lambda/models/grid/")
		} else {
			copy.Copy(AbsolutePath+"initialModels/datagrid/models/", "lambda/models/grid/")
			if config.Config.SysAdmin.UUID || user_with_uuid == "true" {
				copy.Copy(AbsolutePath+"initialModels/dataform/modelsUUID/", "lambda/models/form/")
			} else {
				copy.Copy(AbsolutePath+"initialModels/dataform/models/", "lambda/models/form/")
			}
		}

	}
	fmt.Println("MODEL INIT DONE")
}

func GetColumnsFromTableMeta(columns []lambdaModels.TableMeta, hiddenColumns []string) *map[string]map[string]string {

	columnDataTypes := make(map[string]map[string]string)
	for _, tableColumn := range columns {
		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns {
			if hiddenColumn == tableColumn.Model {
				isHidden = true
			}
		}
		if isHidden == false {

			columnDataTypes[tableColumn.Model] = map[string]string{"value": tableColumn.DbType, "nullable": tableColumn.Nullable, "primary": tableColumn.Key, "scale": tableColumn.Scale}
		}
	}

	return &columnDataTypes

}
func GetModelAlias(modelName string) string {
	return DBSchema.FmtFieldName(DBSchema.StringifyFirstChar(modelName))
}

func TableMetaToStruct(columns []lambdaModels.TableMeta, table string, hiddenColumns []string, pkgName string, Subs []string) string {

	if table != "" {

		columnDataTypes := GetColumnsFromTableMeta(columns, hiddenColumns)

		subStchemas := ""

		for _, sub := range Subs {
			subStchemas = subStchemas + "\n    " + strcase.ToCamel(strings.ToLower(sub)) + " []*" + strcase.ToCamel(strings.ToLower(sub)) + "`gorm:\"-:all\"`"
		}

		struc_, _ := DBSchema.GenerateWithImports("", *columnDataTypes, table, GetModelAlias(strings.ToLower(table)), pkgName, true, true, true, subStchemas, "", "")

		return string(struc_)
	}

	return ""

}
func TableMetaToGraphql(columns []lambdaModels.TableMeta, table string, hiddenColumns []string, Subs []string, isInpute bool) string {

	if table != "" {

		columnDataTypes := GetColumnsFromTableMeta(columns, hiddenColumns)

		struc_, _ := DBSchema.GenerateGrapql(*columnDataTypes, table, GetModelAlias(strings.ToLower(table)), "", false, false, true, "", "", Subs, isInpute)

		return string(struc_)
	}

	return ""

}

func TableMetaColumnsWithMeta(columns []lambdaModels.TableMeta, table string, hiddenColumns []string) []map[string]string {

	if table != "" {

		tableColumns := GetColumnsWithMeta(columns, hiddenColumns)

		return tableColumns
	}

	return []map[string]string{}

}
func GetColumnsWithMeta(columns []lambdaModels.TableMeta, hiddenColumns []string) []map[string]string {

	tableColumns := []map[string]string{}
	for _, tableColumn := range columns {
		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns {
			if hiddenColumn == tableColumn.Model {
				isHidden = true
			}
		}
		if isHidden == false {

			newColumn := map[string]string{}

			newColumn["column"] = tableColumn.Model
			newColumn["nullable"] = tableColumn.Nullable
			newColumn["dataType"] = tableColumn.DbType

			tableColumns = append(tableColumns, newColumn)

		}
	}

	return tableColumns

}
