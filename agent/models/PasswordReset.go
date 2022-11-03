package models

import (
	"time"
)

type PasswordReset struct {
	Token     string    `gorm:"column:token;not null;" json:"token"`
	Wrong     int       `gorm:"column:wrong;not null;" json:"wrong"`
	Email     string    `gorm:"column:email;primaryKey;" json:"email"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (v *PasswordReset) TableName() string {
	return "password_resets"
}

type PASSWORDRESETSOracle struct {
	Token     string    `gorm:"column:TOKEN;not null;" json:"token"`
	Wrong     int       `gorm:"column:WRONG;not null;" json:"wrong"`
	Email     string    `gorm:"column:EMAIL;primaryKey;" json:"email"`
	CreatedAt time.Time `gorm:"column:CREATED_AT" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (v *PASSWORDRESETSOracle) TableName() string {
	return "PASSWORD_RESETS"
}
