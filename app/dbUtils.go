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

// TypeTest CRUD

func db_insert_TypeTest(ctx context.Context, item *Table_TypeTest) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "TypeTest" ("Bool","Bool_Arr","Date","DateTime","Date_arr","Datetime_Arr","Float","Float_arr","Int","Int_Arr","Str_Arr","String","Time","Time_Arr") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`)

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

func db_readAll_TypeTest(ctx context.Context) ([]Table_TypeTest, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM "TypeTest"`)

	data := []Table_TypeTest{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_TypeTest{}

		rows.Scan(&item.ID__, &item.Column_Bool, pq.Array(&item.Column_Bool_Arr), &item.Column_Date, &item.Column_DateTime, pq.Array(&item.Column_Date_arr), pq.Array(&item.Column_Datetime_Arr), &item.Column_Float, pq.Array(&item.Column_Float_arr), &item.Column_Int, pq.Array(&item.Column_Int_Arr), pq.Array(&item.Column_Str_Arr), &item.Column_String, &item.Column_Time, pq.Array(&item.Column_Time_Arr))

		data = append(data, item)
	}

	return data, nil
}

func db_read_TypeTest_ByPK(ctx context.Context, id string) (Table_TypeTest, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM "TypeTest" WHERE "__ID" = $1`)

	item := Table_TypeTest{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRowContext(ctx, id).Scan(&item.ID__, &item.Column_Bool, pq.Array(&item.Column_Bool_Arr), &item.Column_Date, &item.Column_DateTime, pq.Array(&item.Column_Date_arr), pq.Array(&item.Column_Datetime_Arr), &item.Column_Float, pq.Array(&item.Column_Float_arr), &item.Column_Int, pq.Array(&item.Column_Int_Arr), pq.Array(&item.Column_Str_Arr), &item.Column_String, &item.Column_Time, pq.Array(&item.Column_Time_Arr)); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_TypeTest(ctx context.Context, id string, item *Table_TypeTest) error {
	stmt, err := db.PrepareContext(ctx, `UPDATE "TypeTest" SET "__ID" = $1,"Bool" = $2,"Bool_Arr" = $3,"Date" = $4,"DateTime" = $5,"Date_arr" = $6,"Datetime_Arr" = $7,"Float" = $8,"Float_arr" = $9,"Int" = $10,"Int_Arr" = $11,"Str_Arr" = $12,"String" = $13,"Time" = $14,"Time_Arr" = $15 WHERE "__ID" == $16`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, item.ID__, &item.Column_Bool, pq.Array(&item.Column_Bool_Arr), &item.Column_Date, &item.Column_DateTime, pq.Array(&item.Column_Date_arr), pq.Array(&item.Column_Datetime_Arr), &item.Column_Float, pq.Array(&item.Column_Float_arr), &item.Column_Int, pq.Array(&item.Column_Int_Arr), pq.Array(&item.Column_Str_Arr), &item.Column_String, &item.Column_Time, pq.Array(&item.Column_Time_Arr), id)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_TypeTest(ctx context.Context, id string) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM "TypeTest" WHERE "__ID" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)

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
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "branches" ("Branch_Id","Branch_Name","Course_Id","Teachers") VALUES ($1,$2,$3,$4)`)

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

func db_readAll_branches(ctx context.Context) ([]Table_branches, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM "branches"`)

	data := []Table_branches{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_branches{}

		rows.Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Column_Course_Id, pq.Array(&item.Column_Teachers))

		data = append(data, item)
	}

	return data, nil
}

func db_read_branches_ByPK(ctx context.Context, id string) (Table_branches, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM "branches" WHERE "Branch_Id" = $1`)

	item := Table_branches{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRowContext(ctx, id).Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Column_Course_Id, pq.Array(&item.Column_Teachers)); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_branches(ctx context.Context, id string, item *Table_branches) error {
	stmt, err := db.PrepareContext(ctx, `UPDATE "branches" SET "Branch_Id" = $1,"Branch_Name" = $2,"Course_Id" = $3,"Teachers" = $4 WHERE "Branch_Id" = $5`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &item.Column_Branch_Id, &item.Column_Branch_Name, &item.Column_Course_Id, pq.Array(&item.Column_Teachers), id)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_branches(ctx context.Context, id string) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM "branches" WHERE "Branch_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)

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
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "courses" ("Course_Id","Course_Name","Lateral_Allowed") VALUES ($1,$2,$3)`)

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

func db_readAll_courses(ctx context.Context) ([]Table_courses, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM "courses"`)

	data := []Table_courses{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_courses{}

		rows.Scan(&item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed)

		data = append(data, item)
	}

	return data, nil
}

