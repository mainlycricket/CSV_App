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
			name:         "string",
			args:         args{value: "text", datatype: "string"},
			convertedVal: "text",
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
			fields:  fields{DataType: "[]integer", Default: []any{1, 2}},
			wantErr: false,
		},
		{
			name:    "invalid int arr",
			fields:  fields{DataType: "[]integer", Default: []any{3, 4}, minArrLen: 3},
			wantErr: true,
		},
		{
			name:    "str arr",
			fields:  fields{DataType: "[]string", Default: []any{"te", "dd"}},
			wantErr: false,
		},
		{
			name:    "invalid str",
			fields:  fields{DataType: "string", Default: "text", Enums: []any{"other"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := &Column{
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
			}
			if err := column.validateDefaultValue(); (err != nil) != tt.wantErr {
				t.Errorf("Column.validateDefaultValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
