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
	ID        int        `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:NAME" json:"name"`
	Schema    string     `gorm:"column:SCHEMA;type:LONG" json:"schema"`
	Type      string     `gorm:"column:TYPE" json:"type"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *MenuGrid) TableName() string {
	return "VB_SCHEMAS"
}

type MenuGridMain struct {
	ID        int        `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:NAME" json:"name"`
	Schema    string     `gorm:"column:SCHEMA;type:LONG" json:"schema"`
	Type      string     `gorm:"column:TYPE" json:"type"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *MenuGridMain) TableName() string {
	return "VB_SCHEMAS"
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
