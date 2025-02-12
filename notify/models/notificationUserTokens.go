package models

type UserFcmTokens struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID   int    `gorm:"column:user_id" json:"user_id"`
	FcmToken string `gorm:"column:fcm_token" json:"fcm_token"`
}

func (u *UserFcmTokens) TableName() string {
	return "notification_user_tokens"
}

type UserFcmTokensUUID struct {
	ID       string `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID   string `gorm:"column:user_id;type:uuid" json:"user_id"`
	FcmToken string `gorm:"column:fcm_token" json:"fcm_token"`
}

func (u *UserFcmTokensUUID) TableName() string {
	return "notification_user_tokens"
}
