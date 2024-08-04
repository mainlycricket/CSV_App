package main

import "testing"

func Test_detectDataType(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "real_arr",
			args: args{value: "[1, -2, 4.5, 5, -7.5, 6.3, \"+3\", \"+5.68\"]"},
			want: "real[]",
		},
		{
			name: "mixed arr",
			args: args{value: "[1, 2, 4.5, false, 5, 3]"},
			want: "text[]",
		},
		{
			name: "bool arr",
			args: args{value: "[false, true, false]"},
			want: "boolean[]",
		},
		{
			name: "text arr",
			args: args{value: "[\"text1\", \"text2\", \"text3\"]"},
			want: "text[]",
		},
		{
			name: "date arr",
			args: args{value: "[\"2024-01-01\", \"2024-12-30\"]"},
			want: "date[]",
		},
		{
			name: "time arr",
			args: args{value: "[\"12:25:52\", \"00:12:30\"]"},
			want: "time[]",
		},
		{
			name: "datetime arr",
			args: args{value: "[\"2022-01-01T12:25:52+05:30\"]"},
			want: "timestamptz[]",
		},
		{
			name: "invalid date",
			args: args{value: "2024-02-30"},
			want: "text",
		},
		{
			name: "invalid time",
			args: args{value: "12:59:61"},
			want: "text",
		},
		{
			name: "invalid datetime",
			args: args{value: "2024-01-01T12:59:61+05:30"},
			want: "text",
		},
		{
			name: "text",
			args: args{value: "data"},
			want: "text",
		},
		{
			name: "integer",
			args: args{value: "+4"},
			want: "integer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectDataType(tt.args.value); got != tt.want {
				t.Errorf("detectDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}
