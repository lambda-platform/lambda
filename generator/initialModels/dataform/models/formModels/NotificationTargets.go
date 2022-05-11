package formModels

import (
	"github.com/lambda-platform/lambda/DB"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type NotificationTargets struct {
	Body          string `gorm:"column:body" json:"body"`
	Condition     string `gorm:"column:condition" json:"condition"`
	ID            int    `gorm:"column:id;primary_key" json:"id"`
	Link          string `gorm:"column:link" json:"link"`
	SchemaID      int    `gorm:"column:schema_id" json:"schema_id"`
	TargetActions string `gorm:"column:target_actions" json:"target_actions"`
	TargetRole    int    `gorm:"column:target_role" json:"target_role"`
	Title         string `gorm:"column:title" json:"title"`
}

//  TableName sets the insert table name for this struct type
func (n *NotificationTargets) TableName() string {
	return "notification_targets"
}
