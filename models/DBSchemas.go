package models

type Formula struct {
	Targets  []Target `json:"targets"`
	Template string   `json:"template"`
	Form     string   `json:"form"`
	Model    string   `json:"model"`
}
type Target struct {
	Field string `json:"field"`
	Prop  string `json:"prop"`
}

type FormItem struct {
	Model       string      `json:"model"`
	Title       string      `json:"title"`
	DbType      string      `json:"dbType"`
	Table       string      `json:"table"`
	Key         string      `json:"key"`
	Extra       string      `json:"extra"`
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Identity    string      `json:"identity"`
	Label       string      `json:"label"`
	PlaceHolder string      `json:"placeHolder"`
	Param       string      `json:"param"`
	Hidden      bool        `json:"hidden"`
	Disabled    bool        `json:"disabled"`
	Default     interface{} `json:"default"`
	Prefix      string      `json:"prefix"`
	Ifshowhide  string      `json:"ifshowhide"`
	Rules       []struct {
		Type string `json:"type"`
		Msg  string `json:"msg"`
	} `json:"rules"`
	HasTranslation bool   `json:"hasTranslation"`
	HasUserID      bool   `json:"hasUserId"`
	HasEquation    bool   `json:"hasEquation"`
	Equations      string `json:"equations"`
	IsGridSearch   bool   `json:"isGridSearch"`
	GridSearch     struct {
		Grid     interface{} `json:"grid"`
		Key      interface{} `json:"key"`
		Labels   interface{} `json:"labels"`
		Multiple bool        `json:"multiple"`
	} `json:"gridSearch"`
	IsFkey   bool   `json:"isFkey"`
	InfoUrl  string `json:"info_url"`
	Relation struct {
		MicroserviceID     int           `json:"microservice_id"`
		Table              string        `json:"table"`
		Key                interface{}   `json:"key"`
		Fields             []interface{} `json:"fields"`
		FilterWithUser     []interface{} `json:"filterWithUser"`
		SortField          interface{}   `json:"sortField"`
		SortOrder          string        `json:"sortOrder"`
		Multiple           bool          `json:"multiple"`
		Filter             string        `json:"filter"`
		ParentFieldOfForm  string        `json:"parentFieldOfForm"`
		ParentFieldOfTable string        `json:"parentFieldOfTable"`
	} `json:"relation"`
	Span struct {
		Xs int `json:"xs"`
		Sm int `json:"sm"`
		Md int `json:"md"`
		Lg int `json:"lg"`
	} `json:"span"`
	Trigger        string `json:"trigger"`
	TriggerTimeout int    `json:"triggerTimeout"`
	File           struct {
		IsMultiple bool   `json:"isMultiple"`
		Count      int    `json:"count"`
		MaxSize    int    `json:"maxSize"`
		Type       string `json:"type"`
	} `json:"file"`
	Options          []interface{} `json:"options"`
	PasswordOption   interface{}   `json:"passwordOption"`
	GeographicOption interface{}   `json:"GeographicOption"`
	EditorType       interface{}   `json:"editorType"`
	SchemaID         string        `json:"schemaID"`

	//subForm data
	Name                           string                 `json:"name"`
	SubType                        string                 `json:"subtype"`
	Parent                         string                 `json:"parent"`
	FormId                         uint64                 `json:"formId"`
	MapID                          string                 `json:"mapID"`
	SelectedType                   string                 `json:"selectedType"`
	FormType                       string                 `json:"formType"`
	MinHeight                      string                 `json:"min_height"`
	DisableDelete                  bool                   `json:"disableDelete"`
	DisableEdit                    bool                   `json:"disableEdit"`
	DisableCreate                  bool                   `json:"disableCreate"`
	ShowRowNumber                  bool                   `json:"showRowNumber"`
	UseTableType                   bool                   `json:"useTableType"`
	CheckEmpty                     bool                   `json:"checkEmpty"`
	AddFromGrid                    bool                   `json:"addFromGrid"`
	NoFormat                       bool                   `json:"no_format"`
	Precision                      interface{}            `json:"precision"`
	TableTypeColumn                string                 `json:"tableTypeColumn"`
	TableTypeValue                 string                 `json:"tableTypeValue"`
	EmptyErrorMsg                  string                 `json:"EmptyErrorMsg"`
	SourceGridModalTitle           string                 `json:"sourceGridModalTitle"`
	SourceGridParentBasedCondition string                 `json:"sourceGridParentBasedCondition"`
	GSOption                       map[string]interface{} `json:"GSOption"`
	SourceGridTitle                string                 `json:"sourceGridTitle"`
	SourceGridDescription          string                 `json:"sourceGridDescription"`
	SourceGridUserCondition        string                 `json:"sourceGridUserCondition"`
	SourceUniqueField              string                 `json:"sourceUniqueField"`
	SourceMicroserviceID           interface{}            `json:"sourceMicroserviceID"`
	SourceGridID                   interface{}            `json:"sourceGridID"`
	SourceGridTargetColumns        []interface{}          `json:"sourceGridTargetColumns"`
	Schema                         []FormItem             `json:"schema"`
	Nullable                       string                 `json:"nullable"`
	Scale                          string                 `json:"scale"`
	DefaultValue                   string                 `json:"default_value"`
	TableSchema                    string                 `json:"table_schema"`
	Warn                           string                 `json:"warn"`
	Data                           any                    `json:"data"`
	Rule                           any                    `json:"rule"`
	TrKey                          any                    `json:"trKey"`
}

