package utils

import (
	"errors"
	"reflect"
)

func GetFieldValue(structVar interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(structVar)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("not a struct")
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, errors.New("no such field")
	}

	return field.Interface(), nil
}
