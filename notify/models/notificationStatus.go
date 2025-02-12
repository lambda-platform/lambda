package models

import "time"

type NotificationStatus struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NotifID    int64     `gorm:"column:notif_id" json:"notif_id"`
	ReceiverID int64     `gorm:"column:receiver_id" json:"receiver_id"`
	Seen       int       `gorm:"column:seen" json:"seen"`
	SeenTime   time.Time `gorm:"column:seen_time" json:"seen_time"`
}

func (n *NotificationStatus) TableName() string {
	return "notification_status"
}

type NotificationStatusUUID struct {
	ID         string    `gorm:"column:id;primaryKey;type:varchar;default:gen_random_uuid()" json:"id"`
	NotifID    string    `gorm:"column:notif_id;type:uuid" json:"notif_id"`
	ReceiverID string    `gorm:"column:receiver_id;type:uuid" json:"receiver_id"`
	Seen       int       `gorm:"column:seen" json:"seen"`
	SeenTime   time.Time `gorm:"column:seen_time" json:"seen_time"`
}

func (n *NotificationStatusUUID) TableName() string {
	return "notification_status"
}
