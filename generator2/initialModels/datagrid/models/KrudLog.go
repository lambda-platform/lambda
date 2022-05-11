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

type CrudLog struct {
	Action    string     `gorm:"column:action" json:"action"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	FirstName string     `gorm:"column:first_name" json:"first_name"`
	ID        int        `gorm:"column:id" json:"id"`
	Input     string     `gorm:"column:input" json:"input"`
	IP        string     `gorm:"column:ip" json:"ip"`
	LastName  string     `gorm:"column:last_name" json:"last_name"`
	Name      string     `gorm:"column:name" json:"name"`
	RowID     string     `gorm:"column:row_id" json:"row_id"`
	UserAgent string     `gorm:"column:user_agent" json:"user_agent"`
}

type CrudLogMainTable struct {
	Action    string     `gorm:"column:action" json:"action"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	ID        int        `gorm:"column:id;primary_key" json:"id"`
	Input     string     `gorm:"column:input" json:"input"`
	IP        string     `gorm:"column:ip" json:"ip"`
	RowID     string     `gorm:"column:row_id" json:"row_id"`
	SchemaID  int        `gorm:"column:schemaId" json:"schemaId"`
	UserAgent string     `gorm:"column:user_agent" json:"user_agent"`
	UserID    int        `gorm:"column:user_id" json:"user_id"`
}

func (m *CrudLogMainTable) TableName() string {
	return "ds_crud_log"
}

//  TableName sets the insert table name for this struct type
func (d *CrudLog) TableName() string {
	return "crud_log"
}

var CrudLogDatagrid datagrid.Datagrid = datagrid.Datagrid{
	Name:      "Cистем лог",
	Identity:  "id",
	DataTable: "ds_crud_log",
	MainTable: "crud_log",
	DataModel: new(CrudLog),
	Data:      new([]CrudLog),
	MainModel: new(CrudLogMainTable),
	Columns: []datagrid.Column{
		datagrid.Column{Model: "name", Label: "Маягт"},
		datagrid.Column{Model: "last_name", Label: "Хэрэглэгчийн овог"},
		datagrid.Column{Model: "first_name", Label: "Хэрэглэгчийн нэр"},
		datagrid.Column{Model: "ip", Label: "IP хаяг"},
		datagrid.Column{Model: "user_agent", Label: "хэрэглэгчийн агент"},
		datagrid.Column{Model: "action", Label: "Үйлдэл"},
		datagrid.Column{Model: "row_id", Label: "бүртгэлийн ID"},
		datagrid.Column{Model: "input", Label: "Оролтын мэдээлэл"},
		datagrid.Column{Model: "created_at", Label: "Огноо"},
	},
	ColumnList: []string{"id", "name", "last_name", "first_name", "ip", "user_agent", "action", "row_id", "input", "created_at"},
	Filters: map[string]string{
		"user_id": "Select",

		"schemaId": "Select",
	},
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
