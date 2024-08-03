package main

import (
	"fmt"
	"log"
)

func main() {
	input := "sql"

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

		fmt.Println("Schema Validated!")

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

		if err := writeFile("db.sql", createBuffer, insertionBuffer, foreignBuffer); err != nil {
			log.Fatalf("error while creating db.sql: %v", err)
		}
	}
}
