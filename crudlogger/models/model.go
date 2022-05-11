package models

import "time"

type CrudLog struct {
	ID        int     `gorm:"column:id;primary_key" json:"id"`
	UserId        int64     `gorm:"column:user_id;" json:"user_id"`
	Ip        string     `gorm:"column:ip;" json:"ip"`
	UserAgent        string     `gorm:"column:user_agent;" json:"user_agent"`
	Action        string     `gorm:"column:action;" json:"action"`
	SchemaId        int64     `gorm:"column:schemaId;" json:"schemaId"`
	RowId        string     `gorm:"column:row_id;" json:"row_id"`
	Input        string     `gorm:"column:input;type:TEXT" json:"input"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
func (v *CrudLog) TableName() string {
	return "crud_log"
}
type CrudLogMSSQL struct {
	ID        int     `gorm:"column:id;primary_key" json:"id"`
	UserId        int64     `gorm:"column:user_id;" json:"user_id"`
	Ip        string     `gorm:"column:ip;" json:"ip"`
	UserAgent        string     `gorm:"column:user_agent;" json:"user_agent"`
	Action        string     `gorm:"column:action;" json:"action"`
	SchemaId        int64     `gorm:"column:schemaId;" json:"schemaId"`
	RowId        string     `gorm:"column:row_id;" json:"row_id"`
	Input        string     `gorm:"column:input;type:NTEXT" json:"input"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
func (v *CrudLogMSSQL) TableName() string {
	return "crud_log"
}

type CrudResponse struct {
	Data struct{
		ID        int     `gorm:"column:id;" json:"id"`
	} `json:"data"`
}


