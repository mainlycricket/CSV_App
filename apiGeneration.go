package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

type TemplateTableData struct {
	TableName   string
	PrimaryKey  string
	Columns     []Column
	IsAuthTable bool
	TableConfig TableConfig
}

type TemplateFnCall struct {
	filePath      string
	templatePath  string
	templateFuncs template.FuncMap
	data          any
}

func (dbSchema *DB) writeAppFiles(appPath string, appConfig *AppCongif) error {
	basePath, err := os.Getwd()
	if err != nil {
		return err
	}
	basePath = filepath.Join(basePath, "templates")

	templatesData := dbSchema.getTemplatesMetaData(basePath, appPath, appConfig)

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

func (dbSchema *DB) getTemplatesMetaData(basePath, appPath string, appConfig *AppCongif) []TemplateFnCall {
	slicedTableData := dbSchema.getSlicedTableData(appConfig)

	getOrgFields := func() []string { return appConfig.OrgFields }

	templateData := []TemplateFnCall{
		// dbUtils
		{
			filePath:     filepath.Join(appPath, "dbUtils.go"),
			templatePath: filepath.Join(basePath, "db.tmpl"),
			data:         slicedTableData,
			templateFuncs: template.FuncMap{
				"HasSuffix":         strings.HasSuffix,
				"increase":          increase,
				"decrease":          decrease,
				"getTableColumns":   getTableColumnsFn(slicedTableData),
				"getOrgFields":      getOrgFields,
				"capitalize":        capitalize,
				"sliceContains":     slices.Contains[[]string, string],
				"getColumnDataType": getDataTypeFn(slicedTableData),
			},
		},

		// models
		{
			filePath:     filepath.Join(appPath, "models.go"),
			templatePath: filepath.Join(basePath, "model.tmpl"),
			templateFuncs: template.FuncMap{
				"getDbType":         getDbType,
				"getOrgFields":      getOrgFields,
				"capitalize":        capitalize,
				"sliceContains":     slices.Contains[[]string, string],
				"getColumnDataType": getDataTypeFn(slicedTableData),
			},
			data: slicedTableData,
		},

		// httpUtils
		{
			filePath:     filepath.Join(appPath, "httpUtils.go"),
			templatePath: filepath.Join(basePath, "http.tmpl"),
			templateFuncs: template.FuncMap{
				"getPkType":    getPkType,
				"HasSuffix":    strings.HasSuffix,
				"increase":     increase,
				"decrease":     decrease,
				"capitalize":   capitalize,
				"getOrgFields": getOrgFields,
			},
			data: slicedTableData,
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
			filePath:      filepath.Join(appPath, "utils.go"),
			templatePath:  filepath.Join(basePath, "utils.tmpl"),
			templateFuncs: template.FuncMap{"getOrgFields": getOrgFields},
			data:          slicedTableData,
		},

		// main
		{
			filePath:     filepath.Join(appPath, "main.go"),
			templatePath: filepath.Join(basePath, "main.tmpl"),
		},

		// setup.sh
		{
			filePath:     filepath.Join(appPath, "setup.sh"),
			templatePath: filepath.Join(basePath, "setup.tmpl"),
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

	if strings.HasSuffix(filePath, ".sh") {
		os.Chmod(filePath, 0o755)
	}

	if err := template.Execute(fp, templateData); err != nil {
		mainError = err
		return
	}
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

func getPkType(table TemplateTableData) string {
	if table.PrimaryKey == "" {
		return "CustomNullInt"
	}

	for _, column := range table.Columns {
		if column.ColumnName == table.PrimaryKey {
			return getDbType(column.DataType)
		}
	}

	return ""
}

// Returns table along with its columns (ordered by names) in slice, instead of maps
func (dbSchema *DB) getSlicedTableData(appConfig *AppCongif) []TemplateTableData {
	tablesData := []TemplateTableData{}

	for _, table := range dbSchema.Tables {
		item := TemplateTableData{
			TableName:   table.TableName,
			PrimaryKey:  table.PrimaryKey,
			Columns:     []Column{},
			IsAuthTable: table.TableName == appConfig.AuthTable,
			TableConfig: appConfig.Tables[table.TableName],
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
func getTableColumnsFn(tablesData []TemplateTableData) func(tableName string) []Column {
	return func(tableName string) []Column {
		for _, table := range tablesData {
			if table.TableName == tableName {
				return table.Columns
			}
		}
		return []Column{}
	}
}

// Returns a function to get datatype of a column in table name
func getDataTypeFn(tablesData []TemplateTableData) func(tableName, columnName string) string {
	return func(tableName, columnName string) string {
		for _, table := range tablesData {
			if table.TableName == tableName {
				for _, column := range table.Columns {
					if column.ColumnName == columnName {
						return getDbType(column.DataType)
					}
				}
			}
		}
		return ""
	}
}
