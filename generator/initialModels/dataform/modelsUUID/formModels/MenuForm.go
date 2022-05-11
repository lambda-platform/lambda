package formModels

import (
	"github.com/lambda-platform/lambda/DB"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type VbSchemas struct {
	ID        int        `gorm:"column:id;primary_key" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Schema    string     `gorm:"column:schema" json:"schema"`
	Type      string     `gorm:"column:type" json:"type"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

//  TableName sets the insert table name for this struct type
func (v *VbSchemas) TableName() string {
	return "vb_schemas"
}
