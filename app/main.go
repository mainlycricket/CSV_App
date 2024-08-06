package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = connectDB()

	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}

	server := startServer()
	server.ListenAndServe()
}
