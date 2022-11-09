package utils

import (
	"encoding/json"
	"fmt"
	//"encoding/json"
	//"fmt"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/crudlogger/models"
	puzzleModels "github.com/lambda-platform/lambda/models"
	"os"
	//puzzleModels "github.com/lambda-platform/lambda/DB/DBSchema/models"
	//analyticModels "github.com/lambda-platform/lambda/lambda/plugins/dataanalytic/models"
	//"os"
)

func AutoMigrateSeed() {

	if config.Config.Database.Connection == "mssql" {
		DB.DB.AutoMigrate(

			&models.CrudLogMSSQL{},
		)
	} else if config.Config.Database.Connection == "oracle" {
		DB.DB.AutoMigrate(

			&models.CrudLogOracle{},
		)
	} else {
		DB.DB.AutoMigrate(

			&models.CrudLog{},
		)
	}

	if config.Config.App.Seed == "true" {
		if config.Config.Database.Connection == "oracle" {
			var vbs []puzzleModels.VBSchemaOracle
			DB.DB.Where("NAME = ?", "Cистем лог").Find(&vbs)
			if len(vbs) <= 0 {
				seedData()
			}
		} else {
			var vbs []puzzleModels.VBSchema
			DB.DB.Where("name = ?", "Cистем лог").Find(&vbs)
			if len(vbs) <= 0 {
				seedData()
			}
		}
	}
}
func seedData() {

	AbsolutePath := AbsolutePath()

	var vbs2 []puzzleModels.VBSchema

	file := "initialData/vb_schemas.json"
	if config.Config.Database.Connection == "oracle" {
		file = "initialData/vb_schemasOracle.json"
	}
	dataFile2, err2 := os.Open(AbsolutePath + file)
	defer dataFile2.Close()
	if err2 != nil {
		fmt.Println("PUZZLE SEED ERROR")
	}
	jsonParser2 := json.NewDecoder(dataFile2)
	err := jsonParser2.Decode(&vbs2)
	if err != nil {
		fmt.Println(err)
		fmt.Println("PUZZLE SEED DATA ERROR")
	}
	//fmt.Println(len(vbs))

	for _, vb := range vbs2 {

		DB.DB.Create(&vb)

	}
	if config.Config.Database.Connection == "oracle" {
		query := `CREATE VIEW DS_CRUD_LOG AS SELECT
		CRUD_LOG.ID AS ID,
			CRUD_LOG.USER_ID AS USER_ID,
			CRUD_LOG.IP AS IP,
			CRUD_LOG.USER_AGENT AS USER_AGENT,
			CRUD_LOG.ACTION AS ACTION,
			CRUD_LOG.SCHEMA_ID AS SCHEMA_ID,
			CRUD_LOG.ROW_ID AS ROW_ID,
			CRUD_LOG.INPUT AS INPUT,
			CRUD_LOG.CREATED_AT AS CREATED_AT,
			VB_SCHEMAS.NAME AS NAME,
			USERS.ROLE AS ROLE,
			USERS.LOGIN AS LOGIN,
			USERS.EMAIL AS EMAIL,
			USERS.FIRST_NAME AS FIRST_NAME,
			USERS.LAST_NAME AS LAST_NAME
		FROM
		CRUD_LOG
		LEFT JOIN VB_SCHEMAS ON CRUD_LOG.SCHEMA_ID = VB_SCHEMAS.ID
		LEFT JOIN USERS ON CRUD_LOG.USER_ID = USERS.ID`

		DB.DB.Exec(query)
	} else {
		query := `CREATE VIEW ds_crud_log as SELECT
		crud_log.id as id,
			crud_log.user_id as user_id,
			crud_log.ip as ip,
			crud_log.user_agent as user_agent,
			crud_log.action as action,
			crud_log.schema_id as schema_id,
			crud_log.row_id as row_id,
			crud_log.input as input,
			crud_log.created_at as created_at,
			vb_schemas.name as name,
			users.role as role,
			users.login as login,
			users.email as email,
			users.first_name as first_name,
			users.last_name as last_name
		FROM
		crud_log
		LEFT JOIN vb_schemas on crud_log.schema_id = vb_schemas.id
		LEFT JOIN users on crud_log.user_id = users.id`

		DB.DB.Exec(query)
	}
}
