package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var basicTypes = map[string]interface{}{
	"integer": int(0),
	"real":    float64(0),
	"text":    "",
	"boolean": false,
}

var datetimeFormats = map[string]string{
	"date":        time.DateOnly, // 2024-01-01
	"time":        time.TimeOnly, // 14:30:00
	"timestamptz": time.RFC3339,  // 2024-07-01T12:30:00+05:30
}

var typeConversionFuncs = map[string]func(string) (any, error){
	"text": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		return fmt.Sprintf("%v", value), nil
	},
	"integer": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		return strconv.Atoi(value)
	},
	"real": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		return strconv.ParseFloat(value, 64)
	},
	"boolean": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		return strconv.ParseBool(value)
	},
	"positiveInt": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		if parsed <= 0 {
			return nil, errors.New("only positive integer is allowed")
		}
		return parsed, err
	},
	"date": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		return time.Parse(datetimeFormats["date"], value)
	},
	"time": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		return time.Parse(datetimeFormats["time"], value)
	},
	"timestamptz": func(value string) (any, error) {
		if len(value) == 0 {
			return nil, nil
		}
		return time.Parse(datetimeFormats["timestamptz"], value)
	},
}

// used to validate Column.DataType
func isValidTypeName(datatype string) bool {
	datatype = strings.TrimSuffix(datatype, "[]")
	_, isBasic := basicTypes[datatype]
	_, isTime := datetimeFormats[datatype]
	return isBasic || isTime
}

// validates and sets min max constraints (array & individual)
func (column *Column) setMinMaxConstraint() error {
	datatype := strings.TrimSpace(column.DataType)
	column.Min = strings.TrimSpace(column.Min)
	column.Max = strings.TrimSpace(column.Max)
	min := column.Min
	max := column.Max

	if strings.HasSuffix(datatype, "[]") {
		minArr := strings.SplitN(min, ",", 2)
		maxArr := strings.SplitN(max, ",", 2)

		// length
		minLenInterface, ok := validateValueByType(minArr[0], "positiveInt")
		if !ok {
			return errors.New("invalid min array length")
		}

		maxLenInterface, ok := validateValueByType(maxArr[0], "positiveInt")
		if !ok {
			return errors.New("invalid max array length")
		}

		if minLenInterface != nil {
			column.minArrLen = minLenInterface.(int)
		}

		if maxLenInterface != nil {
			column.maxArrLen = maxLenInterface.(int)
		}

		if column.minArrLen > column.maxArrLen {
			return errors.New("array min length can't be greater than max length")
		}

		// set individual min & max
		min = ""
		if len(minArr) == 2 {
			min = strings.TrimSpace(minArr[1])
		}

		max = ""
		if len(maxArr) == 2 {
			max = strings.TrimSpace(maxArr[1])
		}

		datatype = strings.TrimSuffix(datatype, "[]")
	}

	// Individual Constraints

	if datatype == "boolean" && (len(min) > 0 || len(max) > 0) {
		return errors.New("individual boolean values can't have min/max constraints")
	}

	if datatype == "text" {
		datatype = "positiveInt"
	}

	minInterface, ok := validateValueByType(min, datatype)
	if !ok {
		return errors.New("invalid min individual value")
	}

	maxInterface, ok := validateValueByType(max, datatype)
	if !ok {
		return errors.New("invalid max individual value")
	}

	column.minIndividual = minInterface
	column.maxIndividual = maxInterface

	if minInterface != nil && maxInterface != nil {
		if res, ok := compareTypeValues(minInterface, maxInterface, datatype); !ok || res == 1 {
			return errors.New("min value can't be greater than max value")
		}
	}

	return nil
}

// Individual elements are validated for array types
func (column *Column) validateEnums() error {
	datatype := strings.TrimSuffix(column.DataType, "[]")
	column.enumMap = make(map[any]bool)

	for idx, value := range column.Enums {
		interfaceVal, ok := validateValueByType(value, datatype)
		column.Enums[idx] = interfaceVal
		if !ok {
			errorMessage := fmt.Sprintf("%v is not of %v datatype", value, datatype)
			return errors.New(errorMessage)
		}

		if err := column.validateValueByMinMax(interfaceVal); err != nil {
			return err
		}

		column.enumMap[interfaceVal] = true
	}

	return nil
}

func (column *Column) validateDefaultValue() error {
	if column.Default == nil {
		return nil
	}

	interfaceVal, err := column.validateValueByConstraints(column.Default)
	if err != nil {
		return err
	}

	column.Default = interfaceVal

	return nil
}

