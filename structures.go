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
	ForeignTable  string        `json:"foreignTable"` // if datatype is NOT array
	ForeignField  string        `json:"foreignField"` // if datatype is NOT array
	minIndividual interface{}
	maxIndividual interface{}
	minArrLen     int // 0 indicates unset
	maxArrLen     int // 0 indicates unset
	enumMap       map[any]bool
	values        map[string]bool // to check unique values
	lookup        map[string]int  // for foreign look up
}

type AppCongif struct {
	SchemaPath    string                     `json:"schemaPath"`
	AuthTable     string                     `json:"authTable"`
	UserField     string                     `json:"userField"`
	PasswordField string                     `json:"passwordField"`
	OrgField      string                     `json:"orgField"`
	TablesConfig  map[string]TableAuthConfig `json:"tablesConfig"`
}

type TableAuthConfig struct {
	UserField string     `json:"userField"`
	OrgField  string     `json:"orgField"`
	ReadAuth  AuthConfig `json:"readAuth"`  // GET
	WriteAuth AuthConfig `json:"writeAuth"` // INSERT, UPDATE, DELETE
}

type AuthConfig struct {
	BasicAuth    bool     `json:"basicAuth"`
	AllowedRoles []string `json:"allowedRoles"`
	SpecialRoles []string `json:"specialRoles"`
}
