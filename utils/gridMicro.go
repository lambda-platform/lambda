package utils

import (
	pb "github.com/lambda-platform/lambda/grpc/proto"
	"github.com/lambda-platform/lambda/models"
)

func IndexOf(connectionValue interface{}, grpcRows *pb.StringRows) int {

	keyValue := GetString(connectionValue)

	for k, v := range grpcRows.Rows {
		if keyValue == v.Key {
			return k
		}
	}
	return -1
}

func IndexOfMicro(MicroserviceID int, microservices []models.Microservice) int {

	for i, m := range microservices {
		if MicroserviceID == m.ProjectID {
			return i
		}
	}
	return -1
}
