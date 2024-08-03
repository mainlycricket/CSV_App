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

func (dbSchema *DB) createStatements() (bytes.Buffer, error) {
	var createBuffer bytes.Buffer

	basePath, err := os.Getwd()
	if err != nil {
		return createBuffer, err
	}

	funcs := template.FuncMap{
		"HasSuffix":                strings.HasSuffix,
		"TrimSuffix":               strings.TrimSuffix,
		"templateValue":            templateValue,
		"decrease":                 decrease,
		"getArrayValidatorArgs":    getArrayValidatorArgs,
		"templateCheckConstraints": templateCheckConstraints,
	}

	fileName := "create.tmpl"
	templatePath := filepath.Join(basePath, fileName)

	template, err := template.New(fileName).Funcs(funcs).ParseFiles(templatePath)

	if err != nil {
		return createBuffer, err
	}

	buffer := bufio.NewWriter(&createBuffer)

	// TABLES
	if err := template.ExecuteTemplate(buffer, "Tables", dbSchema.Tables); err != nil {
		return createBuffer, err
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
				if err := template.ExecuteTemplate(buffer, "array_validator_function", datatype); err != nil {
					return createBuffer, err
				}
				datatypes[datatype] = true
			}
		}
	}

	// Table Validator Trigger Functions
	if err := template.ExecuteTemplate(buffer, "TableValidatorTrigger", dbSchema.Tables); err != nil {
		return createBuffer, err
	}

	if err := buffer.Flush(); err != nil {
		return createBuffer, err
	}

	return createBuffer, nil
}

func (dbSchema *DB) foreignKeyStatements() (bytes.Buffer, error) {
	var foreignBuffer bytes.Buffer

	basePath, err := os.Getwd()
	if err != nil {
		return foreignBuffer, err
	}

	fileName := "create.tmpl"

	funcs := template.FuncMap{
		"HasSuffix":                strings.HasSuffix,
		"TrimSuffix":               strings.TrimSuffix,
		"templateValue":            templateValue,
		"decrease":                 decrease,
		"getArrayValidatorArgs":    getArrayValidatorArgs,
		"templateCheckConstraints": templateCheckConstraints,
	}

	templatePath := filepath.Join(basePath, fileName)

	template, err := template.New(fileName).Funcs(funcs).ParseFiles(templatePath)

	if err != nil {
		return foreignBuffer, err
	}

	buffer := bufio.NewWriter(&foreignBuffer)

	if err := template.ExecuteTemplate(buffer, "ForeignKeys", dbSchema.Tables); err != nil {
		return foreignBuffer, err
	}

	if err := buffer.Flush(); err != nil {
		return foreignBuffer, err
	}

	return foreignBuffer, nil
}

func (dbSchema *DB) dataInsertion() (bytes.Buffer, error) {
	var insertionBuffer bytes.Buffer

	for tableName, table := range dbSchema.Tables {
		buffer := bufio.NewWriter(&insertionBuffer)

		text := fmt.Sprintf("-- DATA INSERTION \"%s\"\n", tableName)
		text += fmt.Sprintf(`INSERT INTO "%s"`, tableName)
		buffer.Write([]byte(text))

		filePath := filepath.Join(dbSchema.BasePath, table.FileName)

		if err := writeTableRows(filePath, &table, buffer); err != nil {
			return insertionBuffer, err
		}

		dbSchema.Tables[tableName] = table
	}

	if err := validateForeignValues(dbSchema.Tables); err != nil {
		errorMessage := fmt.Sprintf("error while validating foreign values: %v", err)
		return insertionBuffer, errors.New(errorMessage)
	}

	return insertionBuffer, nil
}

func writeTableRows(filePath string, table *Table, buffer *bufio.Writer) error {
	tableName := filepath.Base(filePath)

	fp, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		fp.Close()
	}()

	reader := csv.NewReader(fp)

	headers, err := reader.Read()
	if err != nil {
		return err
	}

	if _, err := buffer.Write([]byte("(")); err != nil {
		return err
	}

	for idx, header := range headers {
		arr := strings.SplitN(header, ":", 2)
		if len(arr) == 2 {
			header = arr[1]
		}

		columnName := sanitize_db_label(header)

		if _, ok := table.Columns[columnName]; !ok {
			errorMessage := fmt.Sprintf("%s column not found in %s table schema", header, tableName)
			return errors.New(errorMessage)
		}

		headers[idx] = columnName

		if idx == len(headers)-1 {
			if _, err := buffer.Write([]byte(fmt.Sprintf(`"%s"`, columnName))); err != nil {
				return err
			}

		} else {
			if _, err := buffer.Write([]byte(fmt.Sprintf(`"%s",`, columnName))); err != nil {
				return err
			}
		}
	}

	if _, err := buffer.Write([]byte(")\nVALUES\n")); err != nil {
		return err
	}

	rowIdx := 2

	for {
		row, err := reader.Read()
		if err == io.EOF {
			buffer.Write([]byte(";\n"))
			break
		}

		if rowIdx > 2 {
			buffer.Write([]byte(",\n"))
		}

		if err != nil {
			return err
		}

		if _, err := buffer.Write([]byte("(")); err != nil {
			return err
		}

		for idx, value := range row {
			columnName := headers[idx]
			column := table.Columns[columnName]
			val, err := column.validateValueByConstraints(value, true)

			if err != nil {
				errorMessage := fmt.Sprintf("error in row no. %d in %s column of %s table: %v", rowIdx, columnName, tableName, err)
				return errors.New(errorMessage)
			}

			str := templateValue(val, column.DataType)

			if len(column.ForeignField) > 0 && str != "NULL" {
				column.lookup[str] = rowIdx
			}

			if idx < len(headers)-1 {
				str += ", "
			}

			if _, err := buffer.Write([]byte(str)); err != nil {
				return err
			}

			table.Columns[columnName] = column
		}

		if _, err := buffer.Write([]byte(")")); err != nil {
			return err
		}
		rowIdx++
	}

	buffer.Write([]byte("\n"))
	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
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