/*
Receives value and converts the value in the specified type.
Returns the parsed value in interface type and bool indicates success/failure.
If length of value is 0, true is returned.
Array types aren't supported.
*/
func validateValueByType(value any, datatype string) (any, bool) {
	strVal := fmt.Sprintf("%v", value)
	strVal = strings.TrimSpace(strVal)

	if len(strVal) == 0 {
		return nil, true
	}

	convFunc, ok := typeConversionFuncs[datatype]
	if !ok {
		return nil, false
	}

	convertedInterface, err := convFunc(strVal)
	if err != nil {
		return nil, false
	}

	var parsed interface{}

	switch datatype {
	case "integer":
		parsed, ok = convertedInterface.(int)
	case "real":
		parsed, ok = convertedInterface.(float64)
	case "text":
		parsed, ok = convertedInterface.(string)
	case "boolean":
		parsed, ok = convertedInterface.(bool)
	case "date":
		parsed, ok = convertedInterface.(time.Time)
	case "time":
		parsed, ok = convertedInterface.(time.Time)
	case "timestamptz":
		parsed, ok = convertedInterface.(time.Time)
	case "positiveInt":
		parsed, ok := convertedInterface.(int)
		if !ok || parsed <= 0 {
			return nil, false
		}
		return parsed, ok
	}

	return parsed, ok
}

func (column *Column) validateValueByConstraints(value any) (any, error) {
	if strings.HasSuffix(column.DataType, "[]") {
		datatype := strings.TrimSuffix(column.DataType, "[]")

		interfaceArr, err := column.validateValArrLen(value)
		if err != nil {
			return nil, err
		}

		for idx, value := range interfaceArr {
			interfaceVal, ok := validateValueByType(value, datatype)
			if !ok {
				errorMessage := fmt.Sprintf("%v is not of %v datatype", value, column.DataType)
				return nil, errors.New(errorMessage)
			}

			if err := column.validateValueByMinMax(interfaceVal); err != nil {
				return nil, err
			}

			if err := column.validateValueByEnum(interfaceVal); err != nil {
				return nil, err
			}

			interfaceArr[idx] = interfaceVal
		}

		return interfaceArr, nil
	}

	interfaceVal, ok := validateValueByType(column.Default, column.DataType)
	if !ok {
		errorMessage := fmt.Sprintf("%v should be of %v type", value, column.DataType)
		return nil, errors.New(errorMessage)
	}

	if err := column.validateValueByMinMax(value); err != nil {
		return nil, err
	}

	if err := column.validateValueByEnum(value); err != nil {
		return nil, err
	}

	return interfaceVal, nil
}

// returns if
func (column *Column) validateValArrLen(value any) ([]any, error) {
	interfaceArr, ok := value.([]any)
	if !ok {
		return nil, errors.New("failed to convert to array")
	}

	if column.minArrLen != 0 {
		res, ok := compareTypeValues(len(interfaceArr), column.minArrLen, "integer")
		if !ok || res == -1 {
			errorMessage := fmt.Sprintf("need at least %v elements in array", column.minArrLen)
			return nil, errors.New(errorMessage)
		}
	}

	if column.maxArrLen != 0 {
		res, ok := compareTypeValues(len(interfaceArr), column.maxArrLen, "integer")
		if !ok || res == 1 {
			errorMessage := fmt.Sprintf("need at most %v elements in array", column.maxArrLen)
			return nil, errors.New(errorMessage)
		}
	}

	return interfaceArr, nil
}

// single value
func (column *Column) validateValueByMinMax(value any) error {
	datatype := strings.TrimSuffix(column.DataType, "[]")

	if column.minIndividual != nil {
		res, ok := compareTypeValues(value, column.minIndividual, datatype)
		if !ok || res == -1 {
			return errors.New("min constraint not satisfied")
		}
	}

	if column.maxIndividual != nil {
		res, ok := compareTypeValues(value, column.maxIndividual, datatype)
		if !ok || res == 1 {
			return errors.New("max constraint not satisfied")
		}
	}

	return nil
}

// single value
func (column *Column) validateValueByEnum(value any) error {
	if len(column.enumMap) == 0 {
		return nil
	}

	_, ok := column.enumMap[value]
	if !ok {
		errorMessage := fmt.Sprintf("value %v not present in enum", value)
		return errors.New(errorMessage)
	}

	return nil
}

