package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"
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
				"HasSuffix":                strings.HasSuffix,
				"increase":                 increase,
				"decrease":                 decrease,
				"getTableColumns":          getTableColumnsFn(slicedTableData),
				"getOrgFields":             getOrgFields,
				"capitalize":               capitalize,
				"sliceContains":            slices.Contains[[]string, string],
				"getColumnDataType":        getDataTypeFn(slicedTableData),
				"getProtectedValuesByRole": getProtectedValuesByRole,
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
				"getPkType":          getPkType,
				"HasSuffix":          strings.HasSuffix,
				"TrimPrefix":         strings.TrimPrefix,
				"increase":           increase,
				"decrease":           decrease,
				"capitalize":         capitalize,
				"getOrgFields":       getOrgFields,
				"getDbType":          getDbType,
				"templateProtectMap": templateProtectMap,
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

func (appConfig *AppCongif) validateAppConfig(dbSchema *DB) error {
	if appConfig.SchemaPath != filepath.Join(dbSchema.BasePath, "schema.json") {
		return fmt.Errorf("invalid schema path: %v", appConfig.SchemaPath)
	}

	var authTable Table
	var rolesEnum []string

	// Auth Table Validation
	if len(appConfig.AuthTable) > 0 {
		var exists bool

		authTable, exists = dbSchema.Tables[appConfig.AuthTable]
		if !exists {
			return fmt.Errorf("auth table %s doesn't exist in schema", appConfig.AuthTable)
		}

		if authTable.PrimaryKey != "username" {
			return fmt.Errorf("auth table %s doesn't have 'username' primary key", authTable.TableName)
		}

		usernameColumn := authTable.Columns["username"]
		if usernameColumn.DataType != "text" || usernameColumn.Hash {
			return errors.New(`"username" field should be a text field with hash disabled`)
		}

		passwordColumn, ok := authTable.Columns["password"]
		if !ok {
			return fmt.Errorf("auth table %s doesn't have 'password' field", authTable.TableName)
		}
		if passwordColumn.DataType != "text" || !passwordColumn.Hash || !passwordColumn.NotNull {
			return errors.New(`"password" field should be a not null text field with hash enabled`)
		}

		roleField, ok := authTable.Columns["role"]
		if ok && (roleField.DataType != "text" || roleField.Hash || !roleField.NotNull || len(roleField.Enums) == 0) {
			return errors.New(`"role" field should be a not null text field with enums enabled and hash disabled`)
		}

		for _, orgField := range appConfig.OrgFields {
			orgColumn, ok := authTable.Columns[orgField]
			if !ok {
				return fmt.Errorf(`"%s" org field not found in auth table`, orgField)
			}
			if orgColumn.DataType != "text" || orgColumn.Hash {
				return fmt.Errorf(`org field "%s" should be a field with hash disabled`, orgField)
			}
		}
	}

	if authTable.TableName != "" {
		rolesEnum, _ = assertAnyArr[string](authTable.Columns["role"].Enums)
	}

	// Tables Validaton
	appTablesCount, schemaTablesCount := len(appConfig.Tables), len(dbSchema.Tables)
	if appTablesCount != schemaTablesCount {
		return fmt.Errorf("app config contains %d tables while schema contains %d tables", appTablesCount, schemaTablesCount)
	}

	for tableName, tableConfig := range appConfig.Tables {
		table, ok := dbSchema.Tables[tableName]
		if !ok {
			return fmt.Errorf(`"%s" table not found in schema`, tableName)
		}

		if err := tableConfig.ReadAllConfig.validateReadConfig(dbSchema, tableName); err != nil {
			return fmt.Errorf(`invalid readAllConfig in "%s" table: %w`, tableName, err)
		}

		if err := tableConfig.ReadByPkConfig.validateReadConfig(dbSchema, tableName); err != nil {
			return fmt.Errorf(`invalid readByPkConfig in "%s" table: %w`, tableName, err)
		}

		if err := tableConfig.ReadAllAuth.validateAuthInfo(rolesEnum, appConfig.OrgFields, table.Columns, authTable, tableName == authTable.TableName); err != nil {
			return fmt.Errorf(`invalid readAllAuth for %s table: %w`, tableName, err)
		}

		if err := tableConfig.ReadByPkAuth.validateAuthInfo(rolesEnum, appConfig.OrgFields, table.Columns, authTable, tableName == authTable.TableName); err != nil {
			return fmt.Errorf(`invalid readByPkAuth for %s table: %w`, tableName, err)
		}

		if err := tableConfig.InsertAuth.validateAuthInfo(rolesEnum, appConfig.OrgFields, table.Columns, authTable, tableName == authTable.TableName); err != nil {
			return fmt.Errorf(`invalid insertAuth for %s table: %w`, tableName, err)
		}

		if err := tableConfig.UpdateAuth.validateAuthInfo(rolesEnum, appConfig.OrgFields, table.Columns, authTable, tableName == authTable.TableName); err != nil {
			return fmt.Errorf(`invalid updateAuth for %s table: %w`, tableName, err)
		}

		if err := tableConfig.DeleteAuth.validateAuthInfo(rolesEnum, appConfig.OrgFields, table.Columns, authTable, tableName == authTable.TableName); err != nil {
			return fmt.Errorf(`invalid deleteAuth for %s table: %w`, tableName, err)
		}

		if tableConfig.DefaultPagination == 0 {
			return fmt.Errorf(`invalid default pagination for "%s" table`, tableName)
		}
	}

	return nil
}

