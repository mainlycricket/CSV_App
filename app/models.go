package main

import "github.com/golang-jwt/jwt/v5"

type Column struct {
	ColumnName    string
	DataType      string
	NotNull       bool
	Hash          bool
	minIndividual interface{}
	maxIndividual interface{}
	minArrLen     int // 0 indicates unset or non-array type
	maxArrLen     int // 0 indicates unset or non-array type
	Enums         []interface{}
	pgType        string
}

type Table_college struct {
	Column_college_id   CustomNullString `json:"college_id"`
	Column_college_name CustomNullString `json:"college_name"`
	Column_principal_id CustomNullString `json:"principal_id"`
}

var Map_college = map[string]Column{

	"college_id": {
		ColumnName: "college_id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"college_name": {
		ColumnName: "college_name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"principal_id": {
		ColumnName: "principal_id",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
		Hash:       false,
	},
}

type Table_college_response struct {
	Column_college_id   CustomNullString `json:"college_id"`
	Column_college_name CustomNullString `json:"college_name"`
	Fkey_principal_id   Table_login      `json:"principal_id"`
}

type Table_courses struct {
	Column_Course_Id       CustomNullString `json:"Course_Id"`
	Column_Course_Name     CustomNullString `json:"Course_Name"`
	Column_Lateral_Allowed CustomNullBool   `json:"Lateral_Allowed"`
	Column_added_by        CustomNullString `json:"added_by"`
	Column_college_id      CustomNullString `json:"college_id"`
}

