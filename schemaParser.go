package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type tableResponse struct {
	table Table
	err   error
}

func generateInititalSchema() error {
	BasePath, err := os.Getwd()
	if err != nil {
		return err
	}

	dataPath := filepath.Join(BasePath, "data")
	dirList, err := os.ReadDir(dataPath)
	if err != nil {
		return err
	}

	tableRespChannel := make(chan tableResponse, 5)
	var tablesCount int
	var mutex sync.Mutex
	primaryKeys := make(map[string]string, 5)

	csvFiles := make(map[string]bool, 5)

	appConfig := AppCongif{
		SchemaPath:   filepath.Join(dataPath, "schema.json"),
		TablesConfig: make(map[string]TableAuthConfig, 5),
	}

	for _, file := range dirList {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}

		fileName := file.Name()
		tableName := sanitize_db_label(strings.TrimSuffix(fileName, ".csv"))

		if len(tableName) == 0 {
			return errors.New("an unnamed csv file is found")
		}

		if _, ok := csvFiles[tableName]; ok {
			message := fmt.Sprintf("table %s already exists", tableName)
			return errors.New(message)
		}

		appConfig.TablesConfig[tableName] = TableAuthConfig{}
		csvFiles[tableName] = true
		mutex.Lock()
		tablesCount += 1
		mutex.Unlock()

		filePath := filepath.Join(dataPath, fileName)
		go createTableSchema(filePath, tableRespChannel)
	}

	dbSchema := DB{BasePath: dataPath, Tables: make(map[string]Table, 5)}

	// receive table schemas
	for resp := range tableRespChannel {
		if resp.err != nil {
			return resp.err
		} else {
			fileName := resp.table.FileName
			tableName := sanitize_db_label(strings.TrimSuffix(fileName, ".csv"))
			dbSchema.Tables[tableName] = resp.table
			key := resp.table.PrimaryKey + ":" + resp.table.Columns[resp.table.PrimaryKey].DataType
			primaryKeys[key] = tableName
		}

		mutex.Lock()
		tablesCount -= 1
		mutex.Unlock()

		if tablesCount == 0 {
			close(tableRespChannel)
		}
	}

	dbSchema.setForeignKeys(primaryKeys)

	if err := writeJsonFile(appConfig.SchemaPath, dbSchema); err != nil {
		return err
	}

	err = writeJsonFile(filepath.Join(dataPath, "appConfig.json"), appConfig)
	return err
}

func (dbSchema *DB) setForeignKeys(primaryKeys map[string]string) {
	for tableName, table := range dbSchema.Tables {
		for columnName, column := range table.Columns {
			if column.ForeignField == "__" {
				referencedTable, ok := primaryKeys[columnName+":"+column.DataType]
				if ok {
					column.ForeignTable = referencedTable
					column.ForeignField = columnName
					table.Columns[columnName] = column
				}
			}
		}
		dbSchema.Tables[tableName] = table
	}
}

// Parses a CSV file and writes the response to channel
func createTableSchema(filePath string, tableResponseChannel chan<- tableResponse) {
	fileName := filepath.Base(filePath)
	tableName := sanitize_db_label(strings.TrimSuffix(fileName, ".csv"))

	table := Table{FileName: fileName, TableName: tableName}
	table.Columns = make(map[string]Column, 20)
	var mainError error

	fp, err := os.Open(filePath)

	defer func() {
		fp.Close()
		tableResponseChannel <- tableResponse{table: table, err: mainError}
	}()

	if err != nil {
		mainError = err
		return
	}

	reader := csv.NewReader(fp)
	headers, err := reader.Read()

	if err == io.EOF {
		mainError = errors.New("no headers found in " + fileName)
		return
	}

	if err != nil {
		mainError = err
		return
	}

	// Initialize Columns
	for idx, header := range headers {
		column := Column{}

		columnName := column.setTableConstraints(&table, header)

		if len(columnName) == 0 {
			message := fmt.Sprintf("empty column in %s table", fileName)
			mainError = errors.New(message)
			return
		}

		column.ColumnName = columnName

		if _, ok := table.Columns[columnName]; ok {
			message := fmt.Sprintf("column %s already exists in %s table", columnName, fileName)
			mainError = errors.New(message)
			return
		}

		headers[idx] = columnName
		table.Columns[columnName] = column
	}

	// Detect DataTypes by traversing rows
	err = setColumnTypes(reader, &table, headers)
	if err != nil {
		message := fmt.Sprintf("error while parsing %s table data: %v", fileName, err)
		mainError = errors.New(message)
		return
	}
}