func (readConfig *ReadConfig) validateReadConfig(dbSchema *DB, tableName string) error {
	for _, field := range readConfig.Columns {
		column, exists := dbSchema.Tables[tableName].Columns[field]

		if !exists {
			return fmt.Errorf(`"%s" field not found in table columns`, field)
		}

		if column.Hash {
			return fmt.Errorf(`"%s" field has hash enabled`, field)
		}
	}

	for tableField, foreignFields := range readConfig.ForeignColumns {
		if !slices.Contains(readConfig.Columns, tableField) {
			return fmt.Errorf(`"%s" foreign field not found in selected fields`, tableField)
		}

		column := dbSchema.Tables[tableName].Columns[tableField]
		if column.Hash {
			return fmt.Errorf(`"%s" foreign field has hash enabled`, tableField)
		}

		foreignTable := dbSchema.Tables[column.ForeignField]
		for _, foreignField := range foreignFields {
			if _, exists := foreignTable.Columns[foreignField]; exists {
				return fmt.Errorf(`"%s" foreign field not found in "%s" referenced table schema`, foreignField, foreignTable.TableName)
			}
		}
	}

	return nil
}

func (authInfo *AuthInfo) validateAuthInfo(rolesEnum, appOrgFields []string, columnsMap map[string]Column, authTable Table, isAuthTable bool) error {
	if userField := authInfo.UserField; len(userField) > 0 {
		if authTable.TableName == "" {
			return fmt.Errorf(`user field "%s" exists without authTable`, userField)
		}

		userFieldCol, ok := columnsMap[userField]
		if !ok {
			return fmt.Errorf(`user field "%s" not found in table schema`, userField)
		}

		flag := false

		if isAuthTable && userFieldCol.ColumnName == authTable.PrimaryKey {
			flag = true
		}

		if !flag && (userFieldCol.ForeignTable != authTable.TableName || userFieldCol.ForeignField != authTable.PrimaryKey) {
			return fmt.Errorf(`user field "%s" in table schema is not referencing "username" in auth table`, userField)
		}
	}

	if len(authInfo.OrgFields) > 0 && len(appOrgFields) == 0 {
		return fmt.Errorf(`orgFields exist without application level orgFields`)
	}

	for tableField, authField := range authInfo.OrgFields {
		if _, exists := columnsMap[tableField]; !exists {
			return fmt.Errorf(`org field "%s" not found in table schema`, tableField)
		}

		if !slices.Contains(appOrgFields, authField) {
			return fmt.Errorf(`org field "%s" not found in appConfig orgFields`, tableField)
		}
	}

	if !authInfo.BasicAuth && len(authInfo.AllowedRoles) > 0 {
		return errors.New("basic auth should be true to enable role based authorization")
	}

	for _, role := range authInfo.AllowedRoles {
		if !slices.Contains(rolesEnum, role) {
			return fmt.Errorf(`invalid allowedRoles: "%s" role not present in roles enum of auth table`, role)
		}
	}

	for field, explicitSetters := range authInfo.Privileges {
		if field != authInfo.UserField && authInfo.OrgFields[field] == "" {
			return fmt.Errorf(`invalid priviliges: "%s" field, neither an userField nor an orgField`, field)
		}

		for _, explicitSetter := range explicitSetters {
			if explicitSetter == "" {
				if authInfo.BasicAuth {
					return fmt.Errorf(`invalid priviliges: "%s" field allows non-logged-in users to be explicit setters, but basic auth is enabled`, field)
				}
				continue
			}

			if len(authInfo.AllowedRoles) > 0 && !slices.Contains(authInfo.AllowedRoles, explicitSetter) {
				return fmt.Errorf(`invalid priviliges: explicitSetter "%s" for "%s" field isn't found in allowedRoles`, explicitSetter, field)
			}

			if len(authInfo.AllowedRoles) == 0 && !slices.Contains(rolesEnum, explicitSetter) {
				return fmt.Errorf(`invalid priviliges: explicitSetter "%s" for "%s" field isn't found in login roles enum`, explicitSetter, field)
			}
		}
	}

	for field, valueSettersMap := range authInfo.ProtectedFields {
		column, ok := columnsMap[field]
		if !ok {
			return fmt.Errorf(`invalid protectedFields: "%s" column is not present in table schema`, field)
		}

		if len(column.Enums) == 0 && !strings.HasPrefix(column.DataType, "boolean") {
			return fmt.Errorf(`invalid protectedFields: "%s" column should be either boolean/boolean[] or have Enums enabled`, field)
		}

		if err := validateProtectMap(valueSettersMap, column.DataType, column.Enums); err != nil {
			return fmt.Errorf(`invalid protectedFields: "%s" column has invalid values %w`, field, err)
		}

		for _, explicitSetters := range valueSettersMap {
			for _, explicitSetter := range explicitSetters {
				if explicitSetter == "" {
					if authInfo.BasicAuth {
						return fmt.Errorf(`invalid protectedFields: "%s" field allows non-logged-in users to be explicit setters, but basic auth is enabled`, field)
					}

					continue
				}

				if len(authInfo.AllowedRoles) > 0 && !slices.Contains(authInfo.AllowedRoles, explicitSetter) {
					return fmt.Errorf(`invalid protectedFields: explicitSetter "%s" for "%s" field isn't found in allowedRoles`, explicitSetter, field)
				}

				if len(authInfo.AllowedRoles) == 0 && !slices.Contains(rolesEnum, explicitSetter) {
					return fmt.Errorf(`invalid protectedFields: explicitSetter "%s" for "%s" field isn't found in login roles enum`, explicitSetter, field)
				}
			}
		}
	}

	return nil
}

