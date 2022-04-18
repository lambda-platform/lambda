package models

import "time"

type Microservice struct {
	GRPCURL        string `json:"grpc_url"`
	ProductionURL string `json:"production_url"`
	ProjectID     int    `json:"project_id"`
}
type Projects struct {
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	DbSchemaPath    string     `gorm:"column:db_schema_path" json:"db_schema_path"`
	ID              int        `gorm:"column:id;primary_key" json:"id"`
	Name            string     `gorm:"column:name" json:"name"`
	ProjectType            string     `gorm:"column:project_type" json:"project_type"`
	OrganizationsID int        `gorm:"column:organizations_id" json:"organizations_id"`
	ProjectKey      string     `gorm:"column:project_key" json:"project_key"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
	UsersID         int        `gorm:"column:users_id" json:"users_id"`
}

func (p *Projects) TableName() string {
	return "projects"
}
type ProjectSettings struct {
	GRPCURL        string `gorm:"column:grpc_url" json:"grpc_url"`
	ID            int    `gorm:"column:id" json:"id"`
	ProductionURL string `gorm:"column:production_url" json:"production_url"`
	DevUrl string `gorm:"column:dev_url" json:"dev_url"`
	ProjectID     int    `gorm:"column:project_id" json:"project_id"`
}

func (p *ProjectSettings) TableName() string {
	return "project_settings"
}