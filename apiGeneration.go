package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

type SlicedTableData struct {
	TableName  string
	PrimaryKey string
	Columns    []Column
}

type TemplateFnCall struct {
	filePath      string
	templatePath  string
	templateFuncs template.FuncMap
	data          any
}

func (dbSchema *DB) writeAppFiles(appPath string) error {
	basePath, err := os.Getwd()
	if err != nil {
		return err
	}
	basePath = filepath.Join(basePath, "templates")

	templatesData := dbSchema.getTemplatesData(basePath, appPath)

	FILES_COUNT := len(templatesData)

	errorChannel := make(chan error, FILES_COUNT)

	for _, item := range templatesData {
		go executeTemplate(item.filePath, item.templatePath, item.data, item.templateFuncs, errorChannel)
	}

	count := 0
	for err := range errorChannel {
		if err != nil {
			return err
		}
		count++
		if count == FILES_COUNT {
			close(errorChannel)
		}
	}

	return nil
}

func (dbSchema *DB) getTemplatesData(basePath, appPath string) []TemplateFnCall {
	slicedTableData := dbSchema.getSlicedTableData()

	templateData := []TemplateFnCall{
		// dbUtils
		{
			filePath:     filepath.Join(appPath, "dbUtils.go"),
			templatePath: filepath.Join(basePath, "db.tmpl"),
			data:         slicedTableData,
			templateFuncs: template.FuncMap{
				"HasSuffix":       strings.HasSuffix,
				"increase":        increase,
				"decrease":        decrease,
				"getTableColumns": getTableColumnsFn(slicedTableData),
			},
		},

		// models
		{
			filePath:      filepath.Join(appPath, "models.go"),
			templatePath:  filepath.Join(basePath, "model.tmpl"),
			templateFuncs: template.FuncMap{"getDbType": getDbType},
			data:          dbSchema.Tables,
		},

		// httpUtils
		{
			filePath:      filepath.Join(appPath, "httpUtils.go"),
			templatePath:  filepath.Join(basePath, "http.tmpl"),
			templateFuncs: template.FuncMap{"getPkType": getPkType},
			data:          dbSchema.Tables,
		},

		// .env
		{
			filePath:     filepath.Join(appPath, ".env"),
			templatePath: filepath.Join(basePath, "env.tmpl"),
		},

		// nullTypes
		{
			filePath:     filepath.Join(appPath, "nullTypes.go"),
			templatePath: filepath.Join(basePath, "nullTypes.tmpl"),
		},

		// utils
		{
			filePath:     filepath.Join(appPath, "utils.go"),
			templatePath: filepath.Join(basePath, "utils.tmpl"),
		},

		// main
		{
			filePath:     filepath.Join(appPath, "main.go"),
			templatePath: filepath.Join(basePath, "main.tmpl"),
		},
	}

	return templateData
}

func executeTemplate(filePath, templatePath string, templateData any, templateFuncs template.FuncMap, channel chan<- error) {
	var mainError error

	fileName := filepath.Base(filePath)
	templateName := filepath.Base(templatePath)

	defer func() {
		if mainError != nil {
			errorMessage := fmt.Sprintf("error while writing %s file: %v", fileName, mainError)
			mainError = errors.New(errorMessage)
		}
		channel <- mainError
	}()

	template, err := template.New(templateName).Funcs(templateFuncs).ParseFiles(templatePath)
	if err != nil {
		mainError = err
		return
	}

	fp, err := os.Create(filePath)
	if err != nil {
		mainError = err
		return
	}

	defer fp.Close()

	if err := template.Execute(fp, templateData); err != nil {
		mainError = err
		return
	}
}

func executeAppCommands(appPath string) error {
	if err := os.Chdir(appPath); err != nil {
		log.Fatalf("error while changing directory: %v", err)
	}

	commands := []string{
		"go fmt",
		"go mod init app.com/app",
		"go get github.com/lib/pq",
		"go mod tidy",
	}

	for _, command := range commands {
		arr := strings.Split(command, " ")
		cmd := exec.Command(arr[0], arr[1:]...)
		if err := cmd.Run(); err != nil {
			errorMessage := fmt.Sprintf("error while executing %s command: %v", command, err)
			return errors.New(errorMessage)
		}
	}

	return nil
}

func getDbType(datatype string) string {
	res := ""

	if strings.HasSuffix(datatype, "[]") {
		res += "[]"
		datatype = strings.TrimSuffix(datatype, "[]")
	}

	switch datatype {
	case "integer":
		res += "CustomNullInt"
	case "real":
		res += "CustomNullFloat"
	case "text":
		res += "CustomNullString"
	case "boolean":
		res += "CustomNullBool"
	case "date":
		res += "CustomNullDate"
	case "time":
		res += "CustomNullTime"
	case "timestamptz":
		res += "CustomNullDateTime"
	}

	return res
}

func getPkType(table Table) string {
	if table.PrimaryKey == "" {
		return "CustomNullInt"
	}

	return getDbType(table.Columns[table.PrimaryKey].DataType)
}

// Returns table along with its columns (ordered by names) in slice, instead of maps
func (dbSchema *DB) getSlicedTableData() []SlicedTableData {
	tablesData := []SlicedTableData{}

	for _, table := range dbSchema.Tables {
		item := SlicedTableData{
			TableName:  table.TableName,
			PrimaryKey: table.PrimaryKey,
			Columns:    []Column{},
		}

		for _, column := range table.Columns {
			item.Columns = append(item.Columns, column)
		}

		slices.SortFunc(item.Columns, func(col1, col2 Column) int {
			if col1.ColumnName > col2.ColumnName {
				return 1
			}

			if col1.ColumnName < col2.ColumnName {
				return -1
			}

			return 0
		})

		tablesData = append(tablesData, item)
	}

	return tablesData
}

// Returns a function to get column associated with a table
func getTableColumnsFn(tablesData []SlicedTableData) func(tableName string) []Column {
	return func(tableName string) []Column {
		for _, table := range tablesData {
			if table.TableName == tableName {
				return table.Columns
			}
		}
		return []Column{}
	}
}
