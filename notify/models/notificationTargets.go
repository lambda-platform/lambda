package models

type NotificationTarget struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Link          string `gorm:"column:link" json:"link"`
	Title         string `gorm:"column:title" json:"title"`
	Body          string `gorm:"column:body" json:"body"`
	SchemaId      int    `gorm:"column:schema_id" json:"schema_id"`
	TargetRole    int    `gorm:"column:target_role" json:"target_role"`
	Condition     string `gorm:"column:condition" json:"condition"`
	TargetActions string `gorm:"column:target_actions" json:"target_actions"`
}

func (n *NotificationTarget) TableName() string {
	return "notification_targets"
}
