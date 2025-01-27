package dataform

import (
	"fmt"
	"github.com/lambda-platform/lambda/models"
	"reflect"
	"strconv"
)

type Dataform struct {
	Name               string
	Identity           string
	Table              string
	Model              interface{}
	FieldTypes         map[string]string
	Formulas           []models.Formula
	ValidationRules    map[string][]string
	ValidationMessages map[string][]string
	SubForms           []map[string]interface{}
	BeforeInsert       func(interface{})
	BeforeUpdate       func(interface{})
	AfterInsert        func(interface{})
	AfterUpdate        func(interface{})
	TriggerNameSpace   string
	Relations          map[string]models.Relation
}

func (d *Dataform) getStringField(field string) (string, error) {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		return string(f.String()), nil
	} else {
		return "", fmt.Errorf("Field not found: " + field)
	}

}

func (d *Dataform) getFieldType(field string) (string, error) {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		return f.Type().String(), nil
	} else {
		return "", fmt.Errorf("Field not found: " + field)
	}
}
func (d *Dataform) setStringField(field string, value string) error {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		f.SetString(value)
		return nil
	} else {
		return fmt.Errorf("Field not found: " + field)
	}

}
func (d *Dataform) setIntField(field string, value int) error {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	if f.IsValid() {
		f.SetInt(int64(value))
		return nil
	} else {
		return fmt.Errorf("Field not found: " + field)
	}
}

func (d *Dataform) getIntField(field string) (int, error) {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return 0, fmt.Errorf("Field not found: " + field)
	}
	return int(f.Int()), nil
}
func (d *Dataform) getFieldValue(field string) (interface{}, error) {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil, fmt.Errorf("Field not found: " + field)
	}
	return f.Interface(), nil
}
func (d *Dataform) getModelFieldValue(Model interface{}, field string) (interface{}, error) {
	r := reflect.ValueOf(Model)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil, fmt.Errorf("Field not found: " + field)
	}
	return f.Interface(), nil
}
func (d *Dataform) setModelField(Model interface{}, field string, value interface{}) error {
	r := reflect.ValueOf(Model)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return fmt.Errorf("Field not found: " + field)
	}
	switch vtype := value.(type) {
	case string:
		f.SetString(value.(string))
	default:
		fmt.Println(vtype)
		f.SetInt(int64(GetInt(value)))
	}
	return nil

}
func GetInt(value interface{}) int {
	intValue := 0

	switch v := value.(type) {
	case int:
		intValue = value.(int)
	case int64:
		intValue = int(value.(int64))
	case int32:
		intValue = int(value.(int32))
	case float64:
		intValue = int(value.(float64))
	case float32:
		intValue = int(value.(float32))
	case string:
		i, _ := strconv.Atoi(value.(string))
		intValue = i
	default:
		fmt.Println(v)
	}
	return intValue
}

func (d *Dataform) setModelFieldValue(Model interface{}, field string) (interface{}, error) {
	r := reflect.ValueOf(Model)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return nil, fmt.Errorf("Field not found: " + field)
	}
	return f.Interface(), nil
}
func Clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

//
//func (d *Dataform) getFloatField(field string) float64 {
//	r := reflect.ValueOf(d.Model)
//	f := reflect.Indirect(r).FieldByName(field)
//	return float64(f.Float())
//}
//
//func (d *Dataform) getInterfaceField(field string) interface{} {
//	r := reflect.ValueOf(d.Model)
//	f := reflect.Indirect(r).FieldByName(field)
//	return f.Interface().(interface{})
//}
//

type Relations struct {
	Relations map[string]Ralation_ `json:"relations"`
}

type Ralation_ struct {
	Fields             []string            `json:"Fields"`
	FilterWithUser     []map[string]string `json:"filterWithUser"`
	Filter             string              `json:"filter"`
	Key                string              `json:"key"`
	Multiple           bool                `json:"multiple"`
	ParentFieldOfForm  string              `json:"parentFieldOfForm"`
	ParentFieldOfTable string              `json:"parentFieldOfTable"`
	SortField          string              `json:"sortField"`
	SortOrder          string              `json:"sortOrder"`
	Table              string              `json:"table"`
}

type RalationOption struct {
	Fields             []string            `json:"Fields"`
	FilterWithUser     []map[string]string `json:"filterWithUser"`
	Filter             string              `json:"filter"`
	Key                string              `json:"key"`
	SortField          string              `json:"sortField"`
	SortOrder          string              `json:"sortOrder"`
	Table              string              `json:"table"`
	ParentFieldOfForm  string              `json:"parentFieldOfForm"`
	ParentFieldOfTable string              `json:"parentFieldOfTable"`
}
