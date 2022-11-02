package formModels

import (
	"github.com/lambda-platform/lambda/DB"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type VbSchemas struct {
	ID        int        `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:NAME" json:"name"`
	Schema    string     `gorm:"column:SCHEMA;type:LONG" json:"schema"`
	Type      string     `gorm:"column:TYPE" json:"type"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *VbSchemas) TableName() string {
	return "VB_SCHEMAS"
}
