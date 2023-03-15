package models

import (
	"gorm.io/gorm"
	"time"
)

type Krud struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title     string         `gorm:"column:title" json:"title"`
	Template  string         `gorm:"column:template" json:"template"`
	Grid      int            `gorm:"column:grid" json:"grid"`
	Form      int            `gorm:"column:form" json:"form"`
	Actions   string         `gorm:"column:actions" json:"actions"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *Krud) TableName() string {
	return "krud"
}

type KrudOracle struct {
	ID        int            `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Actions   *string        `gorm:"column:ACTIONS" json:"actions"`
	Form      *int           `gorm:"column:FORM" json:"form"`
	Grid      *int           `gorm:"column:GRID" json:"grid"`
	Template  *string        `gorm:"column:TEMPLATE" json:"template"`
	Title     *string        `gorm:"column:TITLE" json:"title"`
	CreatedAt *time.Time     `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"column:UPDATED_AT" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"-"`
}

// TableName sets the insert table name for this struct type
func (v *KrudOracle) TableName() string {
	return "KRUD"
}

type ProjectCruds struct {
	ID           int             `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Form         int             `gorm:"column:form" json:"form"`
	Grid         int             `gorm:"column:grid" json:"grid"`
	ProjectsID   int             `gorm:"column:projects_id" json:"projects_id"`
	Template     string          `gorm:"column:template" json:"template"`
	Title        string          `gorm:"column:title" json:"title"`
	MainTabTitle string          `gorm:"column:main_tab_title" json:"main_tab_title"`
	CreatedAt    *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (p *ProjectCruds) TableName() string {
	return "project_cruds"
}
