package models


type Microservice struct {
	GRPCURL        string `json:"grpc_url"`
	ProductionURL string `json:"production_url"`
	ProjectID     int    `json:"project_id"`
}
