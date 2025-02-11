package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type CustomNullDate struct {
	sql.NullTime
}

func (t *CustomNullDate) Scan(src any) error {
	var err error
	switch src := src.(type) {
	case nil:
		t.Valid = false
	case time.Time:
		t.Time = src
		t.Valid = true
	case []byte:
		t.Time, err = time.Parse(time.DateOnly, string(src))
		if err == nil {
			t.Valid = true
		}
	default:
		err = errors.New("unsupported type")
	}
	return err
}

func (t CustomNullDate) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format(time.DateOnly) + `"`), nil
}

func (t *CustomNullDate) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		t.Valid = false
	} else {
		value, err := time.Parse(`"`+time.DateOnly+`"`, str)
		if err != nil {
			return err
		}
		t.Time = value
		t.Valid = true
	}

	return nil
}

func (t *CustomNullDate) GetValue() time.Time {
	return t.Time
}

func getCustomNullDateArrValues(arr []CustomNullDate) []time.Time {
	res := make([]time.Time, len(arr))

	for idx, value := range arr {
		res[idx] = value.GetValue()
	}

	return res
}

type CustomNullTime struct {
	sql.NullTime
}

func (t *CustomNullTime) Scan(src any) error {
	var err error
	switch src := src.(type) {
	case nil:
		t.Valid = false
	case time.Time:
		t.Time = src
		t.Valid = true
	case []byte:
		t.Time, err = time.Parse(time.TimeOnly, string(src))
		if err == nil {
			t.Valid = true
		}
	default:
		err = errors.New("unsupported type")
	}
	return err
}

func (t CustomNullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format(time.TimeOnly) + `"`), nil
}

func (t *CustomNullTime) UnmarshalJSON(data []byte) error {
	str := string(data)

	if str == "null" {
		t.Valid = false
	} else {
		value, err := time.Parse(`"`+time.TimeOnly+`"`, str)
		if err != nil {
			return err
		}
		t.Time = value
		t.Valid = true
	}

	return nil
}

func (t *CustomNullTime) GetValue() time.Time {
	return t.Time
}

func getCustomNullTimeArrValues(arr []CustomNullTime) []time.Time {
	res := make([]time.Time, len(arr))

	for idx, value := range arr {
		res[idx] = value.GetValue()
	}

	return res
}

type CustomNullDateTime struct {
	sql.NullTime
}

func (t *CustomNullDateTime) Scan(src any) error {
	var err error
	switch src := src.(type) {
	case nil:
		t.Valid = false
	case time.Time:
		t.Time = src
		t.Valid = true
	case []byte:
		t.Time, err = time.Parse("2006-01-02 15:04:05-07:00", string(src))
		if err == nil {
			t.Valid = true
		}
	default:
		err = errors.New("unsupported type")
	}
	return err
}

func (t CustomNullDateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format(time.RFC3339) + `"`), nil
}

func (t *CustomNullDateTime) UnmarshalJSON(data []byte) error {
	str := string(data)

	if str == "null" {
		t.Valid = false
	} else {
		value, err := time.Parse(`"`+time.RFC3339+`"`, str)
		if err != nil {
			return err
		}
		t.Time = value
		t.Valid = true
	}

	return nil
}

func (t *CustomNullDateTime) GetValue() time.Time {
	return t.Time
}

func getCustomNullDateTimeArrValues(arr []CustomNullDateTime) []time.Time {
	res := make([]time.Time, len(arr))

	for idx, value := range arr {
		res[idx] = value.GetValue()
	}

	return res
}

type CustomNullBool struct {
	sql.NullBool
}

func (t CustomNullBool) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%v", t.Bool)), nil
}

func (t *CustomNullBool) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		t.Valid = false
	} else {
		val, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		t.Valid = true
		t.Bool = val
	}

	return nil
}

func (t *CustomNullBool) GetValue() bool {
	return t.Bool
}

func getCustomNullBoolArrValues(arr []CustomNullBool) []bool {
	res := make([]bool, len(arr))

	for idx, value := range arr {
		res[idx] = value.GetValue()
	}

	return res
}

type CustomNullFloat struct {
	sql.NullFloat64
}

func (t CustomNullFloat) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%v", t.Float64)), nil
}

func (t *CustomNullFloat) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		t.Valid = false
	} else {
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		t.Valid = true
		t.Float64 = val
	}

	return nil
}

func (t *CustomNullFloat) GetValue() float64 {
	return t.Float64
}

func getCustomNullFloatArrValues(arr []CustomNullFloat) []float64 {
	res := make([]float64, len(arr))

	for idx, value := range arr {
		res[idx] = value.GetValue()
	}

	return res
}

type CustomNullInt struct {
	sql.NullInt64
}

func (t CustomNullInt) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%v", t.Int64)), nil
}

func (t *CustomNullInt) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		t.Valid = false
	} else {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		t.Valid = true
		t.Int64 = val
	}

	return nil
}

func (t *CustomNullInt) GetValue() int64 {
	return t.Int64
}

func getCustomNullIntArrValues(arr []CustomNullInt) []int64 {
	res := make([]int64, len(arr))

	for idx, value := range arr {
		res[idx] = value.GetValue()
	}

	return res
}

type CustomNullString struct {
	sql.NullString
}

func (t CustomNullString) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%v"`, t.String)), nil
}

func (t *CustomNullString) UnmarshalJSON(data []byte) error {
	value := string(data)

	if value == "null" {
		t.Valid = false
	} else {
		t.String = value[1 : len(value)-1]
		t.Valid = true
	}

	return nil
}

func (t *CustomNullString) GetValue() string {
	return t.String
}

func getCustomNullStringArrValues(arr []CustomNullString) []string {
	res := make([]string, len(arr))

	for idx, value := range arr {
		res[idx] = value.GetValue()
	}

	return res
}
