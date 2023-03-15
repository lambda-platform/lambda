package models

import "time"

type Notification struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Link      string    `gorm:"column:link" json:"link"`
	Sender    uint      `gorm:"column:sender" json:"sender"`
	Title     string    `gorm:"column:title" json:"title"`
	Body      string    `gorm:"column:body" json:"body"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (n *Notification) TableName() string {
	return "notifications"
}

type NotificationUUID struct {
	ID        string    `gorm:"column:id;primaryKey;autoIncrement;default:gen_random_uuid();type:uuid" json:"id"`
	Link      string    `gorm:"column:link" json:"link"`
	Sender    string    `gorm:"column:sender" json:"sender"`
	Title     string    `gorm:"column:title" json:"title"`
	Body      string    `gorm:"column:body" json:"body"`
	Data      string    `gorm:"column:data" json:"data"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (n *NotificationUUID) TableName() string {
	return "notifications"
}
