package models

type KrudTemplate struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TemplateName string `gorm:"column:template_name" json:"template_name"`
}

// TableName sets the insert table name for this struct type
func (v *KrudTemplate) TableName() string {
	return "krud_templates"
}

type KrudTemplateOracle struct {
	ID           int    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	TemplateName string `gorm:"column:TEMPLATE_NAME" json:"template_name"`
}

// TableName sets the insert table name for this struct type
func (v *KrudTemplateOracle) TableName() string {
	return "KRUD_TEMPLATES"
}
