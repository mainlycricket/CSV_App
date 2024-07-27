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
	}
}
