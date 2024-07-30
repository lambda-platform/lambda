package models

import "time"

type ProjectSchemas struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string     `gorm:"column:name" json:"name"`
	ProjectsID int        `gorm:"column:projects_id" json:"projects_id"`
	Schema     string     `gorm:"column:schema" json:"schema"`
	Type       string     `gorm:"column:type" json:"type"`
	Actions    string     `gorm:"column:actions" json:"actions"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (p *ProjectSchemas) TableName() string {
	return "project_schemas"
}

type ProjectCruds struct {
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	Form         int        `gorm:"column:form" json:"form"`
	Grid         int        `gorm:"column:grid" json:"grid"`
	ID           int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProjectsID   int        `gorm:"column:projects_id" json:"projects_id"`
	Template     string     `gorm:"column:template" json:"template"`
	Title        string     `gorm:"column:title" json:"title"`
	MainTabTitle string     `gorm:"column:main_tab_title" json:"main_tab_title"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (p *ProjectCruds) TableName() string {
	return "project_cruds"
}

type SubProjectCruds struct {
	CrudID      int    `gorm:"column:crud_id" json:"crud_id"`
	Description string `gorm:"column:description" json:"description"`
	ID          int    `gorm:"column:id" json:"id"`
	Title       string `gorm:"column:title" json:"title"`
}

func (s *SubProjectCruds) TableName() string {
	return "sub_project_cruds"
}

type SubCrudSection struct {
	ID            int    `gorm:"column:id" json:"id"`
	ParintID      int    `gorm:"column:parint_id" json:"parint_id"`
	SectionCode   int    `gorm:"column:section_code" json:"section_code"`
	SesctionTitle string `gorm:"column:sesction_title" json:"sesction_title"`
}

func (s *SubCrudSection) TableName() string {
	return "sub_crud_section"
}

type SubCrud struct {
	ConnectionField string `gorm:"column:connection_field" json:"connection_field"`
	Form            int    `gorm:"column:form" json:"form"`
	FormField       string `gorm:"column:form_field" json:"form_field"`
	Grid            int    `gorm:"column:grid" json:"grid"`
	GridField       string `gorm:"column:grid_field" json:"grid_field"`
	ID              int    `gorm:"column:id" json:"id"`
	MicroserviceID  int    `gorm:"column:microservice_id" json:"microservice_id"`
	ParentID        int    `gorm:"column:parent_id" json:"parent_id"`
	SectionCode     int    `gorm:"column:section_code" json:"section_code"`
	Title           string `gorm:"column:title" json:"title"`
}

func (s *SubCrud) TableName() string {
	return "sub_crud"
}