type SCHEMA struct {
	FormType      string        `json:"formType"`
	FormSubName   string        `json:"formSubName"`
	Model         string        `json:"model"`
	Identity      string        `json:"identity"`
	Timestamp     bool          `json:"timestamp"`
	LabelPosition string        `json:"labelPosition"`
	LabelWidth    interface{}   `json:"labelWidth"`
	ExtraButtons  []interface{} `json:"extraButtons"`
	Width         string        `json:"width"`
	Padding       int           `json:"padding"`
	Schema        []FormItem    `json:"schema"`
	UI            interface{}   `json:"ui"`
	Formula       []Formula     `json:"formula"`
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
	SortField                string `json:"sortField"`
	SordOrder                string `json:"sordOrder"`
	Use2ColumnLayout         bool   `json:"use2ColumnLayout"`
	WithBackButton           bool   `json:"withBackButton"`
	SaveBtnText              string `json:"save_btn_text"`
	IsWarnText               bool   `json:"isWarnText"`
	WarnText                 string `json:"warnText"`
	FormValidationCustomText string `json:"formValidationCustomText"`
	DisableReset             bool   `json:"disableReset"`
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
		VirtualColumn bool   `json:"virtualColumn"`
		Model         string `json:"model"`
		Title         string `json:"title"`
		DbType        string `json:"dbType"`
		Table         string `json:"table"`
		Key           string `json:"key"`
		Extra         string `json:"extra"`
		Label         string `json:"label"`
		GridType      string `json:"gridType"`
		Width         int    `json:"width"`
		Hide          bool   `json:"hide"`
		Sortable      bool   `json:"sortable"`
		Printable     bool   `json:"printable"`
		Pinned        bool   `json:"pinned"`
		PinPosition   string `json:"pinPosition"`
		Link          string `json:"link"`
		LinkTarget    string `json:"linkTarget"`
		Relation      struct {
			Column          string `json:"column"`
			ConnectionField string `json:"connection_field"`
			MicroserviceID  int    `json:"microservice_id"`
			Table           string `json:"table"`
			Key             string `json:"key"`
			Fields          string `json:"fields"`
			Self            bool   `json:"self"`
			Filter          string `json:"filter"`
		} `json:"relation"`
		Filterable bool `json:"filterable"`
		Filter     struct {
			Type             string `json:"type"`
			Param            any    `json:"param"`
			ParamCompareType string `json:"paramCompareType"`
			Default          any    `json:"default"`
			Relation         struct {
				Table     string `json:"table"`
				Key       any    `json:"key"`
				Fields    []any  `json:"fields"`
				SortField any    `json:"sortField"`
				SortOrder string `json:"sortOrder"`
			} `json:"relation"`
		} `json:"filter"`
		Editable struct {
			Status       bool   `json:"status"`
			Type         string `json:"type"`
			ShouldUpdate bool   `json:"shouldUpdate"`
			ShouldPost   bool   `json:"shouldPost"`
		} `json:"editable"`
		Searchable           bool  `json:"searchable"`
		HasTranslation       bool  `json:"hasTranslation"`
		Options              []any `json:"options"`
		TrKey                any   `json:"trKey"`
		CanExcelImport       bool  `json:"canExcelImport"`
		ExcelImportFieldName any   `json:"excelImportFieldName"`
		ExcelImportRelation  struct {
			Table              any    `json:"table"`
			Key                any    `json:"key"`
			Fields             []any  `json:"fields"`
			SortField          any    `json:"sortField"`
			SortOrder          string `json:"sortOrder"`
			Multiple           bool   `json:"multiple"`
			Filter             string `json:"filter"`
			ParentFieldOfForm  string `json:"parentFieldOfForm"`
			ParentFieldOfTable string `json:"parentFieldOfTable"`
		} `json:"excelImportRelation"`
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
	Theme                      string         `json:"theme"`
	FullText                   bool           `json:"fullText"`
	EditableAction             interface{}    `json:"editableAction"`
	EditFullRow                bool           `json:"editFullRow"`
	EditableShouldSubmit       bool           `json:"editableShouldSubmit"`
	SingleClickEdit            bool           `json:"singleClickEdit"`
	FlashChanges               bool           `json:"flashChanges"`
	ColMenu                    bool           `json:"colMenu"`
	ColFilterButton            bool           `json:"colFilterButton"`
	ShowGrid                   bool           `json:"showGrid"`
	SordOrder                  string         `json:"sordOrder"`
	MainTable                  string         `json:"mainTable"`
	IsPivot                    bool           `json:"isPivot"`
	IsPrint                    bool           `json:"isPrint"`
	PrintSize                  string         `json:"printSize"`
	IsExcel                    bool           `json:"isExcel"`
	IsExcelUpload              bool           `json:"isExcelUpload"`
	ExcelUploadCustomNamespace string         `json:"excelUploadCustomNamespace"`
	ExcelUploadCustomTrigger   string         `json:"excelUploadCustomTrigger"`
	GridTheme                  string         `json:"gridTheme"`
	IsRefresh                  bool           `json:"isRefresh"`
	IsNumbered                 bool           `json:"isNumbered"`
	Microservices              []Microservice `json:"microservices"`
	ExcelUploadSample          any            `json:"excelUploadSample"`
	ExcelImportRowtoStart      any            `json:"excelImportRowtoStart"`
	ExcelUploadCustomURL       any            `json:"excelUploadCustomUrl"`
	IsGlobalSearch             bool           `json:"isGlobalSearch"`
	ExcelImportRelation        struct {
		Table              any    `json:"table"`
		Key                any    `json:"key"`
		Fields             []any  `json:"fields"`
		SortField          any    `json:"sortField"`
		SortOrder          string `json:"sortOrder"`
		Multiple           bool   `json:"multiple"`
		Filter             string `json:"filter"`
		ParentFieldOfForm  string `json:"parentFieldOfForm"`
		ParentFieldOfTable string `json:"parentFieldOfTable"`
	} `json:"excelImportRelation"`
	SaveFilter      bool `json:"saveFilter"`
	AutoSelect      bool `json:"autoSelect"`
	AutoSelectModel any  `json:"autoSelectModel"`
}
