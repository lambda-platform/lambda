package formModels

import (
	"github.com/lambda-platform/lambda/DB"
	"gorm.io/gorm"
	"time"
)

var _ = time.Time{}
var _ = DB.Date{}

type Users struct {
	FcmToken  string     `gorm:"column:FCM_TOKEN" json:"fcm_token"`
	ID        int        `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Password  string     `gorm:"column:PASSWORD" json:"password"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "users"
}

type UserForm struct {
	ID             int64          `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Avatar         *string        `gorm:"column:AVATAR" json:"avatar"`
	Bio            *string        `gorm:"column:BIO" json:"bio"`
	Birthday       DB.Date        `gorm:"column:BIRTHDAY;type:DATE" json:"birthday"`
	Email          string         `gorm:"column:EMAIL" json:"email"`
	FcmToken       *string        `gorm:"column:FCM_TOKEN" json:"fcm_token"`
	FirstName      *string        `gorm:"column:FIRST_NAME" json:"first_name"`
	Gender         string         `gorm:"column:GENDER" json:"gender"`
	LastName       *string        `gorm:"column:LAST_NAME" json:"last_name"`
	Login          string         `gorm:"column:LOGIN" json:"login"`
	Password       string         `gorm:"column:PASSWORD" json:"password"`
	Phone          *string        `gorm:"column:PHONE" json:"phone"`
	RegisterNumber string         `gorm:"column:REGISTER_NUMBER" json:"register_number"`
	Role           int64          `gorm:"column:ROLE" json:"role"`
	CreatedAt      *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

// TableName sets the insert table name for this struct type
func (u *UserForm) TableName() string {
	return "USERS"
}

type UserProfile struct {
	ID             int64   `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Avatar         *string `gorm:"column:AVATAR" json:"avatar"`
	Bio            *string `gorm:"column:BIO" json:"bio"`
	Birthday       DB.Date `gorm:"column:BIRTHDAY;type:DATE" json:"birthday"`
	Email          string  `gorm:"column:EMAIL" json:"email"`
	FcmToken       *string `gorm:"column:FCM_TOKEN" json:"fcm_token"`
	FirstName      *string `gorm:"column:FIRST_NAME" json:"first_name"`
	Gender         string  `gorm:"column:GENDER" json:"gender"`
	LastName       *string `gorm:"column:LAST_NAME" json:"last_name"`
	Login          string  `gorm:"column:LOGIN" json:"login"`
	Password       string  `gorm:"column:PASSWORD" json:"password"`
	Phone          *string `gorm:"column:PHONE" json:"phone"`
	RegisterNumber string  `gorm:"column:REGISTER_NUMBER" json:"register_number"`
}

// TableName sets the insert table name for this struct type
func (u *UserProfile) TableName() string {
	return "USERS"
}
