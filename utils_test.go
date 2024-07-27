package main

import (
	"fmt"
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

	var boolInterface interface{} = false
	var strInterface interface{} = "text"
	var dateInterface interface{} = "2024-01-01"
	dateRes, _ := time.Parse(datetimeFormats["date"], fmt.Sprintf("%s", dateInterface))

	tests := []struct {
		name         string
		args         args
		convertedVal any
		success      bool
	}{
		{
			name:         "date",
			args:         args{value: dateInterface, datatype: "date"},
			convertedVal: dateRes,
			success:      true,
		},
		{
			name:         "false boolean",
			args:         args{value: boolInterface, datatype: "boolean"},
			convertedVal: boolInterface,
			success:      true,
		},
		{
			name:         "string",
			args:         args{value: strInterface, datatype: "string"},
			convertedVal: strInterface,
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
