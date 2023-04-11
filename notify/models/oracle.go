package models

import "time"

type NotificationTargetOracle struct {
	ID            int64  `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Link          string `gorm:"column:LINK" json:"link"`
	Title         string `gorm:"column:TITLE" json:"title"`
	Body          string `gorm:"column:BODY" json:"body"`
	SchemaId      int    `gorm:"column:SCHEMA_ID" json:"schema_id"`
	TargetRole    int    `gorm:"column:TARGET_ROLE" json:"target_role"`
	Condition     string `gorm:"column:CONDITION" json:"condition"`
	TargetActions string `gorm:"column:TARGET_ACTIONS" json:"target_actions"`
}

func (n *NotificationTargetOracle) TableName() string {
	return "NOTIFICATION_TARGETS"
}

type NotificationOracle struct {
	ID        int64     `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Link      string    `gorm:"column:LINK" json:"link"`
	Sender    uint      `gorm:"column:SENDER" json:"sender"`
	Title     string    `gorm:"column:TITLE" json:"title"`
	Body      string    `gorm:"column:BODY" json:"body"`
	Data      string    `gorm:"column:DATA" json:"data"`
	CreatedAt time.Time `gorm:"column:CREATED_AT" json:"created_at"`
}

func (n *NotificationOracle) TableName() string {
	return "NOTIFICATIONS"
}

type UserFcmTokensOracle struct {
	ID       int64  `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	UserID   int    `gorm:"column:USER_ID" json:"user_id"`
	FcmToken string `gorm:"column:FCM_TOKEN" json:"fcm_token"`
}

func (u *UserFcmTokensOracle) TableName() string {
	return "NOTIFICATION_USER_TOKENS"
}

type NotificationStatusOracle struct {
	ID         int64     `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	NotifID    int64     `gorm:"column:NOTIF_ID" json:"notif_id"`
	ReceiverID int64     `gorm:"column:RECEIVER_ID" json:"receiver_id"`
	Seen       int       `gorm:"column:SEEN" json:"seen"`
	SeenTime   time.Time `gorm:"column:SEEN_TIME" json:"seen_time"`
}

func (n *NotificationStatusOracle) TableName() string {
	return "NOTIFICATION_STATUS"
}
