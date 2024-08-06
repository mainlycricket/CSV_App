package main

import "time"

type TypeTest struct {
	Bool         bool        `json:"Bool"`
	Bool_Arr     []bool      `json:"Bool_Arr"`
	Date         time.Time   `json:"Date"`
	DateTime     time.Time   `json:"DateTime"`
	Date_arr     []time.Time `json:"Date_arr"`
	Datetime_Arr []time.Time `json:"Datetime_Arr"`
	Float        float64     `json:"Float"`
	Float_arr    []float64   `json:"Float_arr"`
	Int          int         `json:"Int"`
	Int_Arr      []int       `json:"Int_Arr"`
	Str_Arr      []string    `json:"Str_Arr"`
	String       string      `json:"String"`
	Time         time.Time   `json:"Time"`
	Time_Arr     []time.Time `json:"Time_Arr"`
}

type branches struct {
	Branch_Id   int      `json:"Branch_Id"`
	Branch_Name string   `json:"Branch_Name"`
	Course_Id   int      `json:"Course_Id"`
	Teachers    []string `json:"Teachers"`
}

type courses struct {
	Course_Id       int    `json:"Course_Id"`
	Course_Name     string `json:"Course_Name"`
	Lateral_Allowed bool   `json:"Lateral_Allowed"`
}

type empty struct {
	Col_1 int     `json:"Col_1"`
	Col_2 string  `json:"Col_2"`
	Col_3 float64 `json:"Col_3"`
}

type students struct {
	Branch_Id      int    `json:"Branch_Id"`
	Course_Id      int    `json:"Course_Id"`
	Student_Father string `json:"Student_Father"`
	Student_Id     int    `json:"Student_Id"`
	Student_Name   string `json:"Student_Name"`
}

type subjects struct {
	Branch_Id    int    `json:"Branch_Id"`
	Subject_Id   int    `json:"Subject_Id"`
	Subject_Name string `json:"Subject_Name"`
}
