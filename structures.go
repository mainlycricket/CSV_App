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
	minArrLen     int64           // 0 indicates unset
	maxArrLen     int64           // 0 indicates unset
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
	InsertAuth        AuthInfo   `json:"insertAuth"`
	ReadAllAuth       AuthInfo   `json:"readAllAuth"`
	ReadByPkAuth      AuthInfo   `json:"readByPkAuth"`
	UpdateAuth        AuthInfo   `json:"updateAuth"`
	DeleteAuth        AuthInfo   `json:"deleteAuth"`
	ReadAllConfig     ReadConfig `json:"readAllConfig"`
	ReadByPkConfig    ReadConfig `json:"readByPkConfig"`
	DefaultPagination int        `json:"defaultPagination"`
}

type AuthInfo struct {
	UserField       string              `json:"userField"`
	OrgFields       map[string]string   `json:"orgFields"`
	BasicAuth       bool                `json:"basicAuth"`
	AllowedRoles    []string            `json:"allowedRoles"`
	Privileges      map[string][]string `json:"privileges"`
	ProtectedFields ProtectedFieldsInfo `json:"protectedFields"`
}

type ReadConfig struct {
	Columns        []string            `json:"columns"`
	ForeignColumns map[string][]string `json:"foreignColumns"`
}

type ProtectedFieldsInfo map[string]map[string][]string
