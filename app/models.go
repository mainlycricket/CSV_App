package main

import "time"

type Table_TypeTest struct {
	ID__                uint        `json:"__ID"`
	Column_Bool         bool        `json:"Bool"`
	Column_Bool_Arr     []bool      `json:"Bool_Arr"`
	Column_Date         time.Time   `json:"Date"`
	Column_DateTime     time.Time   `json:"DateTime"`
	Column_Date_arr     []time.Time `json:"Date_arr"`
	Column_Datetime_Arr []time.Time `json:"Datetime_Arr"`
	Column_Float        float64     `json:"Float"`
	Column_Float_arr    []float64   `json:"Float_arr"`
	Column_Int          int         `json:"Int"`
	Column_Int_Arr      []int       `json:"Int_Arr"`
	Column_Str_Arr      []string    `json:"Str_Arr"`
	Column_String       string      `json:"String"`
	Column_Time         time.Time   `json:"Time"`
	Column_Time_Arr     []time.Time `json:"Time_Arr"`
}

type Table_branches struct {
	Column_Branch_Id   int      `json:"Branch_Id"`
	Column_Branch_Name string   `json:"Branch_Name"`
	Column_Course_Id   int      `json:"Course_Id"`
	Column_Teachers    []string `json:"Teachers"`
}

type Table_courses struct {
	Column_Course_Id       int    `json:"Course_Id"`
	Column_Course_Name     string `json:"Course_Name"`
	Column_Lateral_Allowed bool   `json:"Lateral_Allowed"`
}

type Table_empty struct {
	ID__         uint    `json:"__ID"`
	Column_Col_1 int     `json:"Col_1"`
	Column_Col_2 string  `json:"Col_2"`
	Column_Col_3 float64 `json:"Col_3"`
}

type Table_students struct {
	Column_Branch_Id      int    `json:"Branch_Id"`
	Column_Course_Id      int    `json:"Course_Id"`
	Column_Student_Father string `json:"Student_Father"`
	Column_Student_Id     int    `json:"Student_Id"`
	Column_Student_Name   string `json:"Student_Name"`
}

type Table_subjects struct {
	Column_Branch_Id    int    `json:"Branch_Id"`
	Column_Subject_Id   int    `json:"Subject_Id"`
	Column_Subject_Name string `json:"Subject_Name"`
}
