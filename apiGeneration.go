package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func (dbSchema *DB) writeAppFiles(appPath string) error {
	basePath, err := os.Getwd()
	if err != nil {
		return err
	}
	basePath = filepath.Join(basePath, "templates")

	errorChannel := make(chan error, 5)

	envApp, envTmpl := filepath.Join(appPath, ".env"), filepath.Join(basePath, "env.tmpl")
	go writeEnvFile(envApp, envTmpl, errorChannel)

	nullApp, nullTmpl := filepath.Join(appPath, "nullTypes.go"), filepath.Join(basePath, "nullTypes.tmpl")
	go writeNullTypes(nullApp, nullTmpl, errorChannel)

	modelApp, modelTmpl := filepath.Join(appPath, "models.go"), filepath.Join(basePath, "model.tmpl")
	go dbSchema.writeModels(modelApp, modelTmpl, errorChannel)

	dbApp, dbTmpl := filepath.Join(appPath, "dbUtils.go"), filepath.Join(basePath, "db.tmpl")
	go dbSchema.writeDbUtils(dbApp, dbTmpl, errorChannel)

	httpApp, httpTmpl := filepath.Join(appPath, "httpUtils.go"), filepath.Join(basePath, "http.tmpl")
	go dbSchema.writeHttpUtils(httpApp, httpTmpl, errorChannel)

	mainApp, mainTmpl := filepath.Join(appPath, "main.go"), filepath.Join(basePath, "main.tmpl")
	go dbSchema.writeMain(mainApp, mainTmpl, errorChannel)

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

func writeEnvFile(filePath, templatePath string, channel chan<- error) {
	var mainError error

	defer func() {
		if mainError != nil {
			errorMessage := fmt.Sprintf("error while writing .env file: %v", mainError)
			mainError = errors.New(errorMessage)
		}
		channel <- mainError
	}()

	template, err := template.New("env.tmpl").ParseFiles(templatePath)
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

func writeNullTypes(filePath, templatePath string, channel chan<- error) {
	var mainError error

	defer func() {
		if mainError != nil {
			errorMessage := fmt.Sprintf("error while writing nullTypes.go file: %v", mainError)
			mainError = errors.New(errorMessage)
		}
		channel <- mainError
	}()

	template, err := template.New("nullTypes.tmpl").ParseFiles(templatePath)
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

func (dbSchema *DB) writeModels(filePath, templatePath string, channel chan<- error) {
	var mainError error

	defer func() {
		if mainError != nil {
			errorMessage := fmt.Sprintf("error while writing models.go file: %v", mainError)
			mainError = errors.New(errorMessage)
		}
		channel <- mainError
	}()

	funcs := template.FuncMap{"getDbType": getDbType}

	template, err := template.New("model.tmpl").Funcs(funcs).ParseFiles(templatePath)
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

func (dbSchema *DB) writeDbUtils(filePath, templatePath string, channel chan<- error) {
	var mainError error

	defer func() {
		if mainError != nil {
			errorMessage := fmt.Sprintf("error while writing dbUtils.go file: %v", mainError)
			mainError = errors.New(errorMessage)
		}
		channel <- mainError
	}()

	funcs := template.FuncMap{
		"HasSuffix": strings.HasSuffix,
		"increase":  increase,
		"decrease":  decrease,
	}

	template, err := template.New("db.tmpl").Funcs(funcs).ParseFiles(templatePath)
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

func (dbSchema *DB) writeHttpUtils(filePath, templatePath string, channel chan<- error) {
	var mainError error

	defer func() {
		if mainError != nil {
			errorMessage := fmt.Sprintf("error while writing httpUtils.go file: %v", mainError)
			mainError = errors.New(errorMessage)
		}
		channel <- mainError
	}()

	template, err := template.New("http.tmpl").ParseFiles(templatePath)
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

func (dbSchema *DB) writeMain(filePath, templatePath string, channel chan<- error) {
	var mainError error

	defer func() {
		if mainError != nil {
			errorMessage := fmt.Sprintf("error while writing main.go file: %v", mainError)
			mainError = errors.New(errorMessage)
		}
		channel <- mainError
	}()

	template, err := template.New("main.tmpl").ParseFiles(templatePath)
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
