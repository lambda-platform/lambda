package models

type OauthClients struct {
	ID       int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClientID string `gorm:"column:client_id" json:"client_id"`
	Secret   string `gorm:"column:secret" json:"secret"`
	Domain   string `gorm:"column:domain" json:"domain"`
}

func (o *OauthClients) TableName() string {
	return "oauth_clients"
}

type OracleOauthClients struct {
	ID       int    `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	ClientID string `gorm:"column:CLIENT_ID" json:"client_id"`
	Secret   string `gorm:"column:SECRET" json:"secret"`
	Domain   string `gorm:"column:DOMAIN" json:"domain"`
}

func (o *OracleOauthClients) TableName() string {
	return "OAUTH_CLIENTS"
}
