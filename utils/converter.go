package utils

import (
	"fmt"
	"strconv"
)

//	func GetString(value interface{}) string {
//		stringValue := ""
//
//		if reflect.TypeOf(value).String() == "*string" {
//			valPre := value.(*string)
//			if valPre != nil {
//				stringValue = *valPre
//			} else {
//				return ""
//			}
//
//		} else if reflect.TypeOf(value).String() == "*int" {
//			valPre := value.(*int)
//			if valPre != nil {
//				stringValue = fmt.Sprintf("%d", *valPre)
//			} else {
//				return ""
//			}
//		} else if reflect.TypeOf(value).String() == "int" {
//			valPre := value.(int)
//			stringValue = fmt.Sprintf("%d", valPre)
//		} else {
//			stringValue = fmt.Sprintf("%v", value)
//		}
//		return stringValue
//	}
func GetString(value interface{}) string {
	if value == nil {
		return ""
	}

	var valueStr string
	switch v := value.(type) {
	case float64:
		valueStr = fmt.Sprintf("%g", v)
	case float32:
		valueStr = fmt.Sprintf("%g", float64(v))
	case int:
		valueStr = strconv.Itoa(v)
	case int32:
		valueStr = strconv.FormatInt(int64(v), 10)
	case int64:
		valueStr = strconv.FormatInt(v, 10)
	case string:
		valueStr = v
	default:
		// handle other types as you see fit, or return an error
		valueStr = ""
	}
	return valueStr
}
