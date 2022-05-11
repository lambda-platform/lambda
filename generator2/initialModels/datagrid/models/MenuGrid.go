package grid

import (
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/models"
)

import (
	"github.com/lambda-platform/lambda/DB"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type MenuGrid struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Schema    string    `gorm:"column:schema" json:"schema"`
	Type      string    `gorm:"column:type" json:"type"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

//  TableName sets the insert table name for this struct type
func (v *MenuGrid) TableName() string {
	return "vb_schemas"
}

type MenuGridMain struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	ID        int64     `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Schema    string    `gorm:"column:schema" json:"schema"`
	Type      string    `gorm:"column:type" json:"type"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

//  TableName sets the insert table name for this struct type
func (v *MenuGridMain) TableName() string {
	return "vb_schemas"
}

var MenuGridDatagrid datagrid.Datagrid = datagrid.Datagrid{
	Name:      "Цэсний тохиргоо",
	Identity:  "id",
	DataTable: "vb_schemas",
	MainTable: "vb_schemas",
	DataModel: new(MenuGrid),
	Data:      new([]MenuGrid),
	MainModel: new(MenuGridMain),
	Columns: []datagrid.Column{
		datagrid.Column{Model: "name", Label: "Нэр"},
		datagrid.Column{Model: "created_at", Label: "Огноо"},
	},
	ColumnList:  []string{"id", "name", "created_at"},
	Filters:     map[string]string{},
	Relations:   []models.GridRelation{},
	Condition:   "type = 'menu'",
	Aggergation: "",
	Triggers: map[string]interface{}{
		"beforeFetch":        "",
		"beforeFetchStruct":  new(interface{}),
		"afterFetch":         "",
		"afterFetchStruct":   new(interface{}),
		"beforeDelete":       "",
		"beforeDeleteStruct": new(interface{}),
		"afterDelete":        "",
		"afterDeleteStruct":  new(interface{}),
		"beforePrint":        "",
		"beforePrintStruct":  new(interface{}),
	},
	TriggerNameSpace: "",
}
