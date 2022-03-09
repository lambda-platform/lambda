package models

import "time"


type UserRoles struct {
	ID        uint64       `gorm:"column:id;primary_key" json:"id"`
	DisplayName      string    `gorm:"column:display_name" json:"display_name"`

}

func (v *UserRoles) TableName() string {
	return "roles"
}
type VBSchemaList struct {
	ID        uint64       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Type      string    `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
type VBSchema struct {
	ID        uint64       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Schema    string    `gorm:"column:schema;type:TEXT" json:"schema"` //type:LONGTEXT
	Type      string    `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (v *VBSchema) TableName() string {
	return "vb_schemas"
}
type ProjectVBSchema struct {
	ID        uint64       `gorm:"column:id;primary_key" json:"id"`
	ProjectID        int       `gorm:"column:projects_id" json:"projects_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Schema    string    `gorm:"column:schema;type:TEXT" json:"schema"` //type:LONGTEXT
	Type      string    `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (v *ProjectVBSchema) TableName() string {
	return "project_schemas"
}

func (v *VBSchemaList) TableName() string {
	return "vb_schemas"
}

type VBSchemaMSSQL struct {
	ID        uint64       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Schema    string    `gorm:"column:schema;type:NTEXT" json:"schema"` //type:LONGTEXT
	Type      string    `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

//  TableName sets the insert table name for this struct type
func (v *VBSchemaMSSQL) TableName() string {
	return "vb_schemas"
}

type PKColumn struct {
	ColumnName      string    `gorm:"column:COLUMN_NAME" json:"COLUMN_NAME"`
}

type MSTableMata struct {
	ColumnName      string    `gorm:"column:COLUMN_NAME" json:"COLUMN_NAME"`
	DataType      string    `gorm:"column:DATA_TYPE" json:"DATA_TYPE"`
}

type VBSchemaAdmin struct {
	VbId        uint64       `gorm:"column:vb_id;primary_key" json:"vb_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Schema    string    `gorm:"column:schema;type:TEXT" json:"schema"`
	Type      string    `gorm:"column:type" json:"type"`
	Id      string    `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

//  TableName sets the insert table name for this struct type
func (v *VBSchemaAdmin) TableName() string {
	return "vb_schemas_admin"
}

type VBSchemaAdminMSSQL struct {
	VbId        uint64       `gorm:"column:vb_id;primary_key" json:"vb_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Schema    string    `gorm:"column:schema;type:NTEXT" json:"schema"`
	Type      string    `gorm:"column:type" json:"type"`
	Id      string    `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

//  TableName sets the insert table name for this struct type
func (v *VBSchemaAdminMSSQL) TableName() string {
	return "vb_schemas_admin"
}
type TableMeta struct {
	Model  string `json:"model"`
	Title  string `json:"title"`
	DbType string `json:"dbType"`
	Table  string `json:"table"`
	Key    string `json:"key"`
	Extra  string `json:"extra"`
	Nullable  string `json:"nullable"`
}
type DBSCHEMA struct {
	TableList []string               `json:"tableList"`
	ViewList  []string               `json:"viewList"`
	TableMeta map[string][]TableMeta `json:"tableMeta"`
	MicroserviceID int `json:"microservice_id"`
	Microservice string `json:"microservice"`
}