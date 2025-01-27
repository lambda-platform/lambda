package grid

import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/models"
	"gorm.io/gorm"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}
var _ = gorm.DB{}

var KrudGridDatagrid datagrid.Datagrid = datagrid.Datagrid{
	Name:      "Крүд тохиргоо",
	Identity:  "ID",
	DataTable: "KRUD",
	MainTable: "KRUD",
	DataModel: new(KrudGrid),
	Data:      new([]KrudGrid),
	MainModel: new(KrudGridMain),
	Columns: []datagrid.Column{
		datagrid.Column{
			Model: "title",
			Label: "Гарчиг",
		},
	},
	ColumnList:  []string{"ID", "TITLE"},
	Filters:     map[string]string{},
	Relations:   []models.GridRelation{},
	Condition:   "",
	Aggregation: "",
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

type KrudGrid struct {
	ID        int            `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Actions   *string        `gorm:"column:ACTIONS" json:"actions"`
	Form      *int           `gorm:"column:FORM" json:"form"`
	Grid      *int           `gorm:"column:GRID" json:"grid"`
	Template  *string        `gorm:"column:TEMPLATE" json:"template"`
	Title     *string        `gorm:"column:TITLE" json:"title"`
	CreatedAt *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

func (k *KrudGrid) TableName() string {
	return "KRUD"
}

type KrudGridMain struct {
	ID        int            `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Actions   *string        `gorm:"column:ACTIONS" json:"actions"`
	Form      *int           `gorm:"column:FORM" json:"form"`
	Grid      *int           `gorm:"column:GRID" json:"grid"`
	Template  *string        `gorm:"column:TEMPLATE" json:"template"`
	Title     *string        `gorm:"column:TITLE" json:"title"`
	CreatedAt *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

func (k *KrudGridMain) TableName() string {
	return "KRUD"
}
