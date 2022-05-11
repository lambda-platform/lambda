package utils

import (
	//"encoding/json"
	//"fmt"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/crudlogger/models"
	//puzzleModels "github.com/lambda-platform/lambda/DB/DBSchema/models"
	//analyticModels "github.com/lambda-platform/lambda/lambda/plugins/dataanalytic/models"
	//"os"
)

func AutoMigrateSeed() {

	if config.Config.Database.Connection == "mssql" {
		DB.DB.AutoMigrate(

			&models.CrudLogMSSQL{},
		)
	} else {
		DB.DB.AutoMigrate(

			&models.CrudLog{},
		)
	}

	//if config.Config.App.Seed == "true" {
	//
	//	var vbs []puzzleModels.VBSchema
	//	DB.DB.Where("name = ?", "Cистем лог").Find(&vbs)
	//	if len(vbs) <= 0 {
	//		seedData()
	//	}
	//}
}
func seedData() {

	//var vbs []puzzleModels.VBSchemaAdmin
	//AbsolutePath := AbsolutePath()
	//dataFile, err := os.Open(AbsolutePath+"initialData/vb_schemas_admin.json")
	//defer dataFile.Close()
	//if err != nil {
	//	fmt.Println("PUZZLE SEED ERROR")
	//}
	//jsonParser := json.NewDecoder(dataFile)
	//err = jsonParser.Decode(&vbs)
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println("PUZZLE SEED DATA ERROR")
	//}
	////fmt.Println(len(vbs))
	//
	//for _, vb := range vbs {
	//
	//	DB.DB.Create(&vb)
	//}
	//
	//
	//var vbs2 []puzzleModels.VBSchema
	//
	//dataFile2, err2 := os.Open(AbsolutePath+"initialData/vb_schemas.json")
	//defer dataFile2.Close()
	//if err2 != nil {
	//	fmt.Println("PUZZLE SEED ERROR")
	//}
	//jsonParser2 := json.NewDecoder(dataFile2)
	//err = jsonParser2.Decode(&vbs2)
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println("PUZZLE SEED DATA ERROR")
	//}
	////fmt.Println(len(vbs))
	//
	//for _, vb := range vbs2 {
	//
	//	DB.DB.Create(&vb)
	//
	//}
	//
	//query := `CREATE VIEW ds_crud_log as SELECT
	//	crud_log.id as id,
	//		crud_log.user_id as user_id,
	//		crud_log.ip as ip,
	//		crud_log.user_agent as user_agent,
	//		crud_log.action as action,
	//		crud_log.schemaId as schemaId,
	//		crud_log.row_id as row_id,
	//		crud_log.input as input,
	//		crud_log.created_at as created_at,
	//		vb_schemas.name as name,
	//		users.role as role,
	//		users.login as login,
	//		users.email as email,
	//		users.first_name as first_name,
	//		users.last_name as last_name
	//	FROM
	//	crud_log
	//	LEFT JOIN vb_schemas on crud_log.schemaId = vb_schemas.id
	//	LEFT JOIN users on crud_log.user_id = users.id`
	//
	//DB.DB.Exec(query)
}
