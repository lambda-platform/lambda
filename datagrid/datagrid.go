package datagrid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/models"
	"gorm.io/gorm"
)

type Datagrid struct {
	Name               string
	Identity           string
	DataTable          string
	MainTable          string
	DataModel          interface{}
	Data               interface{}
	MainModel          interface{}
	Columns            []Column
	ColumnList         []string
	Filters            map[string]string
	Relations          []models.GridRelation
	Condition          string
	Aggergation        string
	Triggers           map[string]interface{}
	BeforeFetch        func(interface{}, Datagrid, *gorm.DB, *fiber.Ctx) (interface{}, *gorm.DB, bool, bool)
	AfterFetch         func(interface{}, Datagrid, *gorm.DB, *fiber.Ctx) (interface{}, *gorm.DB, bool, bool)
	BeforeDelete       func(interface{}, Datagrid, *gorm.DB, *fiber.Ctx) (interface{}, *gorm.DB, bool, bool)
	AfterDelete        func(interface{}, Datagrid, *gorm.DB, *fiber.Ctx) (interface{}, *gorm.DB, bool, bool)
	BeforePrint        func(interface{}, Datagrid, *gorm.DB, *fiber.Ctx) (interface{}, *gorm.DB, bool, bool)
	TriggerNameSpace   string
	FillVirtualColumns func(interface{}) interface{}
}

type Column struct {
	Model string `json:"model"`
	Label string `json:"label"`
}

type RowUpdateData struct {
	Ids   []interface{} `json:"ids"`
	Value interface{}   `json:"value"`
	Model string        `json:"model"`
}

type CustomHeader struct {
	Render    bool `json:"render"`
	Preview   bool `json:"preview"`
	Structure []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Children []struct {
			ID      string      `json:"id"`
			Type    string      `json:"type"`
			Colspan string      `json:"colspan"`
			Rowspan string      `json:"rowspan"`
			Label   string      `json:"label"`
			Rotate  int         `json:"rotate"`
			Width   string      `json:"width"`
			Height  string      `json:"height"`
			Model   interface{} `json:"model,omitempty"`
		} `json:"children"`
	} `json:"structure"`
}