// It reads the entire CSV file and sets the column types
func setColumnTypes(reader *csv.Reader, table *Table, headers []string) error {
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		count := len(headers)

		// Traversing columns in the row
		for j, value := range row {
			columnName := headers[j]
			column := table.Columns[columnName]

			existingType := column.DataType

			value = strings.TrimSpace(value)

			if len(value) == 0 {
				continue
			}

			if existingType == "text" {
				count--
				continue
			}

			// validate against existing type
			if existingType != "" {
				ok := validateAgainstExistingType(value, existingType)

				if ok {
					continue
				}

				detectedType := detectDataType(value)

				if (existingType == "integer" && detectedType == "real") ||
					(existingType == "real" && detectedType == "integer") {
					column.DataType = "real"
				} else {
					column.DataType = "text"
				}
			} else {
				column.DataType = detectDataType(value)
			}

			table.Columns[columnName] = column
			if column.DataType == "text" {
				count--
			}
		}

		// stop reading if all columns are texts
		if count == 0 {
			break
		}
	}

	return nil
}

// Sets table constraints and returns the sanitized column name
func (column *Column) setTableConstraints(table *Table, columnName string) string {
	columnName = strings.TrimSpace(columnName)
	arr := strings.SplitN(columnName, ":", 2)

	if len(arr) == 1 {
		return sanitize_db_label(arr[0])
	}

	constraint := strings.TrimSpace(arr[0])
	constraint = strings.ToUpper(constraint)
	columnName = sanitize_db_label(arr[1])

	if strings.ContainsRune(constraint, 'P') {
		table.PrimaryKey = columnName
		column.Unique = true
		column.NotNull = true
		return columnName
	}

	for _, char := range constraint {
		if char == 'U' {
			column.Unique = true
		}

		if char == 'N' {
			column.NotNull = true
		}

		if char == 'F' {
			column.ForeignField = "__"
			column.ForeignTable = "__"
		}
	}

	return columnName
}

// checks all types including no type
func detectDataType(value string) string {
	if basicType := detectBasicDataType(value); len(basicType) > 0 {
		return basicType
	}

	// try array
	arr := []any{}
	err := json.Unmarshal([]byte(value), &arr)

	if err != nil {
		return "text"
	}

	if len(arr) == 0 {
		return ""
	}

	firstVal := fmt.Sprintf("%v", arr[0])
	prevType := detectBasicDataType(firstVal)
	if len(prevType) == 0 || prevType == "text" {
		return "text[]"
	}

	// Traversing array column
	for _, subVal := range arr {
		subVal := fmt.Sprintf("%v", subVal)
		detectedType := detectBasicDataType(subVal)

		if (prevType == "integer" && detectedType == "real") ||
			(prevType == "real" && detectedType == "integer") {
			prevType = "real"
			continue
		}

		if detectedType != prevType {
			return "text[]"
		} else {
			prevType = detectedType
		}
	}

	return prevType + "[]"
}

func detectBasicDataType(value string) string {
	_, err := strconv.Atoi(value)
	if err == nil {
		return "integer"
	}

	_, err = strconv.ParseFloat(value, 64)
	if err == nil {
		return "real"
	}

	_, err = strconv.ParseBool(value)
	if err == nil {
		return "boolean"
	}

	_, err = time.Parse(datetimeFormats["date"], value)
	if err == nil {
		return "date"
	}

	_, err = time.Parse(datetimeFormats["time"], value)
	if err == nil {
		return "time"
	}

	_, err = time.Parse(datetimeFormats["timestamptz"], value)
	if err == nil {
		return "timestamptz"
	}

	if !(strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]")) {
		return "text"
	}

	return ""
}

