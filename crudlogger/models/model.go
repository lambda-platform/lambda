package models

import "time"

type CrudLog struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:user_id;" json:"user_id"`
	Ip        string    `gorm:"column:ip;" json:"ip"`
	UserAgent string    `gorm:"column:user_agent;" json:"user_agent"`
	Action    string    `gorm:"column:action;" json:"action"`
	SchemaId  int64     `gorm:"column:schema_id;" json:"schema_id"`
	RowId     string    `gorm:"column:row_id;" json:"row_id"`
	Input     string    `gorm:"column:input;type:TEXT" json:"input"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (v *CrudLog) TableName() string {
	return "crud_log"
}

type CrudLogOracle struct {
	ID        int       `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:USER_ID;" json:"user_id"`
	Ip        string    `gorm:"column:IP;" json:"ip"`
	UserAgent string    `gorm:"column:USER_AGENT;" json:"user_agent"`
	Action    string    `gorm:"column:ACTION;" json:"action"`
	SchemaId  int64     `gorm:"column:SCHEMA_ID;" json:"schema_id"`
	RowId     string    `gorm:"column:ROW_ID;" json:"row_id"`
	Input     string    `gorm:"column:INPUT;type:TEXT" json:"input"`
	CreatedAt time.Time `gorm:"column:CREATED_AT" json:"created_at"`
}

func (v *CrudLogOracle) TableName() string {
	return "CRUD_LOG"
}

type CrudLogMSSQL struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:user_id;" json:"user_id"`
	Ip        string    `gorm:"column:ip;" json:"ip"`
	UserAgent string    `gorm:"column:user_agent;" json:"user_agent"`
	Action    string    `gorm:"column:action;" json:"action"`
	SchemaId  int64     `gorm:"column:schema_id;" json:"schema_id"`
	RowId     string    `gorm:"column:row_id;" json:"row_id"`
	Input     string    `gorm:"column:input;type:NTEXT" json:"input"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (v *CrudLogMSSQL) TableName() string {
	return "crud_log"
}

type CrudResponse struct {
	ID interface{} `gorm:"column:id;" json:"id"`
}
type HistoryRequest struct {
	SchemaID int64       `json:"schemaId"`
	RowID    interface{} `json:"rowId"`
}
type CrudLogFullOracle struct {
	Action    *string    `gorm:"column:ACTION" json:"action"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	FirstName *string    `gorm:"column:FIRST_NAME" json:"first_name"`
	ID        int        `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	Input     *string    `gorm:"column:INPUT" json:"input"`
	IP        *string    `gorm:"column:IP" json:"ip"`
	LastName  *string    `gorm:"column:LAST_NAME" json:"last_name"`
	Name      *string    `gorm:"column:NAME" json:"name"`
	RowID     *string    `gorm:"column:ROW_ID" json:"row_id"`
	SchemaID  *int       `gorm:"column:SCHEMA_ID" json:"schema_id"`
	UserAgent *string    `gorm:"column:USER_AGENT" json:"user_agent"`
	UserID    *int       `gorm:"column:USER_ID" json:"user_id"`
}

func (c *CrudLogFullOracle) TableName() string {
	return "DS_CRUD_LOG"
}

type CrudLogFull struct {
	Action    *string    `gorm:"column:action" json:"action"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	FirstName *string    `gorm:"column:first_name" json:"first_name"`
	ID        int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Input     *string    `gorm:"column:input" json:"input"`
	IP        *string    `gorm:"column:ip" json:"ip"`
	LastName  *string    `gorm:"column:last_name" json:"last_name"`
	Name      *string    `gorm:"column:name" json:"name"`
	RowID     *string    `gorm:"column:row_id" json:"row_id"`
	SchemaID  *int       `gorm:"column:schema_id" json:"schema_id"`
	UserAgent *string    `gorm:"column:user_agent" json:"user_agent"`
	UserID    *int       `gorm:"column:user_id" json:"user_id"`
}

func (c *CrudLogFull) TableName() string {
	return "ds_crud_log"
}
