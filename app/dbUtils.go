package main

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "CSV_App"
)

func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)

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

func db_insert_TypeTest(TypeTest_ *Table_TypeTest) error {
	stmt, err := db.Prepare(`INSERT INTO "TypeTest" ("Bool","Bool_Arr","Date","DateTime","Date_arr","Datetime_Arr","Float","Float_arr","Int","Int_Arr","Str_Arr","String","Time","Time_Arr") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(TypeTest_.Column_Bool, TypeTest_.Column_Bool_Arr, TypeTest_.Column_Date, TypeTest_.Column_DateTime, TypeTest_.Column_Date_arr, TypeTest_.Column_Datetime_Arr, TypeTest_.Column_Float, TypeTest_.Column_Float_arr, TypeTest_.Column_Int, TypeTest_.Column_Int_Arr, TypeTest_.Column_Str_Arr, TypeTest_.Column_String, TypeTest_.Column_Time, TypeTest_.Column_Time_Arr)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_TypeTest() ([]Table_TypeTest, error) {
	rows, err := db.Query(`SELECT * FROM "TypeTest"`)

	data := []Table_TypeTest{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_TypeTest{}

		rows.Scan(&item.ID__, &item.Column_Bool, &item.Column_Bool_Arr, &item.Column_Date, &item.Column_DateTime, &item.Column_Date_arr, &item.Column_Datetime_Arr, &item.Column_Float, &item.Column_Float_arr, &item.Column_Int, &item.Column_Int_Arr, &item.Column_Str_Arr, &item.Column_String, &item.Column_Time, &item.Column_Time_Arr)

		data = append(data, item)
	}

	return data, nil
}

func db_read_TypeTest_ByPK(pk uint) (Table_TypeTest, error) {
	stmt, err := db.Prepare(`SELECT * FROM "TypeTest" WHERE "__ID" = $1`)

	item := Table_TypeTest{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRow(pk).Scan(&item.ID__, &item.Column_Bool, &item.Column_Bool_Arr, &item.Column_Date, &item.Column_DateTime, &item.Column_Date_arr, &item.Column_Datetime_Arr, &item.Column_Float, &item.Column_Float_arr, &item.Column_Int, &item.Column_Int_Arr, &item.Column_Str_Arr, &item.Column_String, &item.Column_Time, &item.Column_Time_Arr); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_TypeTest(item *Table_TypeTest, pk uint) error {
	stmt, err := db.Prepare(`UPDATE "TypeTest" SET "__ID" = $1,"Bool" = $2,"Bool_Arr" = $3,"Date" = $4,"DateTime" = $5,"Date_arr" = $6,"Datetime_Arr" = $7,"Float" = $8,"Float_arr" = $9,"Int" = $10,"Int_Arr" = $11,"Str_Arr" = $12,"String" = $13,"Time" = $14,"Time_Arr" = $15 WHERE "__ID" == $16`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(item.ID__, item.Column_Bool, item.Column_Bool_Arr, item.Column_Date, item.Column_DateTime, item.Column_Date_arr, item.Column_Datetime_Arr, item.Column_Float, item.Column_Float_arr, item.Column_Int, item.Column_Int_Arr, item.Column_Str_Arr, item.Column_String, item.Column_Time, item.Column_Time_Arr, pk)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

func db_delete_TypeTest(pk uint) error {
	stmt, err := db.Prepare(`DELETE FROM "TypeTest" WHERE "__ID" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(pk)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

// branches CRUD

func db_insert_branches(branches_ *Table_branches) error {
	stmt, err := db.Prepare(`INSERT INTO "branches" ("Branch_Id","Branch_Name","Course_Id","Teachers") VALUES ($1,$2,$3,$4)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(branches_.Column_Branch_Id, branches_.Column_Branch_Name, branches_.Column_Course_Id, branches_.Column_Teachers)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_branches() ([]Table_branches, error) {
	rows, err := db.Query(`SELECT * FROM "branches"`)

	data := []Table_branches{}

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		item := Table_branches{}

		rows.Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Column_Course_Id, &item.Column_Teachers)

		data = append(data, item)
	}

	return data, nil
}

func db_read_branches_ByPK(pk int) (Table_branches, error) {
	stmt, err := db.Prepare(`SELECT * FROM "branches" WHERE "Branch_Id" = $1`)

	item := Table_branches{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRow(pk).Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Column_Course_Id, &item.Column_Teachers); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_branches(item *Table_branches, pk int) error {
	stmt, err := db.Prepare(`UPDATE "branches" SET "Branch_Id" = $1,"Branch_Name" = $2,"Course_Id" = $3,"Teachers" = $4 WHERE "Branch_Id" = $5`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(item.Column_Branch_Id, item.Column_Branch_Name, item.Column_Course_Id, item.Column_Teachers, pk)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

func db_delete_branches(pk int) error {
	stmt, err := db.Prepare(`DELETE FROM "branches" WHERE "Branch_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(pk)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

// courses CRUD

func db_insert_courses(courses_ *Table_courses) error {
	stmt, err := db.Prepare(`INSERT INTO "courses" ("Course_Id","Course_Name","Lateral_Allowed") VALUES ($1,$2,$3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(courses_.Column_Course_Id, courses_.Column_Course_Name, courses_.Column_Lateral_Allowed)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_courses() ([]Table_courses, error) {
	rows, err := db.Query(`SELECT * FROM "courses"`)

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

func db_read_courses_ByPK(pk int) (Table_courses, error) {
	stmt, err := db.Prepare(`SELECT * FROM "courses" WHERE "Course_Id" = $1`)

	item := Table_courses{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRow(pk).Scan(&item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_courses(item *Table_courses, pk int) error {
	stmt, err := db.Prepare(`UPDATE "courses" SET "Course_Id" = $1,"Course_Name" = $2,"Lateral_Allowed" = $3 WHERE "Course_Id" = $4`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(item.Column_Course_Id, item.Column_Course_Name, item.Column_Lateral_Allowed, pk)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

func db_delete_courses(pk int) error {
	stmt, err := db.Prepare(`DELETE FROM "courses" WHERE "Course_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(pk)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

// empty CRUD

func db_insert_empty(empty_ *Table_empty) error {
	stmt, err := db.Prepare(`INSERT INTO "empty" ("Col_1","Col_2","Col_3") VALUES ($1,$2,$3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(empty_.Column_Col_1, empty_.Column_Col_2, empty_.Column_Col_3)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_empty() ([]Table_empty, error) {
	rows, err := db.Query(`SELECT * FROM "empty"`)

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

func db_read_empty_ByPK(pk uint) (Table_empty, error) {
	stmt, err := db.Prepare(`SELECT * FROM "empty" WHERE "__ID" = $1`)

	item := Table_empty{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRow(pk).Scan(&item.ID__, &item.Column_Col_1, &item.Column_Col_2, &item.Column_Col_3); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_empty(item *Table_empty, pk uint) error {
	stmt, err := db.Prepare(`UPDATE "empty" SET "__ID" = $1,"Col_1" = $2,"Col_2" = $3,"Col_3" = $4 WHERE "__ID" == $5`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(item.ID__, item.Column_Col_1, item.Column_Col_2, item.Column_Col_3, pk)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

func db_delete_empty(pk uint) error {
	stmt, err := db.Prepare(`DELETE FROM "empty" WHERE "__ID" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(pk)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

// students CRUD

func db_insert_students(students_ *Table_students) error {
	stmt, err := db.Prepare(`INSERT INTO "students" ("Branch_Id","Course_Id","Student_Father","Student_Id","Student_Name") VALUES ($1,$2,$3,$4,$5)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(students_.Column_Branch_Id, students_.Column_Course_Id, students_.Column_Student_Father, students_.Column_Student_Id, students_.Column_Student_Name)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_students() ([]Table_students, error) {
	rows, err := db.Query(`SELECT * FROM "students"`)

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

func db_read_students_ByPK(pk int) (Table_students, error) {
	stmt, err := db.Prepare(`SELECT * FROM "students" WHERE "Student_Id" = $1`)

	item := Table_students{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRow(pk).Scan(&item.Column_Branch_Id, &item.Column_Course_Id, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_students(item *Table_students, pk int) error {
	stmt, err := db.Prepare(`UPDATE "students" SET "Branch_Id" = $1,"Course_Id" = $2,"Student_Father" = $3,"Student_Id" = $4,"Student_Name" = $5 WHERE "Student_Id" = $6`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(item.Column_Branch_Id, item.Column_Course_Id, item.Column_Student_Father, item.Column_Student_Id, item.Column_Student_Name, pk)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

func db_delete_students(pk int) error {
	stmt, err := db.Prepare(`DELETE FROM "students" WHERE "Student_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(pk)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

// subjects CRUD

func db_insert_subjects(subjects_ *Table_subjects) error {
	stmt, err := db.Prepare(`INSERT INTO "subjects" ("Branch_Id","Subject_Id","Subject_Name") VALUES ($1,$2,$3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(subjects_.Column_Branch_Id, subjects_.Column_Subject_Id, subjects_.Column_Subject_Name)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_subjects() ([]Table_subjects, error) {
	rows, err := db.Query(`SELECT * FROM "subjects"`)

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

func db_read_subjects_ByPK(pk int) (Table_subjects, error) {
	stmt, err := db.Prepare(`SELECT * FROM "subjects" WHERE "Subject_Id" = $1`)

	item := Table_subjects{}

	if err != nil {
		return item, err
	}

	if err := stmt.QueryRow(pk).Scan(&item.Column_Branch_Id, &item.Column_Subject_Id, &item.Column_Subject_Name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_subjects(item *Table_subjects, pk int) error {
	stmt, err := db.Prepare(`UPDATE "subjects" SET "Branch_Id" = $1,"Subject_Id" = $2,"Subject_Name" = $3 WHERE "Subject_Id" = $4`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(item.Column_Branch_Id, item.Column_Subject_Id, item.Column_Subject_Name, pk)

	if err != nil {
		return err
	}

	if rowsUpdated, _ := result.RowsAffected(); rowsUpdated == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}

func db_delete_subjects(pk int) error {
	stmt, err := db.Prepare(`DELETE FROM "subjects" WHERE "Subject_Id" = $1`)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(pk)

	if err != nil {
		return err
	}

	if rowsDeleted, _ := result.RowsAffected(); rowsDeleted == 0 {
		return errors.New("no row found with provided pk")
	}

	return nil
}
