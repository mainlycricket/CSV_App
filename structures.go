package main

type DB struct {
	DB_Name  string           `json:"dbName"`
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
