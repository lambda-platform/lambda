package handlers

import (
	"context"
	"encoding/json"
	"github.com/lambda-platform/lambda/DBSchema"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/datasource"
	"github.com/lambda-platform/lambda/generator"
	pb "github.com/lambda-platform/lambda/grpc/consoleProto"
	"google.golang.org/grpc"
	"os"

	"fmt"
	"github.com/lambda-platform/lambda/DB"
	agentModels "github.com/lambda-platform/lambda/agent/models"
	genertarModels "github.com/lambda-platform/lambda/generator/models"
	gUtils "github.com/lambda-platform/lambda/generator/utils"
	"github.com/lambda-platform/lambda/models"
	"strconv"
	"time"
)

func UploadDBSCHEMA() (*pb.Response, error) {

	conn, err := grpc.Dial(config.LambdaConfig.LambdaMainServicePath, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))

	if err != nil {
		return nil, err
	}

	defer conn.Close()
	c := pb.NewConsoleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lambdaConfig, err := os.ReadFile("lambda.json")
	if err != nil {
		return nil, err
	}
	r, err := c.UploadDBSCHEMA(ctx, &pb.SchemaParams{
		ProjectKey:   config.LambdaConfig.ProjectKey,
		DBSchema:     GetDBCHEMA(),
		LambdaConfig: lambdaConfig,
	})

	if err != nil {
		return nil, err
	}
	fmt.Println("DB SCHEMA SENT")
	return r, nil
}

func GetDBCHEMA() []byte {

	DBSchema.GenerateSchemaForCloud()

	b, err := os.ReadFile("app/models/db_schema.json")
	if err != nil {
		panic(err)
	}

	return b

}
func GetLambdaSCHEMA() {

	conn, err := grpc.Dial(config.LambdaConfig.LambdaMainServicePath, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))

	if err != nil {
		fmt.Println(err.Error())
	}

	defer conn.Close()
	c := pb.NewConsoleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := c.LambdaSCHEMA(ctx, &pb.LambdaSchemaParams{
		ProjectKey: config.LambdaConfig.ProjectKey,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	data := CloudData{}
	json.Unmarshal(r.Data, &data)

	//

	for _, ds := range data.DatasourceSchemas {
		datasource.DeleteView(ds.Name)
		errSave := datasource.CreateView(ds.Name, ds.Schema)
		if errSave != nil {
			fmt.Println(errSave.Error())
		}
	}
	fmt.Println("FORM & GRID generation starting")

	dbSchema := DBSchema.GetDBSchema()
	/*
	   Generate Form, Grid
	*/

	generator.ModelInit(dbSchema, data.FormSchemas, data.GridSchemas, true, config.Config.SysAdmin.UUID)

	/*
	   Generate GRAPHQL
	*/
	generator.GQLInit(dbSchema, data.GraphqlSchemas)

	for _, vb := range data.FormSchemas {
		_ = os.WriteFile("lambda/schemas/form/"+fmt.Sprintf("%d", vb.ID)+".json", []byte(vb.Schema), 0700)
	}
	for _, vb := range data.GridSchemas {
		_ = os.WriteFile("lambda/schemas/grid/"+fmt.Sprintf("%d", vb.ID)+".json", []byte(vb.Schema), 0700)
	}
	for _, vb := range data.MenuSchemas {
		_ = os.WriteFile("lambda/schemas/menu/"+fmt.Sprintf("%d", vb.ID)+".json", []byte(vb.Schema), 0700)
	}

	microservicesList := `
package microservices

import "github.com/lambda-platform/lambda/models"

`
	for _, projectData := range data.Projects {
		for _, projectSetting := range data.ProjectSettings {
			if projectData.ID == projectSetting.ProjectID {
				microservicesList = microservicesList + fmt.Sprintf(`
var %s models.Microservice = models.Microservice{
    GRPCURL: "%s",
    ProductionURL: "%s",
    ProjectID: %d,
}`, projectData.Name, projectSetting.GRPCURL, projectSetting.ProductionURL, projectData.ID)
			}

		}

	}

	Werror := gUtils.WriteFileFormat(microservicesList, "lambda/microservices/microservices.go")

	if Werror != nil {

		fmt.Println(Werror)
	}
}

func GetRoleData() error {

	conn, err := grpc.Dial(config.LambdaConfig.LambdaMainServicePath, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))

	if err != nil {
		fmt.Println(err.Error())
	}

	defer conn.Close()
	c := pb.NewConsoleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := c.RoleData(ctx, &pb.LambdaSchemaParams{
		ProjectKey: config.LambdaConfig.ProjectKey,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	data := map[string]interface{}{}

	json.Unmarshal(r.Data, &data)

	roleData := map[int]map[string]interface{}{}
	roleDataPre, _ := json.Marshal(data["roleData"])
	json.Unmarshal(roleDataPre, &roleData)

	Roles := []agentModels.Role{}
	RolesPre, _ := json.Marshal(data["roles"])
	json.Unmarshal(RolesPre, &Roles)

	for k, data := range roleData {

		bolB, _ := json.Marshal(data)
		_ = os.WriteFile("lambda/role_"+strconv.Itoa(k)+".json", bolB, 0700)
	}

	DB.DB.Exec("TRUNCATE roles")
	for _, Role := range Roles {
		DB.DB.Create(&Role)
	}

	return nil
}

type CloudData struct {
	GridSchemas       []genertarModels.ProjectSchemas `json:"grid-schemas"`
	FormSchemas       []genertarModels.ProjectSchemas `json:"form-schemas"`
	MenuSchemas       []genertarModels.ProjectSchemas `json:"menu-schemas"`
	ChartSchemas      []genertarModels.ProjectSchemas `json:"chart-schemas"`
	MoqupSchemas      []genertarModels.ProjectSchemas `json:"moqup-schemas"`
	DatasourceSchemas []genertarModels.ProjectSchemas `json:"datasource-schemas"`
	GraphqlSchemas    []genertarModels.ProjectSchemas `json:"graphql-schemas"`
	Projects          []models.Projects               `json:"projects"`
	ProjectSettings   []models.ProjectSettings        `json:"project-settings"`
}
