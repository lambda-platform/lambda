package grpc

import (
	"context"

	"time"
	pb "github.com/lambda-platform/lambda/grpc/proto"
	"google.golang.org/grpc"
)


func CallStringData(address string, table string, key string, field string, filter []string) (*pb.StringRows, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(500 * time.Millisecond))

	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewLambdaClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetStringData(ctx, &pb.TableOption{Table: table, Key: key, Field: field, Values: filter})
	if err != nil {
		return nil, err
	}
	return r, nil
}