func db_read_courses_ByPK(ctx context.Context, id string) (Table_courses, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM "courses" WHERE "Course_Id" = $1`)

	item := Table_courses{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRowContext(ctx, id).Scan(&item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_courses(ctx context.Context, id string, item *Table_courses) error {
	stmt, err := db.PrepareContext(ctx, `UPDATE "courses" SET "Course_Id" = $1,"Course_Name" = $2,"Lateral_Allowed" = $3 WHERE "Course_Id" = $4`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed, id)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_courses(ctx context.Context, id string) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM "courses" WHERE "Course_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

// empty CRUD

func db_insert_empty(ctx context.Context, item *Table_empty) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "empty" ("Col_1","Col_2","Col_3") VALUES ($1,$2,$3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Col_1, item.Column_Col_2, item.Column_Col_3)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_empty(ctx context.Context) ([]Table_empty, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM "empty"`)

	data := []Table_empty{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_empty{}

		rows.Scan(&item.ID__, &item.Column_Col_1, &item.Column_Col_2, &item.Column_Col_3)

		data = append(data, item)
	}

	return data, nil
}

func db_read_empty_ByPK(ctx context.Context, id string) (Table_empty, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM "empty" WHERE "__ID" = $1`)

	item := Table_empty{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRowContext(ctx, id).Scan(&item.ID__, &item.Column_Col_1, &item.Column_Col_2, &item.Column_Col_3); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_empty(ctx context.Context, id string, item *Table_empty) error {
	stmt, err := db.PrepareContext(ctx, `UPDATE "empty" SET "__ID" = $1,"Col_1" = $2,"Col_2" = $3,"Col_3" = $4 WHERE "__ID" == $5`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, item.ID__, &item.Column_Col_1, &item.Column_Col_2, &item.Column_Col_3, id)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_empty(ctx context.Context, id string) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM "empty" WHERE "__ID" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

// students CRUD

func db_insert_students(ctx context.Context, item *Table_students) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "students" ("Branch_Id","Course_Id","Student_Father","Student_Id","Student_Name") VALUES ($1,$2,$3,$4,$5)`)

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

func db_readAll_students(ctx context.Context) ([]Table_students, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM "students"`)

	data := []Table_students{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_students{}

		rows.Scan(&item.Column_Branch_Id, &item.Column_Course_Id, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name)

		data = append(data, item)
	}

	return data, nil
}

func db_read_students_ByPK(ctx context.Context, id string) (Table_students, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM "students" WHERE "Student_Id" = $1`)

	item := Table_students{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRowContext(ctx, id).Scan(&item.Column_Branch_Id, &item.Column_Course_Id, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_students(ctx context.Context, id string, item *Table_students) error {
	stmt, err := db.PrepareContext(ctx, `UPDATE "students" SET "Branch_Id" = $1,"Course_Id" = $2,"Student_Father" = $3,"Student_Id" = $4,"Student_Name" = $5 WHERE "Student_Id" = $6`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &item.Column_Branch_Id, &item.Column_Course_Id, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name, id)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_students(ctx context.Context, id string) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM "students" WHERE "Student_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)

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
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "subjects" ("Branch_Id","Subject_Id","Subject_Name") VALUES ($1,$2,$3)`)

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

func db_readAll_subjects(ctx context.Context) ([]Table_subjects, error) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM "subjects"`)

	data := []Table_subjects{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_subjects{}

		rows.Scan(&item.Column_Branch_Id, &item.Column_Subject_Id, &item.Column_Subject_Name)

		data = append(data, item)
	}

	return data, nil
}

func db_read_subjects_ByPK(ctx context.Context, id string) (Table_subjects, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT * FROM "subjects" WHERE "Subject_Id" = $1`)

	item := Table_subjects{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRowContext(ctx, id).Scan(&item.Column_Branch_Id, &item.Column_Subject_Id, &item.Column_Subject_Name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_subjects(ctx context.Context, id string, item *Table_subjects) error {
	stmt, err := db.PrepareContext(ctx, `UPDATE "subjects" SET "Branch_Id" = $1,"Subject_Id" = $2,"Subject_Name" = $3 WHERE "Subject_Id" = $4`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, &item.Column_Branch_Id, &item.Column_Subject_Id, &item.Column_Subject_Name, id)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}

func db_delete_subjects(ctx context.Context, id string) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM "subjects" WHERE "Subject_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided id")
	}

	return nil
}
