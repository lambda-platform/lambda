package formModels

import (
	"github.com/lambda-platform/lambda/DB"
	"gorm.io/gorm"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type Krud struct {
	Actions   string         `gorm:"column:actions" json:"actions"`
	CreatedAt *time.Time     `gorm:"column:created_at" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Form      int            `gorm:"column:form" json:"form"`
	Grid      int            `gorm:"column:grid" json:"grid"`
	ID        int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Template  string         `gorm:"column:template" json:"template"`
	Title     string         `gorm:"column:title" json:"title"`
	UpdatedAt *time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

//  TableName sets the insert table name for this struct type
func (k *Krud) TableName() string {
	return "krud"
}
