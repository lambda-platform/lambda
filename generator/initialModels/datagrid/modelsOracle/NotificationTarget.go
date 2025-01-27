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
	ID    int    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Title string `gorm:"column:TITLE" json:"title"`
}

type NotificationTargetMainTable struct {
	Body      string `gorm:"column:BODY" json:"body"`
	Condition string `gorm:"column:CONDITION" json:"condition"`
	ID        int    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Link      string `gorm:"column:LINK" json:"link"`
	SchemaID  int    `gorm:"column:SCHEMA_ID" json:"schema_id"`

	TargetRole int    `gorm:"column:TARGET_ROLE" json:"target_role"`
	Title      string `gorm:"column:TITLE" json:"title"`
}

func (n *NotificationTargetMainTable) TableName() string {
	return "notification_targets"
}

// TableName sets the insert table name for this struct type
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
