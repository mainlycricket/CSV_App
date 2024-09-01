package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type insertionResponse struct {
	err   error
	table *Table
}

func (dbSchema *DB) createStatements() (*bytes.Buffer, error) {
	var createBuffer bytes.Buffer

	basePath, err := os.Getwd()
	if err != nil {
		return &createBuffer, err
	}

	funcs := template.FuncMap{
		"HasSuffix":                strings.HasSuffix,
		"TrimSuffix":               strings.TrimSuffix,
		"templateValue":            templateValue,
		"decrease":                 decrease,
		"getArrayValidatorArgs":    getArrayValidatorArgs,
		"templateCheckConstraints": templateCheckConstraints,
	}

	fileName := "sql.tmpl"
	templatePath := filepath.Join(basePath, "templates", fileName)

	template, err := template.New(fileName).Funcs(funcs).ParseFiles(templatePath)

	if err != nil {
		return &createBuffer, err
	}

	writer := bufio.NewWriter(&createBuffer)

	// TABLES
	if err := template.ExecuteTemplate(writer, "Tables", dbSchema.Tables); err != nil {
		return &createBuffer, err
	}

	// CREATE ARRAY VALIDATORS
	datatypes := map[string]bool{}
	for _, table := range dbSchema.Tables {
		for _, column := range table.Columns {
			datatype := column.DataType
			isArray := strings.HasSuffix(datatype, "[]")
			datatype = strings.TrimSuffix(datatype, "[]")

			if isArray && !datatypes[datatype] && (column.minArrLen > 0 ||
				column.maxArrLen > 0 ||
				column.minIndividual != nil ||
				column.maxIndividual != nil ||
				len(column.Enums) > 0) {
				if err := template.ExecuteTemplate(writer, "array_validator_function", datatype); err != nil {
					return &createBuffer, err
				}
				datatypes[datatype] = true
			}
		}
	}

	// Table Validator Trigger Functions
	if err := template.ExecuteTemplate(writer, "TableValidatorTrigger", dbSchema.Tables); err != nil {
		return &createBuffer, err
	}

	if err := writer.Flush(); err != nil {
		return &createBuffer, err
	}

	return &createBuffer, nil
}

func (dbSchema *DB) foreignKeyStatements() (*bytes.Buffer, error) {
	var foreignBuffer bytes.Buffer

	basePath, err := os.Getwd()
	if err != nil {
		return &foreignBuffer, err
	}

	fileName := "sql.tmpl"

	funcs := template.FuncMap{
		"HasSuffix":                strings.HasSuffix,
		"TrimSuffix":               strings.TrimSuffix,
		"templateValue":            templateValue,
		"decrease":                 decrease,
		"getArrayValidatorArgs":    getArrayValidatorArgs,
		"templateCheckConstraints": templateCheckConstraints,
	}

	templatePath := filepath.Join(basePath, "templates", fileName)

	template, err := template.New(fileName).Funcs(funcs).ParseFiles(templatePath)

	if err != nil {
		return &foreignBuffer, err
	}

	writer := bufio.NewWriter(&foreignBuffer)

	if err := template.ExecuteTemplate(writer, "ForeignKeys", dbSchema.Tables); err != nil {
		return &foreignBuffer, err
	}

	if err := writer.Flush(); err != nil {
		return &foreignBuffer, err
	}

	return &foreignBuffer, nil
}

func (dbSchema *DB) dataInsertion() (*bytes.Buffer, error) {
	var insertionBuffer bytes.Buffer

	responseChannel := make(chan insertionResponse, 4)
	tableCount := len(dbSchema.Tables)

	for _, table := range dbSchema.Tables {
		writer := bufio.NewWriter(&insertionBuffer)
		filePath := filepath.Join(dbSchema.BasePath, table.FileName)
		go writeTableRows(filePath, &table, writer, responseChannel)
	}

	for response := range responseChannel {
		tableCount--
		if response.err != nil {
			return &insertionBuffer, response.err
		}
		table := response.table
		dbSchema.Tables[table.TableName] = *table
		if tableCount == 0 {
			close(responseChannel)
		}
	}

	if err := validateForeignValues(dbSchema.Tables); err != nil {
		errorMessage := fmt.Sprintf("error while validating foreign values: %v", err)
		return &insertionBuffer, errors.New(errorMessage)
	}

	return &insertionBuffer, nil
}

