package utils

import (
	"github.com/lambda-platform/lambda/models"
	"fmt"
	"reflect"
	pb "github.com/lambda-platform/lambda/grpc/proto"
)


func IndexOf(connectionValue interface{}, grpcRows *pb.StringRows) (int) {
	keyValue := ""

	if reflect.TypeOf(connectionValue).String() == "*string" {
		valPre := connectionValue.(*string)
		keyValue = *valPre
	} else if reflect.TypeOf(connectionValue).String() == "*int" {
		valPre := connectionValue.(*int)
		keyValue = fmt.Sprintf("%d", *valPre)
	}  else if reflect.TypeOf(connectionValue).String() == "int" {
		valPre := connectionValue.(int)
		keyValue = fmt.Sprintf("%d", valPre)
	} else {
		keyValue = fmt.Sprintf("%v", connectionValue)
	}


	for k, v := range grpcRows.Rows {
		if keyValue == v.Key {
			return k
		}
	}
	return -1
}

func IndexOfMicro(MicroserviceID int, microservices []models.Microservice) (int) {

	for i, m := range microservices {
		if MicroserviceID == m.ProjectID {
			return i
		}
	}
	return -1
}