/*
Doesn't convert value, just assertion, 1 if a > b, -1 if a < b; 0 otherwise.
false bool indicates failure.
strings are compared by length.
*/
func compareTypeValues(a, b any, datatype string) (int, bool) {
	switch datatype {
	case "text":
		parsedA, ok := a.(string)
		if !ok {
			return 0, false
		}

		parsedB, ok := b.(string)
		if !ok {
			return 0, false
		}

		if len(parsedA) > len(parsedB) {
			return 1, true
		} else if len(parsedA) < len(parsedB) {
			return -1, true
		} else {
			return 0, true
		}

	case "integer":
		parsedA, ok := a.(int)
		if !ok {
			return 0, false
		}

		parsedB, ok := b.(int)
		if !ok {
			return 0, false
		}

		if parsedA > parsedB {
			return 1, true
		}

		if parsedA < parsedB {
			return -1, true
		}

		return 0, true

	case "positiveInt":
		parsedA, ok := a.(uint64)
		if !ok {
			return 0, false
		}

		parsedB, ok := b.(uint64)
		if !ok {
			return 0, false
		}

		if parsedA > parsedB {
			return 1, true
		}

		if parsedA < parsedB {
			return -1, true
		}

		return 0, true

	case "real":
		parsedA, ok := a.(float64)
		if !ok {
			return 0, false
		}

		parsedB, ok := b.(float64)
		if !ok {
			return 0, false
		}

		if parsedA > parsedB {
			return 1, true
		} else if parsedA < parsedB {
			return -1, true
		} else {
			return 0, true
		}
	}

	if datatype == "date" || datatype == "time" || datatype == "timestamptz" {
		parsedA, ok := a.(time.Time)
		if !ok {
			return 0, false
		}

		parsedB, ok := b.(time.Time)
		if !ok {
			return 0, false
		}

		if time.Since(parsedA) > time.Since(parsedB) {
			return 1, true
		} else if time.Since(parsedA) < time.Since(parsedB) {
			return -1, true
		} else {
			return 0, true
		}
	}

	return 0, false
}

// For Normal Variables
func templateCheckConstraints(column Column, columnName string) string {
	args := []string{} // Min, Max, Enum

	if column.minIndividual != nil {
		formatted := templateValue(column.minIndividual, column.DataType)
		args = append(args, fmt.Sprintf("\"%v\" > %v", columnName, formatted))
	}

	if column.maxIndividual != nil {
		formatted := templateValue(column.maxIndividual, column.DataType)
		args = append(args, fmt.Sprintf("\"%v\" < %v", columnName, formatted))
	}

	if len(column.Enums) > 0 {
		formatted := templateValue(column.Enums, column.DataType+"[]")
		args = append(args, fmt.Sprintf("\"%v\" IN (%v)", columnName, formatted))
	}

	if len(args) == 0 {
		return ""
	}

	return fmt.Sprintf(" CHECK ( %v )", strings.Join(args, " AND "))
}

func templateValue(value any, datatype string) string {
	if datatype == "integer" || datatype == "real" || datatype == "bool" {
		return fmt.Sprintf("%v", value)
	}

	if datatype == "text" || datatype == "date" || datatype == "time" || datatype == "timestamptz" {
		return fmt.Sprintf("'%v'", value)
	}

	arr, ok := value.([]any)
	if !ok {
		return ""
	}

	datatype = strings.TrimSuffix(datatype, "[]")
	values := []string{}

	for _, item := range arr {
		formatted := fmt.Sprintf("'%v'", item)

		if datatype == "integer" || datatype == "real" || datatype == "bool" {
			formatted = fmt.Sprintf("%v", item)
		}

		values = append(values, formatted)
	}

	return "array[" + strings.Join(values, ", ") + "]::" + datatype + "[]"
}

func getArrayValidatorArgs(column Column) string {
	if !strings.HasSuffix(column.DataType, "[]") || (column.minArrLen == 0 &&
		column.maxArrLen == 0 &&
		column.minIndividual == nil &&
		column.maxIndividual == nil &&
		len(column.Enums) == 0) {
		return ""
	}

	res := []string{"false", "NULL", "NULL", "NULL", "NULL", "NULL"}
	if column.NotNull {
		res[0] = "true"
	}

	datatype := strings.TrimSuffix(column.DataType, "[]")

	if column.minArrLen > 0 {
		min_arr := templateValue(column.minArrLen, "integer")
		res = append(res, min_arr)
	}

	if column.maxArrLen > 0 {
		max_arr := templateValue(column.maxArrLen, "integer")
		res = append(res, max_arr)
	}

	if column.minIndividual != nil {
		min_ind := templateValue(column.minIndividual, datatype)
		res = append(res, min_ind)
	}

	if column.maxIndividual != nil {
		max_ind := templateValue(column.maxIndividual, datatype)
		res = append(res, max_ind)
	}

	if len(column.Enums) > 0 {
		enums := templateValue(column.Enums, datatype+"[]")
		res = append(res, enums)
	}

	return strings.Join(res, ", ")
}

func decrease(x int) int {
	return x - 1
}

func sanitize_db_label(text string) string {
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	trimmed := strings.TrimSpace(text)
	sanitized := regex.ReplaceAll([]byte(trimmed), []byte("_"))
	return string(sanitized)
}

func checkCSVExist(filePath, tableName string) error {
	fileName := filepath.Base(filePath)

	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("file %s for table %s doesn't exist", fileName, tableName)
		return errors.New(errorMessage)
	}

	if err != nil {
		errorMessage := fmt.Sprintf("failed to find file %s for table %s: %v", fileName, tableName, err)
		return errors.New(errorMessage)
	}

	return nil
}
