package grpc

import (
	"context"
	"fmt"
	"github.com/lambda-platform/lambda/DB"
	pb "github.com/lambda-platform/lambda/grpc/proto"
	"github.com/lambda-platform/lambda/puzzle/handlers"
	"os"
)

func GetIntData(ctx context.Context, in *pb.TableOption) (*pb.IntRows, error) {

	var rows *pb.IntRows = &pb.IntRows{}

	DB.DB.Table(in.Table).Select(in.Key+" as key, "+in.Field+" as value").Where(in.Key+" IN (?)", in.Values).Find(&rows.Rows)

	return rows, nil
}
func GetStringData(ctx context.Context, in *pb.TableOption) (*pb.StringRows, error) {

	var rows *pb.StringRows = &pb.StringRows{}

	DB.DB.Table(in.Table).Select(in.Key+" as key, "+in.Field+" as value").Where(in.Key+" IN (?)", in.Values).Find(&rows.Rows)

	return rows, nil
}

func GetSchemaData(ctx context.Context, in *pb.SchemaParams) (*pb.Response, error) {

	var res *pb.Response = &pb.Response{}

	if in.Type == "form" || in.Type == "grid" || in.Type == "chart" {

		_ = os.WriteFile("lambda/schemas/"+in.Type+"/"+fmt.Sprintf("%d", in.Id)+".json", []byte(in.Schema), 0700)
	}

	return res, nil
}

func GetRoleData(ctx context.Context, in *pb.RoleParams) (*pb.Response, error) {

	var res *pb.Response = &pb.Response{}
	handlers.GetRoleData()

	return res, nil
}
func UploadMySCHEMA(ctx context.Context, in *pb.RoleParams) (*pb.Response, error) {

	var res *pb.Response = &pb.Response{}
	handlers.UploadDBSCHEMA()

	return res, nil
}
