package grpc

import (
	"context"
	"github.com/lambda-platform/lambda/DB"
	pb "github.com/lambda-platform/lambda/grpc/proto"
)


func GetIntData(ctx context.Context, in *pb.TableOption) (*pb.IntRows, error) {

	var rows *pb.IntRows = &pb.IntRows{}



	DB.DB.Table(in.Table).Select(in.Key + " as key, "+in.Field + " as value").Where(in.Key+" IN (?)", in.Values).Find(&rows.Rows)



	return rows, nil
}
func GetStringData(ctx context.Context, in *pb.TableOption) (*pb.StringRows, error) {

	var rows *pb.StringRows = &pb.StringRows{}

	DB.DB.Table(in.Table).Select(in.Key + " as key, "+in.Field + " as value").Where(in.Key+" IN (?)", in.Values).Find(&rows.Rows)

	return rows, nil
}
