package utils

import (
	"github.com/lambda-platform/lambda/DBSchema"
	pb "github.com/lambda-platform/lambda/grpc/proto"
)


func IndexOf(connectionValue int, grpcRows *pb.StringRows) (int) {
	keyValue := int32(connectionValue)
	for k, v := range grpcRows.Rows {
		if keyValue == v.Key {
			return k
		}
	}
	return -1
}

func IndexOfMicro(MicroserviceID int, microservices []DBSchema.Microservice) (int) {

	for i, m := range microservices {
		if MicroserviceID == m.ProjectID {
			return i
		}
	}
	return -1
}
