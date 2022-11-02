package models

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID          int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"column:name;not null;unique" json:"name"`
	DisplayName string         `gorm:"column:display_name" json:"display_name"`
	Description string         `gorm:"column:description" json:"description"`
	Permissions string         `gorm:"column:permissions;type:TEXT" json:"permissions"`
	Extra       string         `gorm:"column:extra;type:TEXT" json:"extra"`
	Menu        string         `gorm:"column:menu" json:"menu"`
	CreatedAt   *time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *Role) TableName() string {
	return "roles"
}

type RoleOracle struct {
	ID          int            `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Description string         `gorm:"column:DESCRIPTION" json:"description"`
	DisplayName string         `gorm:"column:DISPLAY_NAME" json:"display_name"`
	Extra       string         `gorm:"column:EXTRA" json:"extra"`
	Menu        *string        `gorm:"column:MENU" json:"menu"`
	Name        string         `gorm:"column:NAME" json:"name"`
	Permissions string         `gorm:"column:PERMISSIONS;type:LONG" json:"permissions"`
	CreatedAt   *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *RoleOracle) TableName() string {
	return "ROLES"
}

type Permissions struct {
	DefaultMenu string `json:"default_menu"`
	Extra       struct {
		Chart       bool `json:"chart"`
		Datasourcce bool `json:"datasourcce"`
		Datasource  bool `json:"datasource"`
		Moqup       bool `json:"moqup"`
	} `json:"extra"`
	MenuID      int                       `json:"menu_id"`
	Permissions map[string]PermissionData `json:"permissions"`
}
type PermissionData struct {
	C      bool   `json:"c"`
	D      bool   `json:"d"`
	MenuID string `json:"menu_id"`
	R      bool   `json:"r"`
	Show   bool   `json:"show"`
	Title  string `json:"title"`
	U      bool   `json:"u"`
}
