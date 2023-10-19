package utils

import (
	"fmt"
	"reflect"
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

	// Handle pointers by getting their element
	pv := reflect.ValueOf(value)
	if pv.Kind() == reflect.Ptr {
		if pv.IsNil() {
			return ""
		}
		value = pv.Elem().Interface()
	}

	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%g", v)
	case float32:
		return fmt.Sprintf("%g", v)
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	default:
		fmt.Println(fmt.Errorf("unsupported type: %T", v))
		return ""
	}
}

func ConvertToInterfaceSlice(specificSliceInterface interface{}) ([]interface{}, error) {
	var specificSlice reflect.Value

	// Check if it's a pointer and get the underlying element if true
	if reflect.TypeOf(specificSliceInterface).Kind() == reflect.Ptr {
		specificSlice = reflect.ValueOf(specificSliceInterface).Elem()
	} else {
		specificSlice = reflect.ValueOf(specificSliceInterface)
	}

	// Ensure that we're working with a slice
	if specificSlice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a slice, it's a %v", specificSlice.Kind())
	}

	interfaceSlice := make([]interface{}, specificSlice.Len())

	for i := 0; i < specificSlice.Len(); i++ {
		if specificSlice.Index(i).Kind() == reflect.Ptr {
			interfaceSlice[i] = specificSlice.Index(i).Elem().Interface()
		} else {
			interfaceSlice[i] = specificSlice.Index(i).Interface()
		}
	}

	return interfaceSlice, nil
}
