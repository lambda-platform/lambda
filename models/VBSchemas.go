package models

import "time"

type UserRoles struct {
	ID          uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`
}

func (v *UserRoles) TableName() string {
	return "roles"
}

type VBSchemaListOracle struct {
	ID        uint64     `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:NAME" json:"name"`
	Type      string     `gorm:"column:TYPE" json:"type"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

func (v *VBSchemaListOracle) TableName() string {
	return "VB_SCHEMAS"
}

type VBSchemaList struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Type      string     `gorm:"column:type" json:"type"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

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

type ProjectVBSchema struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProjectID int        `gorm:"column:projects_id" json:"projects_id"`
	Name      string     `gorm:"column:name" json:"name"`
	Schema    string     `gorm:"column:schema;type:TEXT" json:"schema"` //type:LONGTEXT
	Type      string     `gorm:"column:type" json:"type"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (v *ProjectVBSchema) TableName() string {
	return "project_schemas"
}

func (v *VBSchemaList) TableName() string {
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

type PKColumn struct {
	ColumnName string `gorm:"column:COLUMN_NAME" json:"COLUMN_NAME"`
}

type MSTableMata struct {
	ColumnName string `gorm:"column:COLUMN_NAME" json:"COLUMN_NAME"`
	DataType   string `gorm:"column:DATA_TYPE" json:"DATA_TYPE"`
}

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

type TableMeta struct {
	Model    string `json:"model"`
	Title    string `json:"title"`
	DbType   string `json:"dbType"`
	Table    string `json:"table"`
	Key      string `json:"key"`
	Extra    string `json:"extra"`
	Nullable string `json:"nullable"`
}
type DBSCHEMA struct {
	TableList      []string               `json:"tableList"`
	ViewList       []string               `json:"viewList"`
	TableMeta      map[string][]TableMeta `json:"tableMeta"`
	MicroserviceID int                    `json:"microservice_id"`
	Microservice   string                 `json:"microservice"`
}
