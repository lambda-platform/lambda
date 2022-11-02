package formModels

import (
	"github.com/lambda-platform/lambda/DB"
	"gorm.io/gorm"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type Krud struct {
	ID        int            `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Form      int            `gorm:"column:FORM" json:"form"`
	Grid      int            `gorm:"column:GRID" json:"grid"`
	Template  string         `gorm:"column:TEMPLATE" json:"template"`
	Title     string         `gorm:"column:TITLE" json:"title"`
	Actions   *string        `gorm:"column:ACTIONS" json:"actions"`
	CreatedAt *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

// TableName sets the insert table name for this struct type
func (k *Krud) TableName() string {
	return "KRUD"
}
