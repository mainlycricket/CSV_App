package main

import "github.com/golang-jwt/jwt/v5"

type Column struct {
	ColumnName    string
	DataType      string
	NotNull       bool
	minIndividual interface{}
	maxIndividual interface{}
	minArrLen     int // 0 indicates unset or non-array type
	maxArrLen     int // 0 indicates unset or non-array type
	Enums         []interface{}
	pgType        string
}

type Table_students struct {
	Column_Branch_Id      CustomNullInt    `json:"Branch_Id"`
	Column_Course_Id      CustomNullInt    `json:"Course_Id"`
	Column_Student_Father CustomNullString `json:"Student_Father"`
	Column_Student_Id     CustomNullInt    `json:"Student_Id"`
	Column_Student_Name   CustomNullString `json:"Student_Name"`
}

type Table_students_response struct {
	Fkey_Branch_Id        Table_branches   `json:"Branch_Id"`
	Fkey_Course_Id        Table_courses    `json:"Course_Id"`
	Column_Student_Father CustomNullString `json:"Student_Father"`
	Column_Student_Id     CustomNullInt    `json:"Student_Id"`
	Column_Student_Name   CustomNullString `json:"Student_Name"`
}

var Map_students = map[string]Column{

	"Branch_Id": {
		ColumnName: "Branch_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Course_Id": {
		ColumnName: "Course_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Student_Father": {
		ColumnName: "Student_Father",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
	},

	"Student_Id": {
		ColumnName: "Student_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Student_Name": {
		ColumnName: "Student_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
	},
}

type Table_subjects struct {
	Column_Branch_Id    CustomNullInt    `json:"Branch_Id"`
	Column_Subject_Id   CustomNullInt    `json:"Subject_Id"`
	Column_Subject_Name CustomNullString `json:"Subject_Name"`
}

type Table_subjects_response struct {
	Fkey_Branch_Id      Table_branches   `json:"Branch_Id"`
	Column_Subject_Id   CustomNullInt    `json:"Subject_Id"`
	Column_Subject_Name CustomNullString `json:"Subject_Name"`
}

var Map_subjects = map[string]Column{

	"Branch_Id": {
		ColumnName: "Branch_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Subject_Id": {
		ColumnName: "Subject_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Subject_Name": {
		ColumnName: "Subject_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
	},
}

type Table_TypeTest struct {
	ID__                CustomNullInt        `json:"__ID"`
	Column_Bool         CustomNullBool       `json:"Bool"`
	Column_Bool_Arr     []CustomNullBool     `json:"Bool_Arr"`
	Column_Date         CustomNullDate       `json:"Date"`
	Column_DateTime     CustomNullDateTime   `json:"DateTime"`
	Column_Date_arr     []CustomNullDate     `json:"Date_arr"`
	Column_Datetime_Arr []CustomNullDateTime `json:"Datetime_Arr"`
	Column_Float        CustomNullFloat      `json:"Float"`
	Column_Float_arr    []CustomNullFloat    `json:"Float_arr"`
	Column_Int          CustomNullInt        `json:"Int"`
	Column_Int_Arr      []CustomNullInt      `json:"Int_Arr"`
	Column_Str_Arr      []CustomNullString   `json:"Str_Arr"`
	Column_String       CustomNullString     `json:"String"`
	Column_Time         CustomNullTime       `json:"Time"`
	Column_Time_Arr     []CustomNullTime     `json:"Time_Arr"`
}

type Table_TypeTest_response struct {
	ID__                CustomNullInt        `json:"__ID"`
	Column_Bool         CustomNullBool       `json:"Bool"`
	Column_Bool_Arr     []CustomNullBool     `json:"Bool_Arr"`
	Column_Date         CustomNullDate       `json:"Date"`
	Column_DateTime     CustomNullDateTime   `json:"DateTime"`
	Column_Date_arr     []CustomNullDate     `json:"Date_arr"`
	Column_Datetime_Arr []CustomNullDateTime `json:"Datetime_Arr"`
	Column_Float        CustomNullFloat      `json:"Float"`
	Column_Float_arr    []CustomNullFloat    `json:"Float_arr"`
	Column_Int          CustomNullInt        `json:"Int"`
	Column_Int_Arr      []CustomNullInt      `json:"Int_Arr"`
	Column_Str_Arr      []CustomNullString   `json:"Str_Arr"`
	Column_String       CustomNullString     `json:"String"`
	Column_Time         CustomNullTime       `json:"Time"`
	Column_Time_Arr     []CustomNullTime     `json:"Time_Arr"`
}

var Map_TypeTest = map[string]Column{

	"Bool": {
		ColumnName: "Bool",
		DataType:   "CustomNullBool",
		NotNull:    false,
		pgType:     "boolean",
	},

	"Bool_Arr": {
		ColumnName: "Bool_Arr",
		DataType:   "[]CustomNullBool",
		NotNull:    false,
		pgType:     "boolean[]",
	},

	"Date": {
		ColumnName: "Date",
		DataType:   "CustomNullDate",
		NotNull:    false,
		pgType:     "date",
	},

	"DateTime": {
		ColumnName: "DateTime",
		DataType:   "CustomNullDateTime",
		NotNull:    false,
		pgType:     "timestamptz",
	},

	"Date_arr": {
		ColumnName: "Date_arr",
		DataType:   "[]CustomNullDate",
		NotNull:    false,
		pgType:     "date[]",
	},

	"Datetime_Arr": {
		ColumnName: "Datetime_Arr",
		DataType:   "[]CustomNullDateTime",
		NotNull:    false,
		pgType:     "timestamptz[]",
	},

	"Float": {
		ColumnName: "Float",
		DataType:   "CustomNullFloat",
		NotNull:    false,
		pgType:     "real",
	},

	"Float_arr": {
		ColumnName: "Float_arr",
		DataType:   "[]CustomNullFloat",
		NotNull:    false,
		pgType:     "real[]",
	},

	"Int": {
		ColumnName: "Int",
		DataType:   "CustomNullInt",
		NotNull:    false,
		pgType:     "integer",
	},

	"Int_Arr": {
		ColumnName: "Int_Arr",
		DataType:   "[]CustomNullInt",
		NotNull:    false,
		pgType:     "integer[]",
	},

	"Str_Arr": {
		ColumnName: "Str_Arr",
		DataType:   "[]CustomNullString",
		NotNull:    false,
		pgType:     "text[]",
	},

	"String": {
		ColumnName: "String",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
	},

	"Time": {
		ColumnName: "Time",
		DataType:   "CustomNullTime",
		NotNull:    false,
		pgType:     "time",
	},

	"Time_Arr": {
		ColumnName: "Time_Arr",
		DataType:   "[]CustomNullTime",
		NotNull:    false,
		pgType:     "time[]",
	},
}

type Table_branches struct {
	Column_Branch_Id   CustomNullInt      `json:"Branch_Id"`
	Column_Branch_Name CustomNullString   `json:"Branch_Name"`
	Column_Course_Id   CustomNullInt      `json:"Course_Id"`
	Column_Teachers    []CustomNullString `json:"Teachers"`
}

type Table_branches_response struct {
	Column_Branch_Id   CustomNullInt      `json:"Branch_Id"`
	Column_Branch_Name CustomNullString   `json:"Branch_Name"`
	Fkey_Course_Id     Table_courses      `json:"Course_Id"`
	Column_Teachers    []CustomNullString `json:"Teachers"`
}

var Map_branches = map[string]Column{

	"Branch_Id": {
		ColumnName: "Branch_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Branch_Name": {
		ColumnName: "Branch_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
	},

	"Course_Id": {
		ColumnName: "Course_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Teachers": {
		ColumnName: "Teachers",
		DataType:   "[]CustomNullString",
		NotNull:    false,
		pgType:     "text[]",
	},
}

type Table_courses struct {
	Column_Course_Id       CustomNullInt    `json:"Course_Id"`
	Column_Course_Name     CustomNullString `json:"Course_Name"`
	Column_Lateral_Allowed CustomNullBool   `json:"Lateral_Allowed"`
}

type Table_courses_response struct {
	Column_Course_Id       CustomNullInt    `json:"Course_Id"`
	Column_Course_Name     CustomNullString `json:"Course_Name"`
	Column_Lateral_Allowed CustomNullBool   `json:"Lateral_Allowed"`
}

var Map_courses = map[string]Column{

	"Course_Id": {
		ColumnName: "Course_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
	},

	"Course_Name": {
		ColumnName: "Course_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
	},

	"Lateral_Allowed": {
		ColumnName: "Lateral_Allowed",
		DataType:   "CustomNullBool",
		NotNull:    false,
		pgType:     "boolean",
	},
}

type Table_login struct {
	Column_password CustomNullString `json:"password"`
	Column_role     CustomNullString `json:"role"`
	Column_username CustomNullString `json:"username"`
}

type CustomJwtClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Login_Input struct {
	Username CustomNullString `json:"username"`
	Password CustomNullString `json:"password"`
}

type Login_Output struct {
	Username CustomNullString `json:"username"`
	Password CustomNullString `json:"password"`
	Role     CustomNullString `json:"role"`
}
