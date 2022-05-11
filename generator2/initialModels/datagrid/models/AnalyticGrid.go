package grid

import (
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/datagrid"
	"github.com/lambda-platform/lambda/models"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

var AnalyticGridDatagrid datagrid.Datagrid = datagrid.Datagrid{
	Name:      "АНАЛИЗ",
	Identity:  "id",
	DataTable: "analytic",
	MainTable: "analytic",
	DataModel: new(AnalyticGrid),
	Data:      new([]AnalyticGrid),
	MainModel: new(AnalyticMainTable),
	Columns: []datagrid.Column{
		datagrid.Column{
			Model: "id",
			Label: "",
		},
		datagrid.Column{
			Model: "title",
			Label: "Нэр",
		},
		datagrid.Column{
			Model: "source",
			Label: "Хүснэгт",
		},
	},
	ColumnList:  []string{"id", "title", "source"},
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

type AnalyticGrid struct {
	ID     int    `gorm:"column:id;primary_key" json:"id"`
	Source string `gorm:"column:source" json:"source"`
	Title  string `gorm:"column:title" json:"title"`
}

type AnalyticMainTable struct {
	ID     int    `gorm:"column:id;primary_key" json:"id"`
	Source string `gorm:"column:source" json:"source"`
	Title  string `gorm:"column:title" json:"title"`
}

func (a *AnalyticMainTable) TableName() string {
	return "analytic"
}

//  TableName sets the insert table name for this struct type
func (a *AnalyticGrid) TableName() string {
	return "analytic"
}
