package models

import "github.com/lambda-platform/dataform/models"


type SCHEMA struct {
	Model         string      `json:"model"`
	Identity      string      `json:"identity"`
	Timestamp     bool        `json:"timestamp"`
	LabelPosition string      `json:"labelPosition"`
	LabelWidth    interface{} `json:"labelWidth"`
	Width         string      `json:"width"`
	Padding       int         `json:"padding"`
	Schema        []models.FormItem  `json:"schema"`
	UI            interface{} `json:"ui"`
	Formula       []models.Formula   `json:"formula"`
	Triggers      struct {
		Namespace string `json:"namespace"`
		Insert    struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"insert"`
		Update struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"update"`
	} `json:"triggers"`
	SortField string `json:"sortField"`
	SordOrder string `json:"sordOrder"`
}

type SCHEMAGRID struct {
	Model          string   `json:"model"`
	IsView         bool     `json:"isView"`
	Identity       string   `json:"identity"`
	Actions        []string `json:"actions"`
	ActionPosition int      `json:"actionPosition"`
	IsContextMenu  bool     `json:"isContextMenu"`
	StaticWidth    bool     `json:"staticWidth"`
	FullWidth      bool     `json:"fullWidth"`
	HasCheckbox    bool     `json:"hasCheckbox"`
	IsClient       bool     `json:"isClient"`
	Width          int      `json:"width"`
	Sort           string   `json:"sort"`
	SortOrder      string   `json:"sortOrder"`
	SoftDelete     bool     `json:"softDelete"`
	Paging         int      `json:"paging"`
	Template       int      `json:"template"`
	Schema         []struct {
		Model       string `json:"model"`
		Title       string `json:"title"`
		DbType      string `json:"dbType"`
		Table       string `json:"table"`
		Key         string `json:"key"`
		Extra       string `json:"extra"`
		Label       string `json:"label"`
		GridType    string `json:"gridType"`
		Width       int    `json:"width"`
		Hide        bool   `json:"hide"`
		Sortable    bool   `json:"sortable"`
		Printable   bool   `json:"printable"`
		Pinned      bool   `json:"pinned"`
		PinPosition string `json:"pinPosition"`
		Link        string `json:"link"`
		LinkTarget  string `json:"linkTarget"`
		Relation    struct {
			Table              interface{}   `json:"table"`
			Key                interface{}   `json:"key"`
			Fields             []interface{} `json:"fields"`
			SortField          interface{}   `json:"sortField"`
			SortOrder          string        `json:"sortOrder"`
			Multiple           bool          `json:"multiple"`
			Filter             string        `json:"filter"`
			ParentFieldOfForm  string        `json:"parentFieldOfForm"`
			ParentFieldOfTable string        `json:"parentFieldOfTable"`
		} `json:"relation"`
		Filterable bool `json:"filterable"`
		Filter     struct {
			Type             string      `json:"type"`
			Param            interface{} `json:"param"`
			ParamCompareType string      `json:"paramCompareType"`
			Default          interface{} `json:"default"`
			Relation         struct {
				Table     interface{}   `json:"table"`
				Key       interface{}   `json:"key"`
				Fields    []interface{} `json:"fields"`
				SortField interface{}   `json:"sortField"`
				SortOrder string        `json:"sortOrder"`
			} `json:"relation"`
		} `json:"filter"`
		Editable struct {
			Status       bool   `json:"status"`
			Type         string `json:"type"`
			ShouldUpdate bool   `json:"shouldUpdate"`
			ShouldPost   bool   `json:"shouldPost"`
		} `json:"editable"`
		Searchable     bool          `json:"searchable"`
		HasTranslation bool          `json:"hasTranslation"`
		Options        []interface{} `json:"options"`
	} `json:"schema"`
	Filter                    []interface{}       `json:"filter"`
	Formula                   []interface{}       `json:"formula"`
	Condition                 string              `json:"condition"`
	ColumnAggregations        []map[string]string `json:"columnAggregations"`
	ColumnAggregationsFormula []interface{}       `json:"columnAggregationsFormula"`
	Header                    struct {
		Render    bool          `json:"render"`
		Preview   bool          `json:"preview"`
		Structure []interface{} `json:"structure"`
	} `json:"header"`
	Triggers struct {
		Namespace    string `json:"namespace"`
		BeforeFetch  string `json:"beforeFetch"`
		AfterFetch   string `json:"afterFetch"`
		BeforeDelete string `json:"beforeDelete"`
		AfterDelete  string `json:"afterDelete"`
		BeforePrint  string `json:"beforePrint"`
	} `json:"triggers"`
	Theme                string      `json:"theme"`
	FullText             bool        `json:"fullText"`
	EditableAction       interface{} `json:"editableAction"`
	EditFullRow          bool        `json:"editFullRow"`
	EditableShouldSubmit bool        `json:"editableShouldSubmit"`
	SingleClickEdit      bool        `json:"singleClickEdit"`
	FlashChanges         bool        `json:"flashChanges"`
	ColMenu              bool        `json:"colMenu"`
	ColFilterButton      bool        `json:"colFilterButton"`
	ShowGrid             bool        `json:"showGrid"`
	SordOrder            string      `json:"sordOrder"`
	MainTable            string      `json:"mainTable"`
	IsPivot              bool        `json:"isPivot"`
	IsPrint              bool        `json:"isPrint"`
	PrintSize            string      `json:"printSize"`
	IsExcel              bool        `json:"isExcel"`
	IsRefresh            bool        `json:"isRefresh"`
	IsNumbered           bool        `json:"isNumbered"`
}
