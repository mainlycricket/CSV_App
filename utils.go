package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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

// Enums are validated against datatype, Min, Max constraints.
// Individual elements are validated for array types
func (column *Column) validateEnums() error {
	datatype := strings.TrimSuffix(column.DataType, "[]")

	if len(column.Enums) > 25 {
		return errors.New("array contains more than 25 values")
	}

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
	}

	return nil
}

// Returns true is Default value is nil or it satisfies the Min, Max & Enums constraints
func (column *Column) validateDefaultValue() error {
	if column.Default == nil {
		return nil
	}

	interfaceVal, err := column.validateValueByConstraints(column.Default, false)
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

// check if the given value (including arrays) satisfies all the constraints
// if insert flag is true, NOT NULL & Unique constraints are also validated
// also returns the interface value of the provided value
func (column *Column) validateValueByConstraints(value any, insert bool) (any, error) {
	if insert {
		str := strings.TrimSpace(fmt.Sprintf("%v", value))

		if len(str) == 0 || value == nil {
			if column.NotNull {
				return nil, errors.New("null values aren't allowed")
			}

			return column.Default, nil
		}

		if column.Unique {
			if column.values[str] {
				return nil, errors.New("unique constraint not satisfied")
			}
			templateVal := templateValue(str, column.DataType)
			column.values[templateVal] = true
		}
	}

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

	interfaceVal, ok := validateValueByType(value, column.DataType)
	if !ok {
		errorMessage := fmt.Sprintf("%v should be of %v type", value, column.DataType)
		return nil, errors.New(errorMessage)
	}

	if err := column.validateValueByMinMax(interfaceVal); err != nil {
		return nil, err
	}

	if err := column.validateValueByEnum(interfaceVal); err != nil {
		return nil, err
	}

	return interfaceVal, nil
}

// receives string or []any array and checks the array min & max length constraints
// also typecastes the array value into an array interface
func (column *Column) validateValArrLen(value any) ([]any, error) {
	text, ok := value.(string)
	if ok {
		if err := json.Unmarshal([]byte(text), &value); err != nil {
			return nil, err
		}
	}

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

// checks if the provided value (non-array) satisfies the min, max constraints
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

// checks if the provided value (non-array) is present in column enums
func (column *Column) validateValueByEnum(value any) error {
	if len(column.Enums) == 0 {
		return nil
	}

	if !slices.Contains(column.Enums, value) {
		errorMessage := fmt.Sprintf("value %v not present in enum", value)
		return errors.New(errorMessage)
	}

	return nil
}

/*
returns 1 if a > b, -1 if a < b and 0 otherwise, doesn't convert values before conversion
false bool indicates failure.
strings are compared by length.
not usable for array or boolean datatypes
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

		// more time passed indicates older datetime (smaller datetime)
		if time.Since(parsedA) > time.Since(parsedB) {
			return -1, true
		} else if time.Since(parsedA) < time.Since(parsedB) {
			return 1, true
		} else {
			return 0, true
		}
	}

	return 0, false
}

// SQL Check Constraints templating for non-array variables
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
		var formattedValues []string
		for _, val := range column.Enums {
			formattedValue := templateValue(val, column.DataType)
			formattedValues = append(formattedValues, formattedValue)
		}
		args = append(args, fmt.Sprintf(`"%s" IN (%s)`, columnName, strings.Join(formattedValues, ",")))
	}

	if len(args) == 0 {
		return ""
	}

	return fmt.Sprintf(" CHECK ( %v )", strings.Join(args, " AND "))
}

// get corresponding SQL value for all datatypes including array ones
func templateValue(value any, datatype string) string {
	if value == nil {
		return "NULL"
	} else if parsed := fmt.Sprintf("%v", value); len(parsed) == 0 {
		return "NULL"
	}

	if datatype == "integer" || datatype == "real" || datatype == "boolean" {
		return fmt.Sprintf("%v", value)
	}

	if datatype == "date" || datatype == "time" || datatype == "timestamptz" {
		parsed, ok := value.(time.Time)
		if !ok {
			return "NULL"
		}
		return fmt.Sprintf("'%v'", parsed.Format(datetimeFormats[datatype]))
	}

	if datatype == "text" {
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

		if datatype == "date" || datatype == "time" || datatype == "timestamptz" {
			parsed, ok := item.(time.Time)
			if !ok {
				formatted = "NULL"
			}
			formatted = fmt.Sprintf("'%v'", parsed.Format(datetimeFormats[datatype]))
		}
		if datatype == "integer" || datatype == "real" || datatype == "boolean" {
			formatted = fmt.Sprintf("%v", item)
		}

		values = append(values, formatted)
	}

	return "array[" + strings.Join(values, ", ") + "]::" + datatype + "[]"
}

// used in SQL trigger generation
func getArrayValidatorArgs(column Column) string {
	if !strings.HasSuffix(column.DataType, "[]") || (column.minArrLen == 0 &&
		column.maxArrLen == 0 &&
		column.minIndividual == nil &&
		column.maxIndividual == nil &&
		len(column.Enums) == 0) {
		return ""
	}

	// NOT NULL, Min Arr Len, Max Arr Len, Min Individual, Max Individual, Enums
	res := []string{"false", "NULL", "NULL", "NULL", "NULL", "NULL"}
	if column.NotNull {
		res[0] = "true"
	}

	datatype := strings.TrimSuffix(column.DataType, "[]")

	if column.minArrLen > 0 {
		res[1] = templateValue(column.minArrLen, "integer")
	}

	if column.maxArrLen > 0 {
		res[2] = templateValue(column.maxArrLen, "integer")
	}

	if column.minIndividual != nil {
		res[3] = templateValue(column.minIndividual, datatype)
	}

	if column.maxIndividual != nil {
		res[4] = templateValue(column.maxIndividual, datatype)
	}

	if len(column.Enums) > 0 {
		res[5] = templateValue(column.Enums, datatype+"[]")
	}

	return strings.Join(res, ", ")
}

func decrease(x int) int {
	return x - 1
}

func increase(x int) int {
	return x + 1
}

func sanitize_db_label(text string) string {
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	trimmed := strings.TrimSpace(text)
	sanitized := regex.ReplaceAll([]byte(trimmed), []byte("_"))
	return string(sanitized)
}

// used in schema validation
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

// creates a file and writes the content of provided buffers one after another
func writeFile(filePath string, buffers ...*bytes.Buffer) error {
	fp, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		fp.Close()
	}()

	for _, buffer := range buffers {
		if _, err := fp.Write(buffer.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func hashText(val any, datatype string) (any, error) {
	if datatype == "text" {
		str, ok := val.(string)

		if !ok {
			return nil, errors.New("failed to typecast to text")
		}

		hashedVal, err := hashPassword(str)
		if err != nil {
			return nil, fmt.Errorf("error while hashing %s: %v", str, err)
		}

		return hashedVal, nil
	}

	if datatype == "text[]" {
		arr, ok := val.([]any)
		if !ok {
			return nil, errors.New("failed to typecast to text[]")
		}

		for idx, item := range arr {
			str, ok := item.(string)

			if !ok {
				return nil, fmt.Errorf("failed to typecast item no. %d to text", (idx + 1))
			}

			hashedVal, err := hashPassword(str)
			if err != nil {
				return nil, fmt.Errorf("error while hashing %s: %v", str, err)
			}

			arr[idx] = hashedVal
		}

		return arr, nil
	}

	return nil, errors.New("invalid datatype")
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func writeJsonFile(filepath string, data any) error {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, jsonData, os.ModePerm)

	return err
}

func readJsonFile(filePath string, ptr any) error {
	schema, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(schema, ptr)

	if err != nil {
		return err
	}

	return nil
}

func capitalize(text string) string {
	if len(text) == 0 {
		return text
	}

	return strings.ToUpper(string(text[0])) + text[1:]
}

func convertAnyArrToStrArr(array []any) ([]string, error) {
	res := make([]string, 0, len(array))

	for _, item := range array {
		strVal, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf(`failed to typecast %v into string`, item)
		}
		res = append(res, strVal)
	}

	return res, nil
}
