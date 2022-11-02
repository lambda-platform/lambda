package formModels

import (
	"github.com/lambda-platform/lambda/DB"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type NotificationTargets struct {
	Body          string `gorm:"column:BODY" json:"body"`
	Condition     string `gorm:"column:CONDITION" json:"condition"`
	ID            int    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Link          string `gorm:"column:LINK" json:"link"`
	SchemaID      int    `gorm:"column:SCHEMA_ID" json:"schema_id"`
	TargetActions string `gorm:"column:TARGET_ACTIONS" json:"target_actions"`
	TargetRole    int    `gorm:"column:TARGET_ROLE" json:"target_role"`
	Title         string `gorm:"column:TITLE" json:"title"`
}

// TableName sets the insert table name for this struct type
func (n *NotificationTargets) TableName() string {
	return "NOTIFICATION_TARGETS"
}