func templateProtectMap(m map[string][]string, datatype string) string {
	datatype = strings.TrimSuffix(datatype, "[]")

	switch datatype {
	case "integer":
		res, _ := assertProtectedMap[int64](m, datatype)
		return fmt.Sprintf("%#v", res)
	case "real":
		res, _ := assertProtectedMap[float64](m, datatype)
		return fmt.Sprintf("%#v", res)
	case "text":
		res, _ := assertProtectedMap[string](m, datatype)
		return fmt.Sprintf("%#v", res)
	case "boolean":
		res, _ := assertProtectedMap[bool](m, datatype)
		return fmt.Sprintf("%#v", res)
	case "date":
		res, _ := assertProtectedMap[time.Time](m, datatype)
		return fmt.Sprintf("%#v", res)
	case "time":
		res, _ := assertProtectedMap[time.Time](m, datatype)
		return fmt.Sprintf("%#v", res)
	case "timestamptz":
		res, _ := assertProtectedMap[time.Time](m, datatype)
		return fmt.Sprintf("%#v", res)
	}

	return ""
}

func validateProtectMap(m map[string][]string, datatype string, enums []any) error {
	datatype = strings.TrimSuffix(datatype, "[]")

	switch datatype {
	case "integer":
		castedMap, err := assertProtectedMap[int64](m, datatype)
		if err != nil {
			return err
		}

		castedEnums, _ := assertAnyArr[int64](enums)
		for castedValue := range castedMap {
			if !slices.Contains(castedEnums, castedValue) {
				return fmt.Errorf(`%v value not found in column enums`, castedValue)
			}
		}
	case "real":
		castedMap, err := assertProtectedMap[float64](m, datatype)
		if err != nil {
			return err
		}

		castedEnums, _ := assertAnyArr[float64](enums)
		for castedValue := range castedMap {
			if !slices.Contains(castedEnums, castedValue) {
				return fmt.Errorf(`%v value not found in column enums`, castedValue)
			}
		}
	case "text":
		castedMap, err := assertProtectedMap[string](m, datatype)
		if err != nil {
			return err
		}

		castedEnums, _ := assertAnyArr[string](enums)
		for castedValue := range castedMap {
			if !slices.Contains(castedEnums, castedValue) {
				return fmt.Errorf(`%v value not found in column enums`, castedValue)
			}
		}
	case "boolean":
		if len(enums) == 0 {
			enums = []any{true, false}
		}

		castedMap, err := assertProtectedMap[bool](m, datatype)
		if err != nil {
			return err
		}

		castedEnums, _ := assertAnyArr[bool](enums)
		for castedValue := range castedMap {
			if !slices.Contains(castedEnums, castedValue) {
				return fmt.Errorf(`%v value not found in column enums`, castedValue)
			}
		}
	case "date":
		castedMap, err := assertProtectedMap[time.Time](m, datatype)
		if err != nil {
			return err
		}

		castedEnums, _ := assertAnyArr[time.Time](enums)
		for castedValue := range castedMap {
			if !slices.Contains(castedEnums, castedValue) {
				return fmt.Errorf(`%v value not found in column enums`, castedValue)
			}
		}
	case "time":
		castedMap, err := assertProtectedMap[time.Time](m, datatype)
		if err != nil {
			return err
		}

		castedEnums, _ := assertAnyArr[time.Time](enums)
		for castedValue := range castedMap {
			if !slices.Contains(castedEnums, castedValue) {
				return fmt.Errorf(`%v value not found in column enums`, castedValue)
			}
		}
	case "timestamptz":
		castedMap, err := assertProtectedMap[time.Time](m, datatype)
		if err != nil {
			return err
		}

		castedEnums, _ := assertAnyArr[time.Time](enums)
		for castedValue := range castedMap {
			if !slices.Contains(castedEnums, castedValue) {
				return fmt.Errorf(`%v value not found in column enums`, castedValue)
			}
		}
	default:
		return fmt.Errorf(`unsupported datatype: %s`, datatype)
	}

	return nil
}

