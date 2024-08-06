package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func (dbSchema *DB) writeModels(filePath string) error {
	basePath, err := os.Getwd()
	if err != nil {
		return err
	}

	templateFileName := "model.tmpl"

	funcs := template.FuncMap{
		"capitalizeStr": capitalizeStr,
		"getGoType":     getGoType,
		"HasPrefix":     strings.HasPrefix,
	}

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	fp, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer fp.Close()

	if err := template.Execute(fp, dbSchema.Tables); err != nil {
		return err
	}

	return nil
}

func (dbSchema *DB) writeMain(filePath string) error {
	basePath, err := os.Getwd()
	if err != nil {
		return err
	}

	templateFileName := "main.tmpl"

	templatePath := filepath.Join(basePath, "templates", templateFileName)

	template, err := template.New(templateFileName).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	fp, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer fp.Close()

	if err := template.Execute(fp, nil); err != nil {
		return err
	}

	return nil
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

func getGoType(datatype string) string {
	res := ""

	if strings.HasSuffix(datatype, "[]") {
		res += "[]"
		datatype = strings.TrimSuffix(datatype, "[]")
	}

	switch datatype {
	case "integer":
		res += "int"
	case "real":
		res += "float64"
	case "text":
		res += "string"
	case "boolean":
		res += "bool"
	case "date":
		res += "time.Time"
	case "time":
		res += "time.Time"
	case "timestamptz":
		res += "time.Time"
	}

	return res
}

func capitalizeStr(text string) string {
	res := ""

	if len(text) > 0 {
		res += strings.ToUpper(string(text[0]))
	}

	if len(text) > 1 {
		res += text[1:]
	}

	return res
}
