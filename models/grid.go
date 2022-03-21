package models


type GridRelation struct {
	Column          string `json:"column"`
	ConnectionField string `json:"connection_field"`
	MicroserviceID  int    `json:"microservice_id"`
	Table           string `json:"table"`
	Key             string `json:"key"`
	Fields          string `json:"fields"`
	Self            bool   `json:"self"`
	Filter          string `json:"filter"`
}
