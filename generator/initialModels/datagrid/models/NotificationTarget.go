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

type NotificationTarget struct {
	ID    int64  `gorm:"column:id;primary_key" json:"id"`
	Title string `gorm:"column:title" json:"title"`
}

type NotificationTargetMainTable struct {
	Body       string `gorm:"column:body" json:"body"`
	Condition  string `gorm:"column:condition" json:"condition"`
	ID         int64  `gorm:"column:id;primary_key" json:"id"`
	Link       string `gorm:"column:link" json:"link"`
	SchemaID   int    `gorm:"column:schema_id" json:"schema_id"`
	TargetRole int    `gorm:"column:target_role" json:"target_role"`
	Title      string `gorm:"column:title" json:"title"`
}

func (n *NotificationTargetMainTable) TableName() string {
	return "notification_targets"
}

//  TableName sets the insert table name for this struct type
func (n *NotificationTarget) TableName() string {
	return "notification_targets"
}

var NotificationTargetDatagrid datagrid.Datagrid = datagrid.Datagrid{
	Name:      "Зорилтод мэдэгдэл",
	Identity:  "id",
	DataTable: "notification_targets",
	MainTable: "notification_targets",
	DataModel: new(NotificationTarget),
	Data:      new([]NotificationTarget),
	MainModel: new(NotificationTargetMainTable),
	Columns: []datagrid.Column{
		datagrid.Column{Model: "title", Label: "Нэр"},
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
