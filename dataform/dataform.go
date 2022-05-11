package dataform

import (
	"github.com/lambda-platform/lambda/models"
	"reflect"
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
}

//
func (d *Dataform) getStringField(field string) string {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func (d *Dataform) getFieldType(field string) string {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Type().String()
}
func (d *Dataform) setStringField(field string, value string) {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	f.SetString(value)
}
func (d *Dataform) setIntField(field string, value int) {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	f.SetInt(int64(value))
}

func (d *Dataform) getIntField(field string) int {
	r := reflect.ValueOf(d.Model)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
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
