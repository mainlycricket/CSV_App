package main

type DB struct {
	BasePath string           `json:"basePath"`
	Tables   map[string]Table `json:"tables"` // key: tableName
}

type Table struct {
	TableName  string            `json:"tableName"`
	FileName   string            `json:"fileName"`
	PrimaryKey string            `json:"primaryKey"`
	Columns    map[string]Column `json:"columns"` // key: columnName
}

type Column struct {
	ColumnName    string        `json:"columnName"`
	DataType      string        `json:"dataType"`
	NotNull       bool          `json:"notNull"`
	Unique        bool          `json:"unique"`
	Hash          bool          `json:"hash"`
	Min           string        `json:"min"`
	Max           string        `json:"max"`
	Enums         []interface{} `json:"enums"`
	Default       interface{}   `json:"default"`
	ForeignTable  string        `json:"foreignTable"`
	ForeignField  string        `json:"foreignField"`
	OnUpdate      string        `json:"onUpdate"`
	OnDelete      string        `json:"onDelete"`
	minIndividual interface{}
	maxIndividual interface{}
	minArrLen     int             // 0 indicates unset
	maxArrLen     int             // 0 indicates unset
	values        map[string]bool // to check unique values
	lookup        map[string]int  // for foreign look up
}

type AppCongif struct {
	SchemaPath string                 `json:"schemaPath"`
	AuthTable  string                 `json:"authTable"`
	OrgFields  []string               `json:"orgFields"`
	Tables     map[string]TableConfig `json:"tables"`
}

type TableConfig struct {
	UserField         string            `json:"userField"`
	OrgFields         map[string]string `json:"orgFields"`
	ReadAuth          AuthInfo          `json:"readAuth"`  // GET
	WriteAuth         AuthInfo          `json:"writeAuth"` // INSERT, UPDATE, DELETE
	ReadAllConfig     ReadConfig        `json:"readAllConfig"`
	ReadByPkConfig    ReadConfig        `json:"readByPkConfig"`
	DefaultPagination int               `json:"defaultPagination"`
}

type AuthInfo struct {
	BasicAuth    bool                `json:"basicAuth"`
	AllowedRoles []string            `json:"allowedRoles"`
	Priviliges   map[string][]string `json:"priviliges"`
}

type ReadConfig struct {
	Columns        []string            `json:"columns"`
	ForeignColumns map[string][]string `json:"foreignColumns"`
}
