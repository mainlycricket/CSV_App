package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lib/pq"
)

func readEnvFile() error {
	basePath, err := os.Getwd()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join(basePath, ".env"))
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		arr := strings.SplitN(line, "=", 2)
		if len(arr) != 2 {
			continue
		}

		key := strings.TrimSpace(arr[0])
		value := strings.TrimSpace(arr[1])
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return nil
}

func connectDB() (*sql.DB, error) {
	if err := readEnvFile(); err != nil {
		errorMessage := fmt.Sprintf("error while reading .env file: %v", err)
		return nil, errors.New(errorMessage)
	}

	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// students CRUD
func db_insert_students(ctx context.Context, item *Table_students) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "students" ("Branch_Id", "Course_Id", "Student_Father", "Student_Id", "Student_Name") VALUES ($1, $2, $3, $4, $5)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Branch_Id, item.Column_Course_Id, item.Column_Student_Father, item.Column_Student_Id, item.Column_Student_Name)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_students(ctx context.Context, clause string, args []any) ([]Table_students_response, error) {
	data := []Table_students_response{}

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "branches"."Course_Id", "branches"."Teachers", "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed", "students"."Student_Father", "students"."Student_Id", "students"."Student_Name" FROM "students" INNER JOIN "branches" ON "students"."Branch_Id" = "branches"."Branch_Id" INNER JOIN "courses" ON "students"."Course_Id" = "courses"."Course_Id"`

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_students_response{}

		rows.Scan(&item.Fkey_Branch_Id.Column_Branch_Id, &item.Fkey_Branch_Id.Column_Branch_Name, &item.Fkey_Branch_Id.Column_Course_Id, pq.Array(&item.Fkey_Branch_Id.Column_Teachers), &item.Fkey_Course_Id.Column_Course_Id, &item.Fkey_Course_Id.Column_Course_Name, &item.Fkey_Course_Id.Column_Lateral_Allowed, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name)

		data = append(data, item)
	}

	return data, nil
}

func db_read_students_ByPK(ctx context.Context, id string) (Table_students_response, error) {
	item := Table_students_response{}

	args := []any{id}

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "branches"."Course_Id", "branches"."Teachers", "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed", "students"."Student_Father", "students"."Student_Id", "students"."Student_Name" FROM "students" INNER JOIN "branches" ON "students"."Branch_Id" = "branches"."Branch_Id" INNER JOIN "courses" ON "students"."Course_Id" = "courses"."Course_Id" WHERE "students"."Student_Id" = $1`

	userVal := ctx.Value("userField")
	if userVal != nil {
		value := userVal.(string)
		args = append(args, value)
		query += fmt.Sprintf(` AND "students"."userField" = $%d`, len(args))
	}

	orgVal := ctx.Value("orgField")
	if orgVal != nil {
		value := orgVal.(string)
		args = append(args, value)
		query += fmt.Sprintf(` AND "students"."orgField" = $%d`, len(args))
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Fkey_Branch_Id.Column_Branch_Id, &item.Fkey_Branch_Id.Column_Branch_Name, &item.Fkey_Branch_Id.Column_Course_Id, pq.Array(&item.Fkey_Branch_Id.Column_Teachers), &item.Fkey_Course_Id.Column_Course_Id, &item.Fkey_Course_Id.Column_Course_Name, &item.Fkey_Course_Id.Column_Lateral_Allowed, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_students(ctx context.Context, id string, item *Table_students) error {
	args := []any{item.Column_Branch_Id, item.Column_Course_Id, item.Column_Student_Father, item.Column_Student_Id, item.Column_Student_Name, id}

	query := `UPDATE "students" SET "Branch_Id" = $1, "Course_Id" = $2, "Student_Father" = $3, "Student_Id" = $4, "Student_Name" = $5 WHERE "Student_Id" = $6`

	userVal := ctx.Value("userField")
	if userVal != nil {
		value := userVal.(string)
		args = append(args, value)
		query += fmt.Sprintf(` AND "students"."userField" = $%d`, len(args))
	}

	orgVal := ctx.Value("orgField")
	if orgVal != nil {
		value := orgVal.(string)
		args = append(args, value)
		query += fmt.Sprintf(` AND "students"."orgField" = $%d`, len(args))
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_students(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "students" WHERE "Student_Id" = $1`

	userVal := ctx.Value("userField")
	if userVal != nil {
		value := userVal.(string)
		args = append(args, value)
		query += fmt.Sprintf(` AND "students"."userField" = $%d`, len(args))
	}

	orgVal := ctx.Value("orgField")
	if orgVal != nil {
		value := orgVal.(string)
		args = append(args, value)
		query += fmt.Sprintf(` AND "students"."orgField" = $%d`, len(args))
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

// subjects CRUD
func db_insert_subjects(ctx context.Context, item *Table_subjects) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "subjects" ("Branch_Id", "Subject_Id", "Subject_Name") VALUES ($1, $2, $3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Branch_Id, item.Column_Subject_Id, item.Column_Subject_Name)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_subjects(ctx context.Context, clause string, args []any) ([]Table_subjects_response, error) {
	data := []Table_subjects_response{}

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "branches"."Course_Id", "branches"."Teachers", "subjects"."Subject_Id", "subjects"."Subject_Name" FROM "subjects" INNER JOIN "branches" ON "subjects"."Branch_Id" = "branches"."Branch_Id"`

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_subjects_response{}

		rows.Scan(&item.Fkey_Branch_Id.Column_Branch_Id, &item.Fkey_Branch_Id.Column_Branch_Name, &item.Fkey_Branch_Id.Column_Course_Id, pq.Array(&item.Fkey_Branch_Id.Column_Teachers), &item.Column_Subject_Id, &item.Column_Subject_Name)

		data = append(data, item)
	}

	return data, nil
}

func db_read_subjects_ByPK(ctx context.Context, id string) (Table_subjects_response, error) {
	item := Table_subjects_response{}

	args := []any{id}

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "branches"."Course_Id", "branches"."Teachers", "subjects"."Subject_Id", "subjects"."Subject_Name" FROM "subjects" INNER JOIN "branches" ON "subjects"."Branch_Id" = "branches"."Branch_Id" WHERE "subjects"."Subject_Id" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Fkey_Branch_Id.Column_Branch_Id, &item.Fkey_Branch_Id.Column_Branch_Name, &item.Fkey_Branch_Id.Column_Course_Id, pq.Array(&item.Fkey_Branch_Id.Column_Teachers), &item.Column_Subject_Id, &item.Column_Subject_Name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_subjects(ctx context.Context, id string, item *Table_subjects) error {
	args := []any{item.Column_Branch_Id, item.Column_Subject_Id, item.Column_Subject_Name, id}

	query := `UPDATE "subjects" SET "Branch_Id" = $1, "Subject_Id" = $2, "Subject_Name" = $3 WHERE "Subject_Id" = $4`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_subjects(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "subjects" WHERE "Subject_Id" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

// TypeTest CRUD
func db_insert_TypeTest(ctx context.Context, item *Table_TypeTest) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "TypeTest" ("Bool", "Bool_Arr", "Date", "DateTime", "Date_arr", "Datetime_Arr", "Float", "Float_arr", "Int", "Int_Arr", "Str_Arr", "String", "Time", "Time_Arr") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Bool, pq.Array(item.Column_Bool_Arr), item.Column_Date, item.Column_DateTime, pq.Array(item.Column_Date_arr), pq.Array(item.Column_Datetime_Arr), item.Column_Float, pq.Array(item.Column_Float_arr), item.Column_Int, pq.Array(item.Column_Int_Arr), pq.Array(item.Column_Str_Arr), item.Column_String, item.Column_Time, pq.Array(item.Column_Time_Arr))

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_TypeTest(ctx context.Context, clause string, args []any) ([]Table_TypeTest_response, error) {
	data := []Table_TypeTest_response{}

	query := `SELECT "TypeTest"."__ID", "TypeTest"."Bool", "TypeTest"."Bool_Arr", "TypeTest"."Date", "TypeTest"."DateTime", "TypeTest"."Date_arr", "TypeTest"."Datetime_Arr", "TypeTest"."Float", "TypeTest"."Float_arr", "TypeTest"."Int", "TypeTest"."Int_Arr", "TypeTest"."Str_Arr", "TypeTest"."String", "TypeTest"."Time", "TypeTest"."Time_Arr" FROM "TypeTest"`

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_TypeTest_response{}

		rows.Scan(&item.ID__, &item.Column_Bool, pq.Array(&item.Column_Bool_Arr), &item.Column_Date, &item.Column_DateTime, pq.Array(&item.Column_Date_arr), pq.Array(&item.Column_Datetime_Arr), &item.Column_Float, pq.Array(&item.Column_Float_arr), &item.Column_Int, pq.Array(&item.Column_Int_Arr), pq.Array(&item.Column_Str_Arr), &item.Column_String, &item.Column_Time, pq.Array(&item.Column_Time_Arr))

		data = append(data, item)
	}

	return data, nil
}

func db_read_TypeTest_ByPK(ctx context.Context, id string) (Table_TypeTest_response, error) {
	item := Table_TypeTest_response{}

	args := []any{id}

	query := `SELECT "TypeTest"."__ID", "TypeTest"."Bool", "TypeTest"."Bool_Arr", "TypeTest"."Date", "TypeTest"."DateTime", "TypeTest"."Date_arr", "TypeTest"."Datetime_Arr", "TypeTest"."Float", "TypeTest"."Float_arr", "TypeTest"."Int", "TypeTest"."Int_Arr", "TypeTest"."Str_Arr", "TypeTest"."String", "TypeTest"."Time", "TypeTest"."Time_Arr" FROM "TypeTest" WHERE "TypeTest"."__ID" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.ID__, &item.Column_Bool, pq.Array(&item.Column_Bool_Arr), &item.Column_Date, &item.Column_DateTime, pq.Array(&item.Column_Date_arr), pq.Array(&item.Column_Datetime_Arr), &item.Column_Float, pq.Array(&item.Column_Float_arr), &item.Column_Int, pq.Array(&item.Column_Int_Arr), pq.Array(&item.Column_Str_Arr), &item.Column_String, &item.Column_Time, pq.Array(&item.Column_Time_Arr)); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_TypeTest(ctx context.Context, id string, item *Table_TypeTest) error {
	args := []any{item.ID__, item.Column_Bool, pq.Array(item.Column_Bool_Arr), item.Column_Date, item.Column_DateTime, pq.Array(item.Column_Date_arr), pq.Array(item.Column_Datetime_Arr), item.Column_Float, pq.Array(item.Column_Float_arr), item.Column_Int, pq.Array(item.Column_Int_Arr), pq.Array(item.Column_Str_Arr), item.Column_String, item.Column_Time, pq.Array(item.Column_Time_Arr), id}

	query := `UPDATE "TypeTest" SET "__ID" = $1, "Bool" = $2, "Bool_Arr" = $3, "Date" = $4, "DateTime" = $5, "Date_arr" = $6, "Datetime_Arr" = $7, "Float" = $8, "Float_arr" = $9, "Int" = $10, "Int_Arr" = $11, "Str_Arr" = $12, "String" = $13, "Time" = $14, "Time_Arr" = $15 WHERE "__ID" = $16`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_TypeTest(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "TypeTest" WHERE "__ID" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

// branches CRUD
func db_insert_branches(ctx context.Context, item *Table_branches) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "branches" ("Branch_Id", "Branch_Name", "Course_Id", "Teachers") VALUES ($1, $2, $3, $4)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Branch_Id, item.Column_Branch_Name, item.Column_Course_Id, pq.Array(item.Column_Teachers))

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_branches(ctx context.Context, clause string, args []any) ([]Table_branches_response, error) {
	data := []Table_branches_response{}

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed", "branches"."Teachers" FROM "branches" INNER JOIN "courses" ON "branches"."Course_Id" = "courses"."Course_Id"`

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_branches_response{}

		rows.Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Fkey_Course_Id.Column_Course_Id, &item.Fkey_Course_Id.Column_Course_Name, &item.Fkey_Course_Id.Column_Lateral_Allowed, pq.Array(&item.Column_Teachers))

		data = append(data, item)
	}

	return data, nil
}

func db_read_branches_ByPK(ctx context.Context, id string) (Table_branches_response, error) {
	item := Table_branches_response{}

	args := []any{id}

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed", "branches"."Teachers" FROM "branches" INNER JOIN "courses" ON "branches"."Course_Id" = "courses"."Course_Id" WHERE "branches"."Branch_Id" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Fkey_Course_Id.Column_Course_Id, &item.Fkey_Course_Id.Column_Course_Name, &item.Fkey_Course_Id.Column_Lateral_Allowed, pq.Array(&item.Column_Teachers)); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_branches(ctx context.Context, id string, item *Table_branches) error {
	args := []any{item.Column_Branch_Id, item.Column_Branch_Name, item.Column_Course_Id, pq.Array(item.Column_Teachers), id}

	query := `UPDATE "branches" SET "Branch_Id" = $1, "Branch_Name" = $2, "Course_Id" = $3, "Teachers" = $4 WHERE "Branch_Id" = $5`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_branches(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "branches" WHERE "Branch_Id" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

// courses CRUD
func db_insert_courses(ctx context.Context, item *Table_courses) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "courses" ("Course_Id", "Course_Name", "Lateral_Allowed") VALUES ($1, $2, $3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Course_Id, item.Column_Course_Name, item.Column_Lateral_Allowed)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_courses(ctx context.Context, clause string, args []any) ([]Table_courses_response, error) {
	data := []Table_courses_response{}

	query := `SELECT "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed" FROM "courses"`

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_courses_response{}

		rows.Scan(&item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed)

		data = append(data, item)
	}

	return data, nil
}

func db_read_courses_ByPK(ctx context.Context, id string) (Table_courses_response, error) {
	item := Table_courses_response{}

	args := []any{id}

	query := `SELECT "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed" FROM "courses" WHERE "courses"."Course_Id" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_courses(ctx context.Context, id string, item *Table_courses) error {
	args := []any{item.Column_Course_Id, item.Column_Course_Name, item.Column_Lateral_Allowed, id}

	query := `UPDATE "courses" SET "Course_Id" = $1, "Course_Name" = $2, "Lateral_Allowed" = $3 WHERE "Course_Id" = $4`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_courses(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "courses" WHERE "Course_Id" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

// AUTH
func db_insert_login(ctx context.Context, item *Table_login) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "login" ("password", "role", "username") VALUES ($1, $2, $3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_password, item.Column_role, item.Column_username)

	if err != nil {
		return err
	}

	return nil
}

func db_auth_login(ctx context.Context, login_data *Login_Input) (Login_Output, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT "username", "password", "role" FROM "login" WHERE username = $1`)

	var item Login_Output

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, login_data.Username).Scan(&item.Username, &item.Password, &item.Role)

	if err != nil {
		return item, err
	}

	return item, nil
}
