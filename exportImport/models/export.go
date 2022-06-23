package models

import (
	"time"
)

type LambdaExportData struct {
	Kruds []Krud `json:"kruds"`
}

type Schemas struct {
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	ID        int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Schema    string     `gorm:"column:schema" json:"schema"`
	Type      string     `gorm:"column:type" json:"type"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Krud struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Title      string     `gorm:"column:title" json:"title"`
	Template   string     `gorm:"column:template" json:"template"`
	Grid       int        `gorm:"column:grid" json:"grid"`
	Form       int        `gorm:"column:form" json:"form"`
	Actions    string     `gorm:"column:actions" json:"actions"`
	FormSchema Schemas    `json:"form_schema"`
	GridSchema Schemas    `json:"grid_schema"`
}
