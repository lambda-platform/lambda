package utils

import pb "github.com/lambda-platform/lambda/grpc/proto"


func IndexOf(connectionValue int, grpcRows *pb.StringRows) (int) {
	keyValue := int32(connectionValue)
	for k, v := range grpcRows.Rows {
		if keyValue == v.Key {
			return k
		}
	}
	return -1
}