var Map_courses = map[string]Column{

	"Course_Id": {
		ColumnName: "Course_Id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Course_Name": {
		ColumnName: "Course_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Lateral_Allowed": {
		ColumnName: "Lateral_Allowed",
		DataType:   "CustomNullBool",
		NotNull:    false,
		pgType:     "boolean",
		Hash:       false,
	},

	"added_by": {
		ColumnName: "added_by",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"college_id": {
		ColumnName: "college_id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},
}

type Table_courses_response struct {
	Column_Course_Id       CustomNullString `json:"Course_Id"`
	Column_Course_Name     CustomNullString `json:"Course_Name"`
	Column_Lateral_Allowed CustomNullBool   `json:"Lateral_Allowed"`
	Fkey_added_by          Table_login      `json:"added_by"`
	Fkey_college_id        Table_college    `json:"college_id"`
}

type Table_login struct {
	Column_added_by   CustomNullString `json:"added_by"`
	Column_branch_id  CustomNullString `json:"branch_id"`
	Column_college_id CustomNullString `json:"college_id"`
	Column_course_id  CustomNullString `json:"course_id"`
	Column_password   CustomNullString `json:"password"`
	Column_role       CustomNullString `json:"role"`
	Column_username   CustomNullString `json:"username"`
}

var Map_login = map[string]Column{

	"added_by": {
		ColumnName: "added_by",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
		Hash:       false,
	},

	"branch_id": {
		ColumnName: "branch_id",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
		Hash:       false,
	},

	"college_id": {
		ColumnName: "college_id",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
		Hash:       false,
	},

	"course_id": {
		ColumnName: "course_id",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
		Hash:       false,
	},

	"password": {
		ColumnName: "password",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       true,
	},

	"role": {
		ColumnName: "role",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"username": {
		ColumnName: "username",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},
}

type CustomJwtClaims struct {
	Username   string `json:"username"`
	Role       string `json:"role"`
	College_id string `json:"college_id"`
	Course_id  string `json:"course_id"`
	Branch_id  string `json:"branch_id"`
	jwt.RegisteredClaims
}

type Login_Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login_Output struct {
	Username   CustomNullString `json:"username"`
	Password   CustomNullString `json:"password"`
	Role       CustomNullString `json:"role"`
	College_id CustomNullString `json:"college_id"`
	Course_id  CustomNullString `json:"course_id"`
	Branch_id  CustomNullString `json:"branch_id"`
}

type ContextKey string

type Table_students struct {
	Column_Branch_Id      CustomNullString `json:"Branch_Id"`
	Column_Course_Id      CustomNullString `json:"Course_Id"`
	Column_Student_Father CustomNullString `json:"Student_Father"`
	Column_Student_Id     CustomNullInt    `json:"Student_Id"`
	Column_Student_Name   CustomNullString `json:"Student_Name"`
	Column_added_by       CustomNullString `json:"added_by"`
	Column_college_id     CustomNullString `json:"college_id"`
}

var Map_students = map[string]Column{

	"Branch_Id": {
		ColumnName: "Branch_Id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Course_Id": {
		ColumnName: "Course_Id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Student_Father": {
		ColumnName: "Student_Father",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Student_Id": {
		ColumnName: "Student_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
		Hash:       false,
	},

	"Student_Name": {
		ColumnName: "Student_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"added_by": {
		ColumnName: "added_by",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"college_id": {
		ColumnName: "college_id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},
}

type Table_students_response struct {
	Fkey_Branch_Id        Table_branches   `json:"Branch_Id"`
	Fkey_Course_Id        Table_courses    `json:"Course_Id"`
	Column_Student_Father CustomNullString `json:"Student_Father"`
	Column_Student_Id     CustomNullInt    `json:"Student_Id"`
	Column_Student_Name   CustomNullString `json:"Student_Name"`
	Fkey_added_by         Table_login      `json:"added_by"`
	Fkey_college_id       Table_college    `json:"college_id"`
}

type Table_subjects struct {
	Column_Branch_Id    CustomNullString `json:"Branch_Id"`
	Column_Subject_Id   CustomNullInt    `json:"Subject_Id"`
	Column_Subject_Name CustomNullString `json:"Subject_Name"`
	Column_added_by     CustomNullString `json:"added_by"`
	Column_college_id   CustomNullString `json:"college_id"`
	Column_course_id    CustomNullString `json:"course_id"`
}

var Map_subjects = map[string]Column{

	"Branch_Id": {
		ColumnName: "Branch_Id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Subject_Id": {
		ColumnName: "Subject_Id",
		DataType:   "CustomNullInt",
		NotNull:    true,
		pgType:     "integer",
		Hash:       false,
	},

	"Subject_Name": {
		ColumnName: "Subject_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"added_by": {
		ColumnName: "added_by",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"college_id": {
		ColumnName: "college_id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"course_id": {
		ColumnName: "course_id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},
}

type Table_subjects_response struct {
	Fkey_Branch_Id      Table_branches   `json:"Branch_Id"`
	Column_Subject_Id   CustomNullInt    `json:"Subject_Id"`
	Column_Subject_Name CustomNullString `json:"Subject_Name"`
	Fkey_added_by       Table_login      `json:"added_by"`
	Fkey_college_id     Table_college    `json:"college_id"`
	Fkey_course_id      Table_courses    `json:"course_id"`
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

var Map_TypeTest = map[string]Column{

	"Bool": {
		ColumnName: "Bool",
		DataType:   "CustomNullBool",
		NotNull:    false,
		pgType:     "boolean",
		Hash:       false,
	},

	"Bool_Arr": {
		ColumnName: "Bool_Arr",
		DataType:   "[]CustomNullBool",
		NotNull:    false,
		pgType:     "boolean[]",
		Hash:       false,
	},

	"Date": {
		ColumnName: "Date",
		DataType:   "CustomNullDate",
		NotNull:    false,
		pgType:     "date",
		Hash:       false,
	},

	"DateTime": {
		ColumnName: "DateTime",
		DataType:   "CustomNullDateTime",
		NotNull:    false,
		pgType:     "timestamptz",
		Hash:       false,
	},

	"Date_arr": {
		ColumnName: "Date_arr",
		DataType:   "[]CustomNullDate",
		NotNull:    false,
		pgType:     "date[]",
		Hash:       false,
	},

	"Datetime_Arr": {
		ColumnName: "Datetime_Arr",
		DataType:   "[]CustomNullDateTime",
		NotNull:    false,
		pgType:     "timestamptz[]",
		Hash:       false,
	},

	"Float": {
		ColumnName: "Float",
		DataType:   "CustomNullFloat",
		NotNull:    false,
		pgType:     "real",
		Hash:       false,
	},

	"Float_arr": {
		ColumnName: "Float_arr",
		DataType:   "[]CustomNullFloat",
		NotNull:    false,
		pgType:     "real[]",
		Hash:       false,
	},

	"Int": {
		ColumnName: "Int",
		DataType:   "CustomNullInt",
		NotNull:    false,
		pgType:     "integer",
		Hash:       false,
	},

	"Int_Arr": {
		ColumnName: "Int_Arr",
		DataType:   "[]CustomNullInt",
		NotNull:    false,
		pgType:     "integer[]",
		Hash:       false,
	},

	"Str_Arr": {
		ColumnName: "Str_Arr",
		DataType:   "[]CustomNullString",
		NotNull:    false,
		pgType:     "text[]",
		Hash:       false,
	},

	"String": {
		ColumnName: "String",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
		Hash:       false,
	},

	"Time": {
		ColumnName: "Time",
		DataType:   "CustomNullTime",
		NotNull:    false,
		pgType:     "time",
		Hash:       false,
	},

	"Time_Arr": {
		ColumnName: "Time_Arr",
		DataType:   "[]CustomNullTime",
		NotNull:    false,
		pgType:     "time[]",
		Hash:       false,
	},
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

type Table_branches struct {
	Column_Branch_Id   CustomNullString   `json:"Branch_Id"`
	Column_Branch_Name CustomNullString   `json:"Branch_Name"`
	Column_Course_Id   CustomNullString   `json:"Course_Id"`
	Column_HoD         CustomNullString   `json:"HoD"`
	Column_Teachers    []CustomNullString `json:"Teachers"`
	Column_added_by    CustomNullString   `json:"added_by"`
	Column_college_id  CustomNullString   `json:"college_id"`
}

var Map_branches = map[string]Column{

	"Branch_Id": {
		ColumnName: "Branch_Id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Branch_Name": {
		ColumnName: "Branch_Name",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"Course_Id": {
		ColumnName: "Course_Id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"HoD": {
		ColumnName: "HoD",
		DataType:   "CustomNullString",
		NotNull:    false,
		pgType:     "text",
		Hash:       false,
	},

	"Teachers": {
		ColumnName: "Teachers",
		DataType:   "[]CustomNullString",
		NotNull:    false,
		pgType:     "text[]",
		Hash:       false,
	},

	"added_by": {
		ColumnName: "added_by",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},

	"college_id": {
		ColumnName: "college_id",
		DataType:   "CustomNullString",
		NotNull:    true,
		pgType:     "text",
		Hash:       false,
	},
}

type Table_branches_response struct {
	Column_Branch_Id   CustomNullString   `json:"Branch_Id"`
	Column_Branch_Name CustomNullString   `json:"Branch_Name"`
	Fkey_Course_Id     Table_courses      `json:"Course_Id"`
	Fkey_HoD           Table_login        `json:"HoD"`
	Column_Teachers    []CustomNullString `json:"Teachers"`
	Fkey_added_by      Table_login        `json:"added_by"`
	Fkey_college_id    Table_college      `json:"college_id"`
}
