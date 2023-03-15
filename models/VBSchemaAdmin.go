package models

import "time"

type VBSchemaAdmin struct {
	VbId      uint64     `gorm:"column:vb_id;primaryKey;autoIncrement" json:"vb_id"`
	Name      string     `gorm:"column:name" json:"name"`
	Schema    string     `gorm:"column:schema;type:TEXT" json:"schema"`
	Type      string     `gorm:"column:type" json:"type"`
	Id        string     `gorm:"column:id" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *VBSchemaAdmin) TableName() string {
	return "vb_schemas_admin"
}

type VBSchemaAdminOracle struct {
	VbID      int        `gorm:"column:VB_ID;primaryKey;autoIncrement" json:"vb_id"`
	ID        string     `gorm:"column:ID" json:"id"`
	Name      string     `gorm:"column:NAME" json:"name"`
	Schema    string     `gorm:"column:SCHEMA;type:LONG" json:"schema"`
	Type      string     `gorm:"column:TYPE" json:"type"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *VBSchemaAdminOracle) TableName() string {
	return "VB_SCHEMAS_ADMIN"
}

type VBSchemaAdminMSSQL struct {
	VbId      uint64     `gorm:"column:vb_id;primaryKey;autoIncrement" json:"vb_id"`
	Name      string     `gorm:"column:name" json:"name"`
	Schema    string     `gorm:"column:schema;type:NTEXT" json:"schema"`
	Type      string     `gorm:"column:type" json:"type"`
	Id        string     `gorm:"column:id" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *VBSchemaAdminMSSQL) TableName() string {
	return "vb_schemas_admin"
}
