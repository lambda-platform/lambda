package models

import "time"

type VBSchema struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Schema    string     `gorm:"column:schema;type:TEXT" json:"schema"` //type:LONGTEXT
	Type      string     `gorm:"column:type" json:"type"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (v *VBSchema) TableName() string {
	return "vb_schemas"
}

type VBSchemaMSSQL struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Schema    string     `gorm:"column:schema;type:NTEXT" json:"schema"` //type:LONGTEXT
	Type      string     `gorm:"column:type" json:"type"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *VBSchemaMSSQL) TableName() string {
	return "vb_schemas"
}

type VBSchemaOracle struct {
	ID        int        `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:NAME" json:"name"`
	Schema    string     `gorm:"column:SCHEMA;type:LONG" json:"schema"`
	Type      string     `gorm:"column:TYPE" json:"type"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

func (v *VBSchemaOracle) TableName() string {
	return "VB_SCHEMAS"
}
