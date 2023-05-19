package utils

import (
	"fmt"
	"reflect"
)

func GetString(value interface{}) string {
	stringValue := ""

	if reflect.TypeOf(value).String() == "*string" {
		valPre := value.(*string)
		if valPre != nil {
			stringValue = *valPre
		} else {
			return ""
		}

	} else if reflect.TypeOf(value).String() == "*int" {
		valPre := value.(*int)
		if valPre != nil {
			stringValue = fmt.Sprintf("%d", *valPre)
		} else {
			return ""
		}
	} else if reflect.TypeOf(value).String() == "int" {
		valPre := value.(int)
		stringValue = fmt.Sprintf("%d", valPre)
	} else {
		stringValue = fmt.Sprintf("%v", value)
	}
	return stringValue
}
