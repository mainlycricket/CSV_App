package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_sanitize_db_label(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Trim",
			args: args{text: "  Table Name "},
			want: "Table_Name",
		},
		{
			name: "Multiple Symbols",
			args: args{text: "  Table(-)Name "},
			want: "Table_Name",
		},
		{
			name: "Start & End Symbols",
			args: args{text: "  !'Table Name'! "},
			want: "_Table_Name_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitize_db_label(tt.args.text); got != tt.want {
				t.Errorf("sanitize_db_label() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateValueByType(t *testing.T) {
	type args struct {
		value    any
		datatype string
	}

	tests := []struct {
		name         string
		args         args
		convertedVal any
		success      bool
	}{
		{
			name: "date",
			args: args{value: "2024-01-01", datatype: "date"},
			convertedVal: func() time.Time {
				dateRes, _ := time.Parse(datetimeFormats["date"], "2024-01-01")
				return dateRes
			}(),
			success: true,
		},
		{
			name:         "false boolean",
			args:         args{value: false, datatype: "boolean"},
			convertedVal: false,
			success:      true,
		},
		{
			name:         "text",
			args:         args{value: "text", datatype: "text"},
			convertedVal: "text",
			success:      true,
		},
		{
			name:         "integer",
			args:         args{value: 1, datatype: "integer"},
			convertedVal: 1,
			success:      true,
		},
		{
			name:         "float",
			args:         args{value: 1.6, datatype: "real"},
			convertedVal: 1.6,
			success:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validateValueByType(tt.args.value, tt.args.datatype)
			if !reflect.DeepEqual(got, tt.convertedVal) {
				t.Errorf("validateValueByType() got = %v, want %v", got, tt.convertedVal)
			}
			if got1 != tt.success {
				t.Errorf("validateValueByType() got1 = %v, want %v", got1, tt.success)
			}
		})
	}
}

