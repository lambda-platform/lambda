package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Status         string         `gorm:"column:status;" json:"status"` //type:ENUM('0','1','2')
	Role           int64          `gorm:"column:role" json:"role"`
	Login          string         `gorm:"column:login;unique_index;not null;unique" json:"login"`
	Email          string         `gorm:"column:email;unique_index;not null;unique" json:"email"`
	RegisterNumber string         `gorm:"column:register_number;not null;unique" json:"register_number"`
	Avatar         string         `gorm:"column:avatar;type:TEXT" json:"avatar"`
	Bio            string         `gorm:"column:bio;type:TEXT" json:"bio"`
	FirstName      string         `gorm:"column:first_name" json:"first_name"`
	LastName       string         `gorm:"column:last_name" json:"last_name"`
	Birthday       time.Time      `gorm:"column:birthday;type:DATE" json:"birthday"`
	Phone          string         `gorm:"column:phone" json:"phone"`
	Gender         string         `gorm:"column:gender;" json:"gender"` //type:ENUM('f','m')
	Password       string         `gorm:"column:password;not null" json:"password"`
	CreatedAt      *time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *User) TableName() string {
	return "users"
}

type USERSOracle struct {
	ID             int64          `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Avatar         *string        `gorm:"column:AVATAR" json:"avatar"`
	Bio            *string        `gorm:"column:BIO" json:"bio"`
	Birthday       time.Time      `gorm:"column:BIRTHDAY;type:DATE" json:"birthday"`
	Email          string         `gorm:"column:EMAIL" json:"email"`
	FirstName      string         `gorm:"column:FIRST_NAME" json:"first_name"`
	Gender         string         `gorm:"column:GENDER" json:"gender"`
	LastName       *string        `gorm:"column:LAST_NAME" json:"last_name"`
	Login          string         `gorm:"column:LOGIN" json:"login"`
	Password       string         `gorm:"column:PASSWORD" json:"password"`
	Phone          *string        `gorm:"column:PHONE" json:"phone"`
	RegisterNumber string         `gorm:"column:REGISTER_NUMBER" json:"register_number"`
	Role           int64          `gorm:"column:ROLE" json:"role"`
	Status         string         `gorm:"column:STATUS" json:"status"`
	CreatedAt      *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *USERSOracle) TableName() string {
	return "USERS"
}

type UserUUID struct {
	ID             string     `gorm:"column:id;primaryKey;autoIncrement;type:varchar;default:gen_random_uuid()" json:"id"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"-"`
	Status         string     `gorm:"column:status;" json:"status"` //type:ENUM('0','1','2')
	Role           int64      `gorm:"column:role" json:"role"`
	Login          string     `gorm:"column:login;unique_index;not null;unique" json:"login"`
	Email          string     `gorm:"column:email;unique_index;not null;unique" json:"email"`
	RegisterNumber string     `gorm:"column:register_number;not null;unique" json:"register_number"`
	Avatar         string     `gorm:"column:avatar;type:TEXT" json:"avatar"`
	Bio            string     `gorm:"column:bio;type:TEXT" json:"bio"`
	FirstName      string     `gorm:"column:first_name" json:"first_name"`
	LastName       string     `gorm:"column:last_name" json:"last_name"`
	Birthday       time.Time  `gorm:"column:birthday;type:DATE" json:"birthday"`
	Phone          string     `gorm:"column:phone" json:"phone"`
	Gender         string     `gorm:"column:gender;" json:"gender"` //type:ENUM('f','m')
	Password       string     `gorm:"column:password;not null" json:"password"`
}

// TableName sets the insert table name for this struct type
func (v *UserUUID) TableName() string {
	return "users"
}

type UserWithoutPassword struct {
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Status         string         `gorm:"column:status;" json:"status"` //type:ENUM('0','1','2')
	Role           int64          `gorm:"column:role" json:"role"`
	Login          string         `gorm:"column:login;unique_index;not null;unique" json:"login"`
	Email          string         `gorm:"column:email;unique_index;not null;unique" json:"email"`
	RegisterNumber string         `gorm:"column:register_number;not null;unique" json:"register_number"`
	Avatar         string         `gorm:"column:avatar;type:TEXT" json:"avatar"`
	Bio            string         `gorm:"column:bio;type:TEXT" json:"bio"`
	FirstName      string         `gorm:"column:first_name" json:"first_name"`
	LastName       string         `gorm:"column:last_name" json:"last_name"`
	Birthday       time.Time      `gorm:"column:birthday;type:DATE" json:"birthday"`
	Phone          string         `gorm:"column:phone" json:"phone"`
	Gender         string         `gorm:"column:gender;" json:"gender"` //type:ENUM('f','m')
	CreatedAt      *time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *UserWithoutPassword) TableName() string {
	return "users"
}

type UserWithoutPasswordOracle struct {
	ID             int64          `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Avatar         *string        `gorm:"column:AVATAR" json:"avatar"`
	Bio            *string        `gorm:"column:BIO" json:"bio"`
	Birthday       time.Time      `gorm:"column:BIRTHDAY;type:DATE" json:"birthday"`
	Email          string         `gorm:"column:EMAIL" json:"email"`
	FirstName      *string        `gorm:"column:FIRST_NAME" json:"first_name"`
	Gender         string         `gorm:"column:GENDER" json:"gender"`
	LastName       *string        `gorm:"column:LAST_NAME" json:"last_name"`
	Login          string         `gorm:"column:LOGIN" json:"login"`
	Phone          *string        `gorm:"column:PHONE" json:"phone"`
	RegisterNumber string         `gorm:"column:REGISTER_NUMBER" json:"register_number"`
	Role           int64          `gorm:"column:ROLE" json:"role"`
	Status         string         `gorm:"column:STATUS" json:"status"`
	CreatedAt      *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *UserWithoutPasswordOracle) TableName() string {
	return "USERS"
}

type UserWithoutPasswordUUID struct {
	ID             string         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Status         string         `gorm:"column:status;" json:"status"` //type:ENUM('0','1','2')
	Role           int64          `gorm:"column:role" json:"role"`
	Login          string         `gorm:"column:login;unique_index;not null;unique" json:"login"`
	Email          string         `gorm:"column:email;unique_index;not null;unique" json:"email"`
	RegisterNumber string         `gorm:"column:register_number;not null;unique" json:"register_number"`
	Avatar         string         `gorm:"column:avatar;type:TEXT" json:"avatar"`
	Bio            string         `gorm:"column:bio;type:TEXT" json:"bio"`
	FirstName      string         `gorm:"column:first_name" json:"first_name"`
	LastName       string         `gorm:"column:last_name" json:"last_name"`
	Birthday       time.Time      `gorm:"column:birthday;type:DATE" json:"birthday"`
	Phone          string         `gorm:"column:phone" json:"phone"`
	Gender         string         `gorm:"column:gender;" json:"gender"` //type:ENUM('f','m')
	CreatedAt      *time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *UserWithoutPasswordUUID) TableName() string {
	return "users"
}
