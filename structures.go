package main

type DB struct {
	DB_Name  string           `json:"dbName"`
	BasePath string           `json:"basePath"`
	Tables   map[string]Table `json:"tables"` // key: tableName
}

type Table struct {
	FileName   string            `json:"fileName"`
	PrimaryKey string            `json:"primaryKey"`
	Columns    map[string]Column `json:"columns"` // key: columnName
}

type Column struct {
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
}
