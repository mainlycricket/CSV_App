package main

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

type Table_branches struct {
	Column_Branch_Id   CustomNullInt      `json:"Branch_Id"`
	Column_Branch_Name CustomNullString   `json:"Branch_Name"`
	Column_Course_Id   CustomNullInt      `json:"Course_Id"`
	Column_Teachers    []CustomNullString `json:"Teachers"`
}

type Table_courses struct {
	Column_Course_Id       CustomNullInt    `json:"Course_Id"`
	Column_Course_Name     CustomNullString `json:"Course_Name"`
	Column_Lateral_Allowed CustomNullBool   `json:"Lateral_Allowed"`
}

type Table_students struct {
	Column_Branch_Id      CustomNullInt    `json:"Branch_Id"`
	Column_Course_Id      CustomNullInt    `json:"Course_Id"`
	Column_Student_Father CustomNullString `json:"Student_Father"`
	Column_Student_Id     CustomNullInt    `json:"Student_Id"`
	Column_Student_Name   CustomNullString `json:"Student_Name"`
}

type Table_subjects struct {
	Column_Branch_Id    CustomNullInt    `json:"Branch_Id"`
	Column_Subject_Id   CustomNullInt    `json:"Subject_Id"`
	Column_Subject_Name CustomNullString `json:"Subject_Name"`
}
