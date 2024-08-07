package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	argsMessage := "Invalid args! Provide:\n'schema' to generate schema.json or \n'sql' to generate db.sql"

	if len(args) != 2 {
		log.Fatalf(argsMessage)
	}

	input := strings.ToLower(strings.TrimSpace(args[1]))

	if input == "schema" {
		err := generateInititalSchema()
		if err != nil {
			log.Fatalf("Failed to generate initial schema: %v", err)
		}
		fmt.Println("Generated Schema successfully!")
	} else if input == "sql" {
		dbSchema, err := readSchema()

		if err != nil {
			log.Fatalf("Failed to parse DB schema: %v", err)
		}

		err = dbSchema.validateSchema()

		if err != nil {
			log.Fatalf("Schema Validation Failed: %v", err)
		}

		insertionBuffer, err := dbSchema.dataInsertion()
		if err != nil {
			log.Fatalf("error while data insertion: %v", err)
		}

		createBuffer, err := dbSchema.createStatements()
		if err != nil {
			log.Fatalf("error while creating sql statements: %v", err)
		}

		foreignBuffer, err := dbSchema.foreignKeyStatements()
		if err != nil {
			log.Fatalf("error while adding foreign key constriants: %v", err)
		}

		basePath, err := os.Getwd()
		if err != nil {
			log.Fatalf("error while reading current directory: %v", err)
		}

		filePath := filepath.Join(basePath, "data", "db.sql")

		if err := writeFile(filePath, createBuffer, insertionBuffer, foreignBuffer); err != nil {
			log.Fatalf("error while creating db.sql: %v", err)
		}

		fmt.Println("db.sql generated")
	} else if input == "app" {
		dbSchema, err := readSchema()

		if err != nil {
			log.Fatalf("Failed to parse DB schema: %v", err)
		}

		err = dbSchema.validateSchema()

		if err != nil {
			log.Fatalf("Schema Validation Failed: %v", err)
		}

		basePath, err := os.Getwd()
		if err != nil {
			log.Fatalf("error while reading current directory path: %v", err)
		}

		appPath := filepath.Join(basePath, "app")

		if err := os.Mkdir(appPath, os.ModePerm); err != nil {
			log.Fatalf("error while creating app: %v", err)
		}

		modelFilePath := filepath.Join(appPath, "models.go")
		if err := dbSchema.writeModels(modelFilePath); err != nil {
			log.Fatalf("error while writing models: %v", err)
		}

		dbFilePath := filepath.Join(appPath, "dbUtils.go")
		if err := dbSchema.writeDbUtils(dbFilePath); err != nil {
			log.Fatalf("error while DB Utils: %v", err)
		}

		httpFilePath := filepath.Join(appPath, "httpUtils.go")
		if err := dbSchema.writeHttpUtils(httpFilePath); err != nil {
			log.Fatalf("error while writing http utils: %v", err)
		}

		mainFilePath := filepath.Join(appPath, "main.go")
		if err := dbSchema.writeMain(mainFilePath); err != nil {
			log.Fatalf("error while writing main: %v", err)
		}

		if err := executeAppCommands(appPath); err != nil {
			log.Fatalf("error while executing commands: %v", err)
		}

		fmt.Println("app generated")
	} else {
		log.Fatalf(argsMessage)
	}
}
