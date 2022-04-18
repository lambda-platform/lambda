package models

type Permissions struct {
	DefaultMenu string      `json:"default_menu"`
	Extra       interface{} `json:"extra"`
	MenuID      int         `json:"menu_id"`
	Permissions interface{} `json:"permissions"`
}