func TestColumn_validateDefaultValue(t *testing.T) {
	type fields struct {
		ColumnName    string
		DataType      string
		NotNull       bool
		Unique        bool
		Min           string
		Max           string
		Enums         []interface{}
		Default       interface{}
		ForeignTable  string
		ForeignField  string
		minIndividual interface{}
		maxIndividual interface{}
		minArrLen     int
		maxArrLen     int
		enumMap       map[any]bool
		values        map[string]bool
		lookup        map[string]int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{DataType: "integer", Default: nil, minIndividual: 45},
			wantErr: false,
		},
		{
			name:    "int",
			fields:  fields{DataType: "integer", Default: 2},
			wantErr: false,
		},
		{
			name:    "int arr",
			fields:  fields{DataType: "integer[]", Default: []any{1, 2}},
			wantErr: false,
		},
		{
			name:    "invalid int arr",
			fields:  fields{DataType: "integer[]", Default: []any{3, 4}, minArrLen: 3},
			wantErr: true,
		},
		{
			name:    "text arr",
			fields:  fields{DataType: "text[]", Default: []any{"te", "dd"}},
			wantErr: false,
		},
		{
			name:    "invalid text",
			fields:  fields{DataType: "text", Default: "text", enumMap: map[any]bool{"other": true}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := &Column{
				ColumnName:    tt.fields.ColumnName,
				DataType:      tt.fields.DataType,
				NotNull:       tt.fields.NotNull,
				Unique:        tt.fields.Unique,
				Min:           tt.fields.Min,
				Max:           tt.fields.Max,
				Enums:         tt.fields.Enums,
				Default:       tt.fields.Default,
				ForeignTable:  tt.fields.ForeignTable,
				ForeignField:  tt.fields.ForeignField,
				minIndividual: tt.fields.minIndividual,
				maxIndividual: tt.fields.maxIndividual,
				minArrLen:     tt.fields.minArrLen,
				maxArrLen:     tt.fields.maxArrLen,
				enumMap:       tt.fields.enumMap,
				values:        tt.fields.values,
				lookup:        tt.fields.lookup,
			}
			if err := column.validateDefaultValue(); (err != nil) != tt.wantErr {
				t.Errorf("Column.validateDefaultValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestColumn_validateValArrLen(t *testing.T) {
	type fields struct {
		ColumnName    string
		DataType      string
		NotNull       bool
		Unique        bool
		Min           string
		Max           string
		Enums         []interface{}
		Default       interface{}
		ForeignTable  string
		ForeignField  string
		minIndividual interface{}
		maxIndividual interface{}
		minArrLen     int
		maxArrLen     int
		enumMap       map[any]bool
		values        map[string]bool
		lookup        map[string]int
	}
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []any
		wantErr bool
	}{
		{
			name:    "valid integer arr",
			fields:  fields{minArrLen: 2, maxArrLen: 5},
			args:    args{value: []any{1, 2, 3}},
			want:    []any{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "invalid string arr",
			fields:  fields{minArrLen: 2, maxArrLen: 5},
			args:    args{value: []any{"text"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid string",
			fields:  fields{minArrLen: 2, maxArrLen: 5},
			args:    args{value: "[1.2, 2.5, 3.4]"},
			want:    []any{1.2, 2.5, 3.4},
			wantErr: false,
		},
		{
			name:    "invalid string",
			fields:  fields{minArrLen: 2, maxArrLen: 5},
			args:    args{value: "not an array"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid string array",
			fields:  fields{minArrLen: 2, maxArrLen: 5},
			args:    args{value: "[1]"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := &Column{
				ColumnName:    tt.fields.ColumnName,
				DataType:      tt.fields.DataType,
				NotNull:       tt.fields.NotNull,
				Unique:        tt.fields.Unique,
				Min:           tt.fields.Min,
				Max:           tt.fields.Max,
				Enums:         tt.fields.Enums,
				Default:       tt.fields.Default,
				ForeignTable:  tt.fields.ForeignTable,
				ForeignField:  tt.fields.ForeignField,
				minIndividual: tt.fields.minIndividual,
				maxIndividual: tt.fields.maxIndividual,
				minArrLen:     tt.fields.minArrLen,
				maxArrLen:     tt.fields.maxArrLen,
				enumMap:       tt.fields.enumMap,
				values:        tt.fields.values,
				lookup:        tt.fields.lookup,
			}
			got, err := column.validateValArrLen(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.validateValArrLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Column.validateValArrLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareTypeValues(t *testing.T) {
	type args struct {
		a        any
		b        any
		datatype string
	}

	date1, _ := time.Parse(time.DateOnly, "2024-02-01")
	date2, _ := time.Parse(time.DateOnly, "2024-01-31")

	time1, _ := time.Parse(time.TimeOnly, "15:26:59")
	time2, _ := time.Parse(time.TimeOnly, "15:26:59")

	datetime1, _ := time.Parse(time.RFC3339, "2024-07-01T12:30:00+05:30")
	datetime2, _ := time.Parse(time.RFC3339, "2024-08-01T12:30:00+05:30")

	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name:  "a > b integer",
			args:  args{a: 6, b: 5, datatype: "integer"},
			want:  1,
			want1: true,
		},
		{
			name:  "a < b float",
			args:  args{a: 1.3, b: 3.5, datatype: "real"},
			want:  -1,
			want1: true,
		},
		{
			name:  "a == b string",
			args:  args{a: "hello", b: "world", datatype: "text"},
			want:  0,
			want1: true,
		},
		{
			name:  "a > b +ve int",
			args:  args{a: uint64(6), b: uint64(5), datatype: "positiveInt"},
			want:  1,
			want1: true,
		},
		{
			name:  "a > b date",
			args:  args{a: date1, b: date2, datatype: "date"},
			want:  1,
			want1: true,
		},
		{
			name:  "a == b time",
			args:  args{a: time1, b: time2, datatype: "time"},
			want:  0,
			want1: true,
		},
		{
			name:  "a < b datetime",
			args:  args{a: datetime1, b: datetime2, datatype: "timestamptz"},
			want:  -1,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := compareTypeValues(tt.args.a, tt.args.b, tt.args.datatype)
			if got != tt.want {
				t.Errorf("compareTypeValues() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("compareTypeValues() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestColumn_validateEnums(t *testing.T) {
	type fields struct {
		ColumnName    string
		DataType      string
		NotNull       bool
		Unique        bool
		Min           string
		Max           string
		Enums         []interface{}
		Default       interface{}
		ForeignTable  string
		ForeignField  string
		minIndividual interface{}
		maxIndividual interface{}
		minArrLen     int
		maxArrLen     int
		enumMap       map[any]bool
		values        map[string]bool
		lookup        map[string]int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "valid int enum",
			fields:  fields{Enums: []any{1, 2, 3}, DataType: "integer", minIndividual: 1},
			wantErr: false,
		},
		{
			name:    "invalid float enum",
			fields:  fields{Enums: []any{1.5, 2.6}, DataType: "real", maxIndividual: 2},
			wantErr: true,
		},
		{
			name:    "invalid string enum",
			fields:  fields{Enums: []any{"hi", "bro"}, DataType: "text", maxIndividual: 2},
			wantErr: true,
		},
		{
			name:    "valid array string enum",
			fields:  fields{Enums: []any{"hi", "bro"}, DataType: "text[]"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := &Column{
				ColumnName:    tt.fields.ColumnName,
				DataType:      tt.fields.DataType,
				NotNull:       tt.fields.NotNull,
				Unique:        tt.fields.Unique,
				Min:           tt.fields.Min,
				Max:           tt.fields.Max,
				Enums:         tt.fields.Enums,
				Default:       tt.fields.Default,
				ForeignTable:  tt.fields.ForeignTable,
				ForeignField:  tt.fields.ForeignField,
				minIndividual: tt.fields.minIndividual,
				maxIndividual: tt.fields.maxIndividual,
				minArrLen:     tt.fields.minArrLen,
				maxArrLen:     tt.fields.maxArrLen,
				enumMap:       tt.fields.enumMap,
				values:        tt.fields.values,
				lookup:        tt.fields.lookup,
			}
			if err := column.validateEnums(); (err != nil) != tt.wantErr {
				t.Errorf("Column.validateEnums() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestColumn_validateValueByConstraints(t *testing.T) {
	type fields struct {
		ColumnName    string
		DataType      string
		NotNull       bool
		Unique        bool
		Min           string
		Max           string
		Enums         []interface{}
		Default       interface{}
		ForeignTable  string
		ForeignField  string
		minIndividual interface{}
		maxIndividual interface{}
		minArrLen     int
		maxArrLen     int
		enumMap       map[any]bool
		values        map[string]bool
		lookup        map[string]int
	}
	type args struct {
		value  any
		insert bool
	}

	date1, _ := time.Parse(time.DateOnly, "2024-01-01")
	date2, _ := time.Parse(time.DateOnly, "2024-02-01")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name:    "invalid string",
			fields:  fields{DataType: "text", enumMap: map[any]bool{"admin": true}},
			args:    args{value: "user", insert: false},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid datetime",
			fields:  fields{DataType: "timestamptz"},
			args:    args{value: date1, insert: false},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid date",
			fields:  fields{DataType: "date", minIndividual: date2},
			args:    args{value: date1, insert: false},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid int arr",
			fields:  fields{DataType: "integer[]", minArrLen: 3, maxArrLen: 3},
			args:    args{value: "[1, 2, 3]", insert: false},
			want:    []any{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "invalid real arr",
			fields:  fields{DataType: "real[]", minIndividual: 3.0},
			args:    args{value: "[1.2, 2.5, 3.3]", insert: false},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := &Column{
				ColumnName:    tt.fields.ColumnName,
				DataType:      tt.fields.DataType,
				NotNull:       tt.fields.NotNull,
				Unique:        tt.fields.Unique,
				Min:           tt.fields.Min,
				Max:           tt.fields.Max,
				Enums:         tt.fields.Enums,
				Default:       tt.fields.Default,
				ForeignTable:  tt.fields.ForeignTable,
				ForeignField:  tt.fields.ForeignField,
				minIndividual: tt.fields.minIndividual,
				maxIndividual: tt.fields.maxIndividual,
				minArrLen:     tt.fields.minArrLen,
				maxArrLen:     tt.fields.maxArrLen,
				enumMap:       tt.fields.enumMap,
				values:        tt.fields.values,
				lookup:        tt.fields.lookup,
			}
			got, err := column.validateValueByConstraints(tt.args.value, tt.args.insert)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.validateValueByConstraints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Column.validateValueByConstraints() = %v, want %v", got, tt.want)
			}
		})
	}
}
