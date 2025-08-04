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
	} else if config.Config.Database.Connection == "postgres" {
		DB.DB.AutoMigrate(

			&models.CrudLogUUID{},
			&models.CrudTableType{},
		)
	} else {
		if config.Config.SysAdmin.UUID {
			DB.DB.AutoMigrate(

				&models.CrudLogUUID{},
			)

		} else {
			DB.DB.AutoMigrate(

				&models.CrudLog{},
			)
		}

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
		query := `CREATE VIEW view_crud_log AS SELECT
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
	} else if config.Config.Database.Connection == "postgres" {
		query := `CREATE OR REPLACE VIEW view_crud_log AS
SELECT 
    crud_log.id,
    crud_log.user_id,
    crud_log.ip,
    crud_log.user_agent,
    crud_log.action,
    crud_log.schema_id,
    crud_log.row_id,
    crud_log.input,
    crud_log.created_at,
    vb_schemas.name,
    ctt.db_table_name,
    users.role,
    users.login,
    users.email,
    users.first_name,
    users.last_name,
    ctt.log_type AS geo_type
FROM crud_log
LEFT JOIN vb_schemas ON crud_log.schema_id = vb_schemas.id
LEFT JOIN users ON crud_log.user_id = users.id
LEFT JOIN crud_table_type ctt ON crud_log.schema_id = ctt.schema_id;

-- On crud_log (main large table)
CREATE INDEX idx_crud_log_created_at ON crud_log (created_at DESC);  -- For ORDER BY created_at DESC and date filters
CREATE INDEX idx_crud_log_action ON crud_log (action);  -- For action filter
CREATE INDEX idx_crud_log_schema_id ON crud_log (schema_id);  -- For join to vb_schemas
CREATE INDEX idx_crud_log_user_id ON crud_log (user_id);  -- For join to users

-- Composite index for common filters + sort (if single indexes aren't enough)
CREATE INDEX idx_crud_log_schema_action_created ON crud_log (schema_id, action, created_at DESC);

-- On vb_schemas
CREATE INDEX idx_vb_schemas_type ON vb_schemas (type);  -- For WHERE type IN ('form', 'grid')
CREATE INDEX idx_vb_schemas_name ON vb_schemas (name);  -- For category/name filter

-- On crud_table_type
CREATE INDEX idx_crud_table_type_log_type ON crud_table_type (log_type);  -- For geo_type/log_type filter

-- On users (if not already indexed; assume id is PK, so it is)
CREATE INDEX idx_users_id ON users (id);  -- Redundant if id is PRIMARY KEY

CREATE OR REPLACE FUNCTION "public"."insert_crud_table_type"()
  RETURNS "pg_catalog"."trigger" AS $BODY$
BEGIN
    IF NEW.type IN ('form', 'grid') THEN
        DECLARE
            tbl_name TEXT := (NEW.schema::jsonb ->> 'model');
            is_geo BOOLEAN;
        BEGIN
            SELECT EXISTS (
                SELECT 1 FROM geometry_columns
                WHERE f_table_schema = 'public' AND f_table_name = tbl_name
            ) INTO is_geo;
            
            INSERT INTO crud_table_type (schema_id, db_table_name, log_type)
            VALUES (NEW.id, tbl_name, CASE WHEN is_geo THEN 'geometry' ELSE 'normal' END)
            ON CONFLICT (schema_id) DO UPDATE SET 
                db_table_name = EXCLUDED.db_table_name, 
                log_type = EXCLUDED.log_type;
        END;
    END IF;
    RETURN NEW;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;

INSERT INTO crud_table_type (schema_id, db_table_name, log_type)
SELECT 
    id, 
    schema::jsonb ->> 'model' AS db_table_name,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM geometry_columns 
            WHERE f_table_schema = 'public' AND f_table_name = (schema::jsonb ->> 'model')
        ) THEN 'geometry' 
        ELSE 'normal' 
    END AS log_type
FROM vb_schemas
WHERE type IN ('form', 'grid')
ON CONFLICT (schema_id) DO NOTHING;
`

		DB.DB.Exec(query)
	} else {
		query := `CREATE VIEW view_crud_log as SELECT
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
