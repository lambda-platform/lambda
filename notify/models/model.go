package models

import "time"

type UserNotifications struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	Link       string     `gorm:"column:link" json:"link"`
	Sender     string     `gorm:"column:sender" json:"sender"`
	Title      string     `gorm:"column:title" json:"title"`
	Body       string     `gorm:"column:body" json:"body"`
	Data       string     `gorm:"column:data" json:"data"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	SID        int64      `gorm:"column:sid" json:"sid"`
	ReceiverID int64      `gorm:"column:receiver_id" json:"receiver_id"`
	Seen       int        `gorm:"column:seen" json:"seen"`
	SeenTime   *time.Time `gorm:"column:seen_time" json:"seen_time"`
	FirstName  string     `gorm:"column:first_name" json:"first_name"`
	Login      string     `gorm:"column:login" json:"login"`
}

type UserNotificationsUUID struct {
	ID         string     `gorm:"column:id" json:"id"`
	Link       string     `gorm:"column:link" json:"link"`
	Sender     string     `gorm:"column:sender" json:"sender"`
	Title      string     `gorm:"column:title" json:"title"`
	Body       string     `gorm:"column:body" json:"body"`
	Data       string     `gorm:"column:data" json:"data"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	SID        string     `gorm:"column:sid" json:"sid"`
	ReceiverID string     `gorm:"column:receiver_id" json:"receiver_id"`
	Seen       int        `gorm:"column:seen" json:"seen"`
	SeenTime   *time.Time `gorm:"column:seen_time" json:"seen_time"`
	FirstName  string     `gorm:"column:first_name" json:"first_name"`
	Login      string     `gorm:"column:login" json:"login"`
}

type Payload struct {
	RegistrationIds []string        `json:"registration_ids"`
	Data            interface{}     `json:"data"`
	Notification    FCMNotification `json:"notification"`
}
type FCMData struct {
	Title       string    `json:"title"`
	Body        string    `gorm:"column:body" json:"body"`
	Sound       string    `json:"sound"`
	Icon        string    `json:"icon"`
	ClickAction string    `gorm:"column:click_action" json:"click_action"`
	Link        string    `gorm:"column:link" json:"link"`
	FirstName   string    `gorm:"column:first_name" json:"first_name"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	ID          int64     `gorm:"column:id" json:"id"`
}
type FCMNotification struct {
	Title       string `json:"title"`
	Body        string `gorm:"column:body" json:"body"`
	Sound       string `json:"sound"`
	Link        string `gorm:"column:link" json:"link"`
	ClickAction string `gorm:"column:click_action" json:"click_action"`
	Icon        string `json:"icon"`
}

type NotificationData struct {
	Users        []int
	Roles        []int
	Data         FCMData         `json:"data"`
	Notification FCMNotification `json:"notification"`
}

type UserNotificationsOracle struct {
	ID         int64      `gorm:"column:ID;primaryKey;autoIncrement;" json:"id"`
	Link       string     `gorm:"column:LINK" json:"link"`
	Sender     string     `gorm:"column:SENDER" json:"sender"`
	Title      string     `gorm:"column:TITLE" json:"title"`
	Body       string     `gorm:"column:BODY" json:"body"`
	CreatedAt  *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	SID        int64      `gorm:"column:SID" json:"sid"`
	ReceiverID int64      `gorm:"column:RECEIVER_ID" json:"receiver_id"`
	Seen       int        `gorm:"column:SEEN" json:"seen"`
	SeenTime   *time.Time `gorm:"column:SEEN_TIME" json:"seen_time"`
	FirstName  string     `gorm:"column:FIRST_NAME" json:"first_name"`
	Login      string     `gorm:"column:LOGIN" json:"login"`
}
