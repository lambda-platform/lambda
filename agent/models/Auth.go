package models

type LoginData struct {
	Token  string      `json:"token"`
	Path   string      `json:"path"`
	OAuth  bool        `json:"oauth"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type Unauthorized struct {
	Error  string `json:"error"`
	Status bool   `json:"status"`
}