func writeTableRows(filePath string, table *Table, writer *bufio.Writer, channel chan<- insertionResponse) {
	tableName := table.TableName
	var mainError error

	fp, err := os.Open(filePath)
	if err != nil {
		channel <- insertionResponse{table: table, err: err}
		return
	}

	defer func() {
		fp.Close()
		channel <- insertionResponse{table: table, err: mainError}
	}()

	reader := csv.NewReader(fp)

	headersSQL := ""
	headers, err := reader.Read()
	if err != nil {
		mainError = err
		return
	}

	for idx, header := range headers {
		arr := strings.SplitN(header, ":", 2)
		if len(arr) == 2 {
			header = arr[1]
		}

		columnName := sanitize_db_label(header)

		if _, ok := table.Columns[columnName]; !ok {
			errorMessage := fmt.Sprintf("%s column not found in %s table schema", header, tableName)
			mainError = errors.New(errorMessage)
			return
		}

		headers[idx] = columnName

		headersSQL += fmt.Sprintf(`"%s"`, columnName)
		if idx < len(headers)-1 {
			headersSQL += ", "
		}
	}

	rowIdx := 2

	for {
		row, err := reader.Read()
		if err == io.EOF {
			if rowIdx > 2 {
				writer.Write([]byte(";\n"))
			}
			break
		}

		if rowIdx == 2 {
			text := fmt.Sprintf("-- DATA INSERTION \"%s\"\n", tableName)
			text += fmt.Sprintf("INSERT INTO \"%s\" (%s)\nVALUES\n", tableName, headersSQL)

			if _, err := writer.Write([]byte(text)); err != nil {
				mainError = err
				return
			}
		}

		if rowIdx > 2 {
			writer.Write([]byte(",\n"))
		}

		if err != nil {
			mainError = err
			return
		}

		if _, err := writer.Write([]byte("(")); err != nil {
			mainError = err
			return
		}

		for idx, value := range row {
			columnName := headers[idx]
			column := table.Columns[columnName]
			val, err := column.validateValueByConstraints(value, true)

			if err != nil {
				errorMessage := fmt.Sprintf("error in row no. %d in %s column of %s table: %v", rowIdx, columnName, tableName, err)
				mainError = errors.New(errorMessage)
				return
			}

			if column.Hash && val != nil {
				hashedVal, err := hashText(val, column.DataType)

				if err != nil {
					errorMessage := fmt.Sprintf("error in row no. %d in %s column of %s table: %v", rowIdx, columnName, tableName, err)
					mainError = errors.New(errorMessage)
					return
				}

				val = hashedVal
			}

			str := templateValue(val, column.DataType)

			if len(column.ForeignField) > 0 && str != "NULL" {
				column.lookup[str] = rowIdx
			}

			if idx < len(headers)-1 {
				str += ", "
			}

			if _, err := writer.Write([]byte(str)); err != nil {
				mainError = err
				return
			}

			table.Columns[columnName] = column
		}

		if _, err := writer.Write([]byte(")")); err != nil {
			mainError = err
			return
		}
		rowIdx++
	}

	writer.Write([]byte("\n"))
	if err := writer.Flush(); err != nil {
		mainError = err
		return
	}
}

func validateForeignValues(tables map[string]Table) error {
	for tableName, table := range tables {
		for columnName, column := range table.Columns {
			if len(column.lookup) > 0 {
				foreignTable := tables[column.ForeignTable]
				foreignColumn := foreignTable.Columns[column.ForeignField]

				for key, rowNum := range column.lookup {
					if !foreignColumn.values[key] {
						errorMessage := fmt.Sprintf("invalid value %s for foreign key column %s in %s table on row num: %d", key, columnName, tableName, rowNum)
						return errors.New(errorMessage)
					}
				}
			}
		}
	}
	return nil
}