func validateAgainstExistingType(value, datatype string) bool {
	if strings.HasSuffix(datatype, "[]") {
		arr := []any{}
		err := json.Unmarshal([]byte(value), &arr)

		if err != nil {
			return false
		}

		if len(arr) == 0 {
			return true
		}

		datatype := strings.TrimSuffix(datatype, "[]")
		for _, value := range arr {
			_, ok := validateValueByType(value, datatype)
			if !ok {
				return false
			}
		}
		return true
	}

	_, ok := validateValueByType(value, datatype)
	return ok
}

func (dbSchema *DB) validateSchema() error {
	basePath := dbSchema.BasePath

	for tableName, table := range dbSchema.Tables {
		if tableName != sanitize_db_label(tableName) {
			errorMessage := fmt.Sprintf("table name %s isn't sanitized", tableName)
			return errors.New(errorMessage)
		}

		filePath := filepath.Join(basePath, table.FileName)

		if err := checkCSVExist(filePath, tableName); err != nil {
			return err
		}

		primaryKeyFlag := false

		for columnName, column := range table.Columns {
			if len(columnName) == 0 {
				errorMessage := fmt.Sprintf("empty column found in table %s", tableName)
				return errors.New(errorMessage)
			}

			if columnName != sanitize_db_label(columnName) {
				errorMessage := fmt.Sprintf("column %s in table %s isn't sanitized", columnName, tableName)
				return errors.New(errorMessage)
			}

			// Data Type
			if validType := isValidTypeName(column.DataType); !validType {
				errorMessage := fmt.Sprintf("invalid type for column %s in table %s", columnName, tableName)
				return errors.New(errorMessage)
			}

			// Primary & Foreign Key
			if columnName == table.PrimaryKey {
				primaryKeyFlag = true
			}

			if len(column.ForeignField) > 0 || len(column.ForeignTable) > 0 {
				if column.ForeignTable == tableName {
					errorMessage := fmt.Sprintf("column %s in table %s refers to same table", columnName, tableName)
					return errors.New(errorMessage)
				}

				referencedTable, ok := dbSchema.Tables[column.ForeignTable]

				if !ok {
					errorMessage := fmt.Sprintf("invalid referenced table by %s column in %s table", columnName, tableName)
					return errors.New(errorMessage)
				}

				referredCol, ok := referencedTable.Columns[column.ForeignField]
				if !ok {
					errorMessage := fmt.Sprintf("invalid referenced column by %s column in %s table", columnName, tableName)
					return errors.New(errorMessage)
				}

				if referencedTable.PrimaryKey != columnName {
					errorMessage := fmt.Sprintf("referenced column by %s column in %s table isn't primary key", columnName, tableName)
					return errors.New(errorMessage)
				}

				if referredCol.DataType != column.DataType {
					errorMessage := fmt.Sprintf("referenced column by %s column in %s table isn't of %s datatype", columnName, tableName, column.DataType)
					return errors.New(errorMessage)
				}

				column.lookup = make(map[string]int)
			}

			// Set Min, Max Constraints
			if err := column.setMinMaxConstraint(); err != nil {
				errorMessage := fmt.Sprintf("invalid min/max constraint for column %s in table %s:\n:%v", columnName, tableName, err)
				return errors.New(errorMessage)
			}

			// Validate Enum Types
			if err := column.validateEnums(); err != nil {
				errorMessage := fmt.Sprintf("invalid enum for column %s in table %s:\n%v", columnName, tableName, err)
				return errors.New(errorMessage)
			}

			// Default Value
			if err := column.validateDefaultValue(); err != nil {
				errorMessage := fmt.Sprintf("invalid default value for column %s in table %s:\n%v", columnName, tableName, err)
				return errors.New(errorMessage)
			}

			// unqiue
			if column.Unique {
				if column.Default != nil {
					errorMessage := fmt.Sprintf("unique column %s has a default value in table %s", columnName, tableName)
					return errors.New(errorMessage)
				}
				column.values = make(map[string]bool)
			}

			table.Columns[columnName] = column
		}

		// Primary key
		if len(table.PrimaryKey) > 0 && !primaryKeyFlag {
			errorMessage := fmt.Sprintf("invalid primary key %s in table %s", table.PrimaryKey, tableName)
			return errors.New(errorMessage)
		}

		dbSchema.Tables[tableName] = table
	}

	return nil
}
