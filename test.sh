psql -h localhost -U postgres -c 'DROP DATABASE "CSV_App"'
rm -r ./app

#!/bin/bash
set -e

go test
go build .
./CSV_App schema 
./CSV_App sql
psql -h localhost -U postgres -c 'CREATE DATABASE "CSV_App"'
psql -h localhost -U postgres -d "CSV_App" -f data/db.sql
./CSV_App app