func assertProtectedMap[T comparable](m map[string][]string, datatype string) (map[T][]string, error) {
	res := make(map[T][]string, len(m))

	for strVal, roles := range m {
		interfaceVal, err := typeConversionFuncs[datatype](strVal)

		if err != nil {
			return nil, fmt.Errorf(`failed to type cast value %s`, strVal)
		}

		casted, ok := interfaceVal.(T)
		if !ok {
			return nil, fmt.Errorf(`failed to type assert value %s`, strVal)
		}

		res[casted] = roles
	}

	return res, nil
}

// returns negative map[role][]values i.e. values which can't be set by a particular role
func getProtectedValuesByRole(m map[string][]string, datatype string) map[string][]any {
	data := make(map[string][]any, 5)
	var roles []string

	// get all roles
	for _, allowedRoles := range m {
		for _, role := range allowedRoles {
			if !slices.Contains(roles, role) {
				roles = append(roles, role)
			}
		}
	}

	datatype = strings.TrimPrefix(datatype, "[]")

	for value, allowedRoles := range m {
		interfaceVal, _ := typeConversionFuncs[datatype](value)

		for _, role := range roles {
			if !slices.Contains(allowedRoles, role) {
				data[role] = append(data[role], interfaceVal)
			}
		}
	}

	return data
}
