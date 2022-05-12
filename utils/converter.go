package utils

import (
	"fmt"
	"reflect"
)

func GetString(value interface{}) string {
	stringValue := ""

	if reflect.TypeOf(value).String() == "*string" {
		valPre := value.(*string)
		stringValue = *valPre
	} else if reflect.TypeOf(value).String() == "*int" {
		valPre := value.(*int)
		stringValue = fmt.Sprintf("%d", *valPre)
	} else if reflect.TypeOf(value).String() == "int" {
		valPre := value.(int)
		stringValue = fmt.Sprintf("%d", valPre)
	} else {
		stringValue = fmt.Sprintf("%v", value)
	}
	return stringValue
}
