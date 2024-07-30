package DBDictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/DBSchema"
	"github.com/lambda-platform/lambda/config"
	genertarModels "github.com/lambda-platform/lambda/generator/models"
	pb "github.com/lambda-platform/lambda/grpc/consoleProto"
	puzzle "github.com/lambda-platform/lambda/puzzle/handlers"
	"google.golang.org/grpc"
	"time"
)

func Dictionary(c *fiber.Ctx) error {

	//tables := map[string][]string{}

	//tables["tables"] = []string{""}
	//tables["views"] = []string{""}

	//dbSchema := DBSchema.GetDBSchemaWithTargets(tables)
	dbSchema := DBSchema.GetDBSchema()

	var FormSchemasPre []genertarModels.ProjectSchemas
	var GridSchemasPre []genertarModels.ProjectSchemas
	var FormSchemas []interface{}
	var GridSchemas []interface{}

	if config.LambdaConfig.LambdaMainServicePath != "" {

		conn, err := grpc.Dial(config.LambdaConfig.LambdaMainServicePath, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(60*time.Second))

		if err != nil {
			fmt.Println(err.Error())
		}

		defer conn.Close()
		cc := pb.NewConsoleClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		r, err := cc.LambdaSCHEMA(ctx, &pb.LambdaSchemaParams{
			ProjectKey: config.LambdaConfig.ProjectKey,
		})

		if err != nil {
			fmt.Println(err.Error())
		}
		data := puzzle.CloudData{}

		json.Unmarshal(r.Data, &data)

		for _, vb := range data.FormSchemas {

			var result FormSCHEMA
			err := json.Unmarshal([]byte(vb.Schema), &result)
			if err == nil {
				result.Name = vb.Name
				FormSchemas = append(FormSchemas, result)
			}
		}
		for _, vb := range data.GridSchemas {
			var result GridSCHEMA
			err := json.Unmarshal([]byte(vb.Schema), &result)
			if err == nil {
				result.Name = vb.Name
				GridSchemas = append(GridSchemas, result)
			}
		}

	} else {

		if config.Config.Database.Connection == "oracle" {

			DB.DB.Select("ID AS \"id\", NAME AS \"name\", SCHEMA AS \"schema\", \"TYPE\" AS \"type\", CREATED_AT AS \"created_at\", UPDATED_AT AS \"updated_at\"").Where("TYPE = ?", "form").Table("VB_SCHEMAS").Find(&FormSchemasPre)
			DB.DB.Select("ID AS \"id\", NAME AS \"name\", SCHEMA AS \"schema\", \"TYPE\" AS \"type\", CREATED_AT AS \"created_at\", UPDATED_AT AS \"updated_at\"").Where("TYPE = ?", "grid").Table("VB_SCHEMAS").Find(&GridSchemasPre)

		} else {
			DB.DB.Where("type = ?", "form").Table("vb_schemas").Find(&FormSchemasPre)
			DB.DB.Where("type = ?", "grid").Table("vb_schemas").Find(&GridSchemasPre)

		}

		for _, vb := range FormSchemasPre {
			var result FormSCHEMA
			err := json.Unmarshal([]byte(vb.Schema), &result)
			if err == nil {
				result.Name = vb.Name
				FormSchemas = append(FormSchemas, result)
			}

		}
		for _, vb := range GridSchemasPre {
			var result GridSCHEMA
			err := json.Unmarshal([]byte(vb.Schema), &result)
			if err == nil {
				result.Name = vb.Name
				GridSchemas = append(GridSchemas, result)
			}
		}

	}
	return c.JSON(map[string]interface{}{
		"tableMetas":  dbSchema.TableMeta,
		"tableList":   dbSchema.TableList,
		"viewList":    dbSchema.ViewList,
		"formSchemas": FormSchemas,
		"gridSchemas": GridSchemas,
	})
}

func Set(e *fiber.App) {
	//e.Get("/db-dic", agentMW.IsLoggedIn(), agentMW.IsAdmin, Dictionary)
	e.Get("/db-dic", Dictionary)

}

type FormItem struct {
	Model    string `json:"model"`
	Label    string `json:"label"`
	Relation struct {
		Table  interface{}   `json:"table"`
		Key    interface{}   `json:"key"`
		Fields []interface{} `json:"fields"`
	} `json:"relation,omitempty"`
	Options []interface{} `json:"options"`
}

type FormSCHEMA struct {
	Name     string     `json:"name"`
	Model    string     `json:"model"`
	Identity string     `json:"identity"`
	Schema   []FormItem `json:"schema"`
}

type GridSCHEMA struct {
	Model    string `json:"model"`
	Name     string `json:"name"`
	Identity string `json:"identity"`
	Schema   []struct {
		Model string `json:"model"`
		Label string `json:"label"`
	} `json:"schema"`
}
