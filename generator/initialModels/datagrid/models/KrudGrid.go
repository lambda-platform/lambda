package grid

import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/models"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

var KrudGridDatagrid datagrid.Datagrid = datagrid.Datagrid{
	Name:      "Крүд тохиргоо",
	Identity:  "id",
	DataTable: "krud",
	MainTable: "krud",
	DataModel: new(KrudGrid),
	Data:      new([]KrudGrid),
	MainModel: new(KrudGridMain),
	Columns: []datagrid.Column{
		datagrid.Column{
			Model: "title",
			Label: "Гарчиг",
		},
	},
	ColumnList:  []string{"id", "title"},
	Filters:     map[string]string{},
	Relations:   []models.GridRelation{},
	Condition:   "",
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

type KrudGrid struct {
	ID        int        `gorm:"column:id;primary_key" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Form      int        `gorm:"column:form" json:"form"`
	Grid      int        `gorm:"column:grid" json:"grid"`
	Actions   string     `gorm:"column:actions" json:"actions"`
	Template  string     `gorm:"column:template" json:"template"`
	Title     string     `gorm:"column:title" json:"title"`
}

func (k *KrudGrid) TableName() string {
	return "krud"
}

type KrudGridMain struct {
	ID        int        `gorm:"column:id;primary_key" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Form      int        `gorm:"column:form" json:"form"`
	Grid      int        `gorm:"column:grid" json:"grid"`
	Actions   string     `gorm:"column:actions" json:"actions"`
	Template  string     `gorm:"column:template" json:"template"`
	Title     string     `gorm:"column:title" json:"title"`
}

func (k *KrudGridMain) TableName() string {
	return "krud"
}
