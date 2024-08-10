package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func (dbSchema *DB) writeAppFiles(appPath string) error {
	errorChannel := make(chan error, 5)

	go writeEnvFile(filepath.Join(appPath, ".env"), errorChannel)
	go writeNullTypes(filepath.Join(appPath, "nullTypes.go"), errorChannel)
	go dbSchema.writeModels(filepath.Join(appPath, "models.go"), errorChannel)
	go dbSchema.writeDbUtils(filepath.Join(appPath, "dbUtils.go"), errorChannel)
	go dbSchema.writeHttpUtils(filepath.Join(appPath, "httpUtils.go"), errorChannel)
	go dbSchema.writeMain(filepath.Join(appPath, "main.go"), errorChannel)

	count := 0
	for err := range errorChannel {
		if err != nil {
			return err
		}
		count++
		if count == 6 {
			close(errorChannel)
		}
	}

	return nil
}

func writeEnvFile(filePath string, channel chan<- error) {
	var mainError error

	defer func() {
		channel <- mainError
	}()

	basePath, err := os.Getwd()
	if err != nil {
		mainError = err
		return
	}

	templateFileName := "env.tmpl"

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).ParseFiles(templatePath)
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

	if err := template.Execute(fp, nil); err != nil {
		mainError = err
		return
	}
}

func writeNullTypes(filePath string, channel chan<- error) {
	var mainError error

	defer func() {
		channel <- mainError
	}()

	basePath, err := os.Getwd()
	if err != nil {
		mainError = err
		return
	}

	templateFileName := "nullTypes.tmpl"

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).ParseFiles(templatePath)
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

	if err := template.Execute(fp, nil); err != nil {
		mainError = err
		return
	}
}

func (dbSchema *DB) writeModels(filePath string, channel chan<- error) {
	var mainError error

	defer func() {
		channel <- mainError
	}()

	basePath, err := os.Getwd()
	if err != nil {
		mainError = err
		return
	}

	templateFileName := "model.tmpl"

	funcs := template.FuncMap{"getDbType": getDbType}

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).Funcs(funcs).ParseFiles(templatePath)
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

	if err := template.Execute(fp, dbSchema.Tables); err != nil {
		mainError = err
		return
	}
}

func (dbSchema *DB) writeDbUtils(filePath string, channel chan<- error) {
	var mainError error

	defer func() {
		channel <- mainError
	}()

	basePath, err := os.Getwd()
	if err != nil {
		mainError = err
		return
	}

	templateFileName := "db.tmpl"

	funcs := template.FuncMap{
		"HasSuffix": strings.HasSuffix,
		"increase":  increase,
		"decrease":  decrease,
	}

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).Funcs(funcs).ParseFiles(templatePath)
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

	if err := template.Execute(fp, dbSchema.Tables); err != nil {
		mainError = err
		return
	}
}

func (dbSchema *DB) writeHttpUtils(filePath string, channel chan<- error) {
	var mainError error

	defer func() {
		channel <- mainError
	}()

	basePath, err := os.Getwd()
	if err != nil {
		mainError = err
		return
	}

	templateFileName := "http.tmpl"

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).ParseFiles(templatePath)
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

	if err := template.Execute(fp, dbSchema.Tables); err != nil {
		mainError = err
		return
	}
}

func (dbSchema *DB) writeMain(filePath string, channel chan<- error) {
	var mainError error

	defer func() {
		channel <- mainError
	}()

	basePath, err := os.Getwd()
	if err != nil {
		mainError = err
		return
	}

	templateFileName := "main.tmpl"

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).ParseFiles(templatePath)
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

	if err := template.Execute(fp, nil); err != nil {
		mainError = err
		return
	}
}

func executeAppCommands(appPath string) error {
	if err := os.Chdir(appPath); err != nil {
		log.Fatalf("error while changing directory: %v", err)
	}

	fmtCmd := exec.Command("go", "fmt")
	if err := fmtCmd.Run(); err != nil {
		return err
	}

	initCommand := exec.Command("go", "mod", "init", "app.com/app")
	if err := initCommand.Run(); err != nil {
		return err
	}

	libCmd := exec.Command("go", "get", "github.com/lib/pq")
	if err := libCmd.Run(); err != nil {
		return err
	}

	tidyCmd := exec.Command("go", "mod", "tidy")
	if err := tidyCmd.Run(); err != nil {
		return err
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
