package generator

import (
	"encoding/json"
	"fmt"
	"github.com/lambda-platform/lambda/DBSchema"
	genertarModels "github.com/lambda-platform/lambda/generator/models"
	"github.com/lambda-platform/lambda/generator/utils"
	lambdaModels "github.com/lambda-platform/lambda/models"
	lambdaUtils "github.com/lambda-platform/lambda/utils"
	"strconv"
	"strings"
)

func WriteGridsModel(dbSchema lambdaModels.DBSCHEMA, grids []genertarModels.ProjectSchemas, copyClienModels bool) {

	genertedGrids := WriteGridModel(dbSchema, grids)
	WriteGridDataCaller(genertedGrids, copyClienModels)
}
func WriteGridModel(dbSchema lambdaModels.DBSCHEMA, grids []genertarModels.ProjectSchemas) []genertarModels.ProjectSchemas {
	genertedGrids := []genertarModels.ProjectSchemas{}
	for _, vb := range grids {
		var schema lambdaModels.SCHEMAGRID

		json.Unmarshal([]byte(vb.Schema), &schema)

		modelAlias := GetModelAlias(schema.Model)
		MainTableAlias := GetModelAlias(schema.MainTable) + "MainTable"
		if schema.MainTable == "" {
			MainTableAlias = GetModelAlias(modelAlias) + "MainTable"

			schema.MainTable = schema.Model
		}

		modelAliasWithID := modelAlias + strconv.FormatInt(int64(vb.ID), 10)
		MainTableAliasWithID := MainTableAlias + strconv.FormatInt(int64(vb.ID), 10)

		/*Grid Columns */
		columns, columnList, virtualColums := createColumns(schema, modelAliasWithID)

		/*GRID Relation & Microservices*/
		relations, MicroserviceCaller, microRelationFound := createRelation(schema, modelAliasWithID)

		/*Create Grid Model*/
		models := createModel(schema, dbSchema, modelAliasWithID, MainTableAliasWithID, virtualColums, microRelationFound)

		/*GRID DEFAULT CONDITION & Filters*/
		filters := createFilter(schema, modelAliasWithID)

		/*GRID Aggergation*/
		aggergations := createAggergation(schema, modelAliasWithID)

		/*Grid Trigger*/
		triggers := createTrigger(schema, modelAliasWithID, modelAlias, vb.ID)

		content := fmt.Sprintf(`%s

var %sDatagrid datagrid.Datagrid = datagrid.Datagrid{
    Name: "%s",
    Identity: "%s",
    DataTable: "%s",
    MainTable: "%s",
    DataModel:  new(%s),
    Data:  new([]%s),
    MainModel:  new(%s),
    Columns:%s,
	ColumnList:%s,
    Filters: %s,
    Relations: %s,
    Condition: "%s",
    Aggergation: "%s",
   	%s
    TriggerNameSpace: "%s",
	FillVirtualColumns: fillVirtualColumns%s,
}

func fillVirtualColumns%s(rowsPre interface{}) interface{}{
    %s
}
`, models, modelAliasWithID, vb.Name, schema.Identity, schema.Model, schema.MainTable, modelAliasWithID, modelAliasWithID, MainTableAliasWithID, columns, columnList, filters, relations, schema.Condition, aggergations, triggers, schema.Triggers.Namespace, modelAliasWithID, modelAliasWithID, MicroserviceCaller)

		Werror := utils.WriteFileFormat(content, "lambda/models/grid/"+modelAlias+strconv.FormatInt(int64(vb.ID), 10)+".go")
		if Werror == nil {
			genertedGrids = append(genertedGrids, vb)
		} else {
			fmt.Println(Werror)
		}

	}

	return genertedGrids

}
func WriteGridDataCaller(grids []genertarModels.ProjectSchemas, copyClienModels bool) {
	content := "package caller\n"

	content = content + "import \"lambda/lambda/models/grid\"\n"
	content = content + "import \"github.com/lambda-platform/lambda/datagrid\"\n"

	content = content + "func GetMODEL(schema_id string) (datagrid.Datagrid) {\n\nswitch schema_id {\n"

	if copyClienModels {
		content = content + ` 

		case "crud_grid":
			return grid.KrudGridDatagrid

		case "crud_log":
			return grid.CrudLogDatagrid

		case "analytic_grid":
			return grid.AnalyticGridDatagrid

 		case "menu_grid":
			return grid.MenuGridDatagrid

 		case "notification_target_grid":
			return grid.NotificationTargetDatagrid
 		
`
	}

	for _, vb := range grids {
		var schema lambdaModels.SCHEMAGRID

		json.Unmarshal([]byte(vb.Schema), &schema)

		modelAlias := GetModelAlias(schema.Model)

		content = content + "\n case \"" + strconv.FormatInt(int64(vb.ID), 10) + "\": \nreturn grid." + modelAlias + strconv.FormatInt(int64(vb.ID), 10) + "Datagrid\n"

	}

	content = content + "\n} \nreturn datagrid.Datagrid{}\n\n}"

	utils.WriteFileFormat(content, "lambda/models/grid/caller/modelCaller.go")
}
func createModel(schema lambdaModels.SCHEMAGRID, dbSchema lambdaModels.DBSCHEMA, modelAliasWithID string, MainTableAliasWithID string, virtualColums string, microRelationFound bool) string {
	hiddenColumns := []string{}

	for _, column := range schema.Schema {
		if column.Hide == true && column.Model != schema.Identity && column.Model != "deleted_at" && column.Model != "created_at" && column.Model != "updated_at" && column.Model != "DELETED_AT" && column.Model != "CREATED_AT" && column.Model != "UPDATED_AT" {
			if column.Label == "" {
				hiddenColumns = append(hiddenColumns, column.Model)
			}
		}
	}

	columnDataTypes := GetColumnsFromTableMeta(dbSchema.TableMeta[schema.Model], hiddenColumns)

	importPackages := ""

	if schema.Triggers.Namespace != "" {

		importPackages = "\n import \"" + schema.Triggers.Namespace + "\" \n"

	}
	if microRelationFound {
		importPackages = importPackages + "\n import \"github.com/lambda-platform/lambda/utils\" \n"

		importPackages = importPackages + "\n import \"github.com/lambda-platform/lambda/grpc\" \n"
	}
	importPackages = importPackages + "\n import \"github.com/lambda-platform/lambda/datagrid\" \n"
	importPackages = importPackages + "\n import \"github.com/lambda-platform/lambda/models\" \n"

	MainTableColumnDataTypes := GetColumnsFromTableMeta(dbSchema.TableMeta[schema.MainTable], []string{})

	MainTableStructs, _ := DBSchema.GenerateOnlyStruct(*MainTableColumnDataTypes, schema.MainTable, MainTableAliasWithID, "", true, true, true, "", "")

	struc, _ := DBSchema.GenerateWithImports(importPackages, *columnDataTypes, schema.Model, modelAliasWithID, "grid", true, true, true, "", string(MainTableStructs), virtualColums)

	return string(struc)

}
func createFilter(schema lambdaModels.SCHEMAGRID, modelAliasWithID string) string {

	gridFilter := `map[string]string{`
	for i := range schema.Schema {
		//if schema.Schema[i].Filterable == true {

		gridFilter = gridFilter + `
					"` + schema.Schema[i].Model + `":"` + schema.Schema[i].Filter.Type + `",
`
		//}
	}

	return gridFilter + `
			}`

}
func createColumns(schema lambdaModels.SCHEMAGRID, modelAliasWithID string) (string, string, string) {
	list := `[]string{"` + schema.Identity + `"`
	gridColumns := `[]datagrid.Column{
`
	virtualColumns := ""

	for i := range schema.Schema {
		if schema.Schema[i].Hide == false && schema.Schema[i].Model != schema.Identity {
			gridColumns = gridColumns + `datagrid.Column{Model: "` + schema.Schema[i].Model + `",Label: "` + schema.Schema[i].Label + `"},
`

			if schema.Schema[i].VirtualColumn {
				virtualColumns = virtualColumns + "\n" + GetModelAlias(schema.Schema[i].Model) + " " + schema.Schema[i].DbType + "  `gorm:\"column:" + schema.Schema[i].Model + "\" json:\"" + schema.Schema[i].Model + "\"`\n"

				if schema.Schema[i].Relation.Table != "" && schema.Schema[i].Relation.Key != "" && schema.Schema[i].Relation.Fields != "" && schema.Schema[i].Relation.Self {

					list = list + `, "` + fmt.Sprintf(`(SELECT %s FROM %s WHERE %s IN (%s) limit 1) as `, schema.Schema[i].Relation.Fields, schema.Schema[i].Relation.Table, schema.Schema[i].Relation.Key, schema.Schema[i].Relation.ConnectionField) + schema.Schema[i].Model + `"`

				}
				if schema.Schema[i].Relation.Table != "" && schema.Schema[i].Relation.Key != "" && schema.Schema[i].Relation.Fields != "" && !schema.Schema[i].Relation.Self {

				}

			} else {
				list = list + `, "` + schema.Schema[i].Model + `"`
			}
		} else {
			if schema.Schema[i].Label != "" && schema.Schema[i].Model != schema.Identity {
				list = list + `, "` + schema.Schema[i].Model + `"`
			}
		}

	}
	gridColumns = gridColumns + `
}`
	return gridColumns, list + `}`, virtualColumns

}
func createTrigger(schema lambdaModels.SCHEMAGRID, modelAliasWithID string, modelAlias string, vbID int) string {

	beforeFetchMethod := `nil`

	if schema.Triggers.BeforeFetch != "" {
		beforeFetchMethod = schema.Triggers.BeforeFetch

	}

	afterFetchMethod := `nil`

	if schema.Triggers.AfterFetch != "" {
		afterFetchMethod = schema.Triggers.AfterFetch

	}

	beforeDeleteMethod := `nil`

	if schema.Triggers.BeforeDelete != "" {
		beforeDeleteMethod = schema.Triggers.BeforeDelete

	}

	afterDeleteMethod := `nil`

	if schema.Triggers.AfterDelete != "" {
		afterDeleteMethod = strings.ReplaceAll(schema.Triggers.AfterDelete, "@", ".")

	}

	beforePrintMethod := `nil`

	if schema.Triggers.BeforePrint != "" {
		beforePrintMethod = schema.Triggers.BeforePrint

	}

	return `BeforeFetch:` + beforeFetchMethod + `,
			
				AfterFetch:` + afterFetchMethod + `,
				
				BeforeDelete:` + beforeDeleteMethod + `,
			
				AfterDelete:` + afterDeleteMethod + `,
				
				BeforePrint:` + beforePrintMethod + `,`
}
func createAggergation(schema lambdaModels.SCHEMAGRID, modelAliasWithID string) string {
	gridAggergation := ``
	for i, aggergation := range schema.ColumnAggregations {

		if i <= 0 {
			gridAggergation = gridAggergation + `` + aggergation["aggregation"] + `(` + aggergation["column"] + `) as ` + aggergation["aggregation"] + `_` + aggergation["column"]
		} else {
			gridAggergation = gridAggergation + `, ` + aggergation["aggregation"] + `(` + aggergation["column"] + `) as ` + aggergation["aggregation"] + `_` + aggergation["column"]
		}

	}
	return gridAggergation
}
func createRelation(schema lambdaModels.SCHEMAGRID, modelAliasWithID string) (string, string, bool) {
	IDvariables := ""
	AppendIDvariables := ""
	microserviceClients := ""
	dataFillers := ""

	microserviceFound := false
	gridRelation := `[]models.GridRelation{`
	for _, column := range schema.Schema {

		if column.Relation.Table != "" && column.Relation.Key != "" && column.Relation.Fields != "" && !column.Relation.Self && column.Relation.MicroserviceID >= 1 {

			indexOfMicro := lambdaUtils.IndexOfMicro(column.Relation.MicroserviceID, schema.Microservices)
			filedAlias := GetModelAlias(column.Relation.Fields)
			connectionAlies := DBSchema.FmtFieldName(column.Relation.ConnectionField)
			if indexOfMicro >= 1 {

				IDvariables = IDvariables + fmt.Sprintf(`
%sIDs := []string{}
`, column.Relation.Fields)
				AppendIDvariables = AppendIDvariables + fmt.Sprintf(`
%sIDs = append(%sIDs, utils.GetString(row.%s))
`, column.Relation.Fields, column.Relation.Fields, connectionAlies)

				microserviceClients = microserviceClients + fmt.Sprintf(`
%sRows, %sError := grpc.CallStringData("%s", "%s", "%s", "%s", %sIDs)
`, column.Relation.Fields, column.Relation.Fields, schema.Microservices[indexOfMicro].GRPCURL, column.Relation.Table, column.Relation.Key, column.Relation.Fields, column.Relation.Fields)

				dataFillers = dataFillers + fmt.Sprintf(`
 			if %sError == nil{
                %sIndex := utils.IndexOf(rows[i].%s , %sRows)
                if %sIndex >= 0 {
                    rows[i].%s = %sRows.Rows[%sIndex].Value
                }
            } else {
                rows[i].%s = "Холболт амжилтгүй"
            }
`, column.Relation.Fields, column.Relation.Fields, connectionAlies, column.Relation.Fields, column.Relation.Fields, filedAlias, column.Relation.Fields, column.Relation.Fields, filedAlias)

			}

			microserviceFound = true
			relation := fmt.Sprintf(`models.GridRelation{
				MicroserviceID: %d,
				Table: "%s",
				Key: "%s",
				Fields: "%s",
				Column: "%s",
				Self: %t,
				ConnectionField: "%s",
			},`, column.Relation.MicroserviceID, column.Relation.Table, column.Relation.Key, column.Relation.Fields, column.Model, column.Relation.Self, column.Relation.ConnectionField)

			gridRelation = gridRelation + relation + "\n"

		}
	}

	gridRelation = gridRelation + `}`

	microserviceCaller := ``

	if microserviceFound {

		microserviceCaller = fmt.Sprintf(`
	
    rowsData, ok := rowsPre.(*[]%s)
 	if ok {
rows := []%s{}
    rows = *rowsData

    %s
    for _, row := range rows {
        %s       
    }

    %s
	

   
        for i, _ := range rows {

           %s


        }
        return rows

    } else {

        return rowsPre
    }
`, modelAliasWithID, modelAliasWithID, IDvariables, AppendIDvariables, microserviceClients, dataFillers)
	} else {
		microserviceCaller = "return rowsPre"
	}

	return gridRelation, microserviceCaller, microserviceFound
}
