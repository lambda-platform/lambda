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

func ModelInit(dbSchema lambdaModels.DBSCHEMA, formSchemas []genertarModels.ProjectSchemas, gridSchemas []genertarModels.ProjectSchemas, copyClientModels bool, WithUUID bool) {
	formPatch := "lambda/models/form/"
	gridPatch := "lambda/models/grid/"
	directories := []string{
		"lambda/models",
		"lambda/schemas",
		"lambda/microservices",
		"lambda/schemas/form",
		"lambda/schemas/menu",
		"lambda/schemas/grid",
		formPatch,
		"lambda/models/form/formModels/",
		"lambda/models/form/caller/",
		gridPatch,
		"lambda/models/grid/caller",
	}

	desiredPermissions := os.FileMode(0700)
	//desiredUser := os.Getuid()
	//desiredGroup := os.Getgid()

	for _, dir := range directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, desiredPermissions)
			//err = os.Chown(dir, desiredUser, desiredGroup)
		} else {
			err = os.Chmod(dir, desiredPermissions)
			//err = os.Chown(dir, desiredUser, desiredGroup)
		}
	}

	WriteGridsModel(dbSchema, gridSchemas, copyClientModels)
	WriteFormsModelData(dbSchema, formSchemas, copyClientModels)
	AbsolutePath := utils.AbsolutePath()
	if copyClientModels {
		if config.Config.Database.Connection == "oracle" {
			copy.Copy(AbsolutePath+"initialModels/dataform/modelsOracle/", "lambda/models/form/")
			copy.Copy(AbsolutePath+"initialModels/datagrid/modelsOracle/", "lambda/models/grid/")
		} else {
			copy.Copy(AbsolutePath+"initialModels/datagrid/models/", "lambda/models/grid/")
			if config.Config.SysAdmin.UUID || WithUUID {
				copy.Copy(AbsolutePath+"initialModels/dataform/modelsUUID/", "lambda/models/form/")
			} else {
				copy.Copy(AbsolutePath+"initialModels/dataform/models/", "lambda/models/form/")
			}
		}
	}
	fmt.Println("MODEL INIT DONE")
}

func GetColumnsFromTableMeta(columns []lambdaModels.TableMeta, hiddenColumns []string) []genertarModels.ColumnData {

	var columnTypes []genertarModels.ColumnData
	for _, tableColumn := range columns {
		var isHidden bool = false

		for _, hiddenColumn := range hiddenColumns {
			if hiddenColumn == tableColumn.Model {
				isHidden = true
			}
		}
		if isHidden == false {

			columnTypes = append(columnTypes, genertarModels.ColumnData{
				Name:     tableColumn.Model,
				DataType: tableColumn.DbType,
				Nullable: tableColumn.Nullable,
				Primary:  tableColumn.Key,
				Scale:    tableColumn.Scale,
			})
		}
	}

	return columnTypes

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

		struc_, _ := DBSchema.GenerateWithImports("", columnDataTypes, table, GetModelAlias(strings.ToLower(table)), pkgName, true, true, true, subStchemas, "", "")

		return string(struc_)
	}

	return ""

}
func TableMetaToGraphql(columns []lambdaModels.TableMeta, table string, hiddenColumns []string, Subs []string, isInpute bool) string {

	if table != "" {

		columnDataTypes := GetColumnsFromTableMeta(columns, hiddenColumns)

		struc_, _ := DBSchema.GenerateGrapql(columnDataTypes, table, GetModelAlias(strings.ToLower(table)), "", false, false, true, "", "", Subs, isInpute)

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
