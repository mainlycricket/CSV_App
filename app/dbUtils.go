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

// AUTH

func db_auth_login(ctx context.Context, login_data *Login_Input) (Login_Output, error) {
	stmt, err := db.PrepareContext(ctx, `SELECT "username", "password", "role", "college_id", "course_id", "branch_id" FROM "login" WHERE username = $1`)

	var item Login_Output

	if err != nil {
	}

	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, login_data.Username).Scan(&item.Username, &item.Password, &item.Role, &item.College_id, &item.Course_id, &item.Branch_id)

	if err != nil {
		return item, err
	}

	return item, nil
}

// login CRUD
func db_insert_login(ctx context.Context, item *Table_login) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "login" ("added_by", "branch_id", "college_id", "course_id", "password", "role", "username") VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_added_by, item.Column_branch_id, item.Column_college_id, item.Column_course_id, item.Column_password, item.Column_role, item.Column_username)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_login(ctx context.Context, clause string, args []any) ([]Table_login_ResponseAll, bool, error) {
	recordsCount := args[len(args)-1].(int)
	data := make([]Table_login_ResponseAll, 0, recordsCount-1)

	query := `SELECT "login"."added_by", "login"."branch_id", "login"."college_id", "login"."course_id", "login"."role", "login"."username" FROM "login"`

	var tokenClauses []string

	ctxVal := ctx.Value(ContextKey("branch_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"login"."branch_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"login"."branch_id" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"login"."college_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"login"."college_id" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"login"."course_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"login"."course_id" IS NULL`))
		}
	}

	if len(tokenClauses) > 0 {
		if strings.HasPrefix(clause, " WHERE") {
			clause = strings.TrimPrefix(clause, " WHERE")
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + " AND " + clause
		} else {
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + clause
		}
	}

	var protectClauses []string
	role := ctx.Value(ContextKey("__role")).(string)
	var rolesMap map[string][]interface{}

	rolesMap = map[string][]interface{}{"admin": []interface{}{"hod", "teacher", "student"}, "hod": []interface{}{"principal", "hod"}, "principal": []interface{}{"principal"}}
	if disallowedValues := rolesMap[role]; len(disallowedValues) > 0 {

		protectClauses = append(protectClauses, fmt.Sprintf(`"login"."role" NOT IN (%s)`, getArgPlaceHolders(len(args)+1, len(disallowedValues))))

		args = append(args, disallowedValues...)
	}

	if len(protectClauses) > 0 {
		if strings.HasPrefix(clause, " WHERE") {
			clause = strings.TrimPrefix(clause, " WHERE")
			clause = " WHERE " + strings.Join(protectClauses, " AND ") + " AND " + clause
		} else {
			clause = " WHERE " + strings.Join(protectClauses, " AND ") + clause
		}
	}

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, false, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, false, err
	}

	defer rows.Close()
	nextFlag := false

	for rows.Next() {
		recordsCount--
		if recordsCount == 0 {
			nextFlag = true
			break
		}

		item := Table_login_ResponseAll{}

		rows.Scan(&item.Column_added_by, &item.Column_branch_id, &item.Column_college_id, &item.Column_course_id, &item.Column_role, &item.Column_username)

		data = append(data, item)
	}

	return data, nextFlag, nil
}

func db_read_login_ByPK(ctx context.Context, id string) (Table_login_ResponsePK, error) {
	item := Table_login_ResponsePK{}

	args := []any{id}

	query := `SELECT "login"."added_by", "login"."branch_id", "login"."college_id", "login"."course_id", "login"."role", "login"."username" FROM "login" WHERE "login"."username" = $1`

	ctxVal := ctx.Value(ContextKey("branch_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."branch_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."branch_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."college_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."course_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."course_id" IS NULL`)
		}
	}

	role := ctx.Value(ContextKey("__role")).(string)
	var rolesMap map[string][]interface{}

	rolesMap = map[string][]interface{}{"admin": []interface{}{"hod", "teacher", "student"}, "hod": []interface{}{"hod", "principal"}, "principal": []interface{}{"principal"}}
	if disallowedValues := rolesMap[role]; len(disallowedValues) > 0 {

		query += fmt.Sprintf(` AND "login"."role" NOT IN (%s)`, getArgPlaceHolders(len(args)+1, len(disallowedValues)))

		args = append(args, disallowedValues...)
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_added_by, &item.Column_branch_id, &item.Column_college_id, &item.Column_course_id, &item.Column_role, &item.Column_username); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_login(ctx context.Context, id string, item *Table_login) error {
	args := []any{item.Column_added_by, item.Column_branch_id, item.Column_college_id, item.Column_course_id, item.Column_password, item.Column_role, item.Column_username, id}

	query := `UPDATE "login" SET "added_by" = $1, "branch_id" = $2, "college_id" = $3, "course_id" = $4, "password" = $5, "role" = $6, "username" = $7 WHERE "username" = $8`

	ctxVal := ctx.Value(ContextKey("username"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."username" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."username" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("branch_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."branch_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."branch_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."college_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."course_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."course_id" IS NULL`)
		}
	}

	role := ctx.Value(ContextKey("__role")).(string)
	var rolesMap map[string][]interface{}

	rolesMap = map[string][]interface{}{"admin": []interface{}{"hod", "teacher"}, "hod": []interface{}{"principal", "hod"}, "principal": []interface{}{"principal"}}
	if disallowedValues := rolesMap[role]; len(disallowedValues) > 0 {

		query += fmt.Sprintf(` AND "login"."role" NOT IN (%s)`, getArgPlaceHolders(len(args)+1, len(disallowedValues)))

		args = append(args, disallowedValues...)
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

func db_delete_login(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "login" WHERE "username" = $1`

	ctxVal := ctx.Value(ContextKey("username"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."username" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."username" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("branch_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."branch_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."branch_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."college_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "login"."course_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "login"."course_id" IS NULL`)
		}
	}

	role := ctx.Value(ContextKey("__role")).(string)
	var rolesMap map[string][]interface{}

	rolesMap = map[string][]interface{}{"admin": []interface{}{"teacher", "hod"}, "hod": []interface{}{"principal", "hod"}, "principal": []interface{}{"principal"}}
	if disallowedValues := rolesMap[role]; len(disallowedValues) > 0 {

		query += fmt.Sprintf(` AND "login"."role" NOT IN (%s)`, getArgPlaceHolders(len(args)+1, len(disallowedValues)))

		args = append(args, disallowedValues...)
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

// students CRUD
func db_insert_students(ctx context.Context, item *Table_students) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "students" ("Branch_Id", "Course_Id", "Student_Father", "Student_Id", "Student_Name", "added_by", "college_id") VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Branch_Id, item.Column_Course_Id, item.Column_Student_Father, item.Column_Student_Id, item.Column_Student_Name, item.Column_added_by, item.Column_college_id)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_students(ctx context.Context, clause string, args []any) ([]Table_students_ResponseAll, bool, error) {
	recordsCount := args[len(args)-1].(int)
	data := make([]Table_students_ResponseAll, 0, recordsCount-1)

	query := `SELECT "Branch_Id_branches"."Branch_Id", "Branch_Id_branches"."Branch_Name", "Course_Id_courses"."Course_Id", "Course_Id_courses"."Course_Name", "students"."Student_Father", "students"."Student_Id", "students"."Student_Name", "students"."added_by", "college_id_college"."college_id", "college_id_college"."college_name" FROM "students" LEFT JOIN "branches" AS "Branch_Id_branches" ON "students"."Branch_Id" = "Branch_Id_branches"."Branch_Id" LEFT JOIN "courses" AS "Course_Id_courses" ON "students"."Course_Id" = "Course_Id_courses"."Course_Id" LEFT JOIN "college" AS "college_id_college" ON "students"."college_id" = "college_id_college"."college_id"`

	var tokenClauses []string

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."added_by" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."added_by" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."Branch_Id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."Branch_Id" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."Course_Id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."Course_Id" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."college_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"students"."college_id" IS NULL`))
		}
	}

	if len(tokenClauses) > 0 {
		if strings.HasPrefix(clause, " WHERE") {
			clause = strings.TrimPrefix(clause, " WHERE")
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + " AND " + clause
		} else {
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + clause
		}
	}

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, false, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, false, err
	}

	defer rows.Close()
	nextFlag := false

	for rows.Next() {
		recordsCount--
		if recordsCount == 0 {
			nextFlag = true
			break
		}

		item := Table_students_ResponseAll{}

		rows.Scan(&item.Column_Branch_Id.Column_Branch_Id, &item.Column_Branch_Id.Column_Branch_Name, &item.Column_Course_Id.Column_Course_Id, &item.Column_Course_Id.Column_Course_Name, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name, &item.Column_added_by, &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name)

		data = append(data, item)
	}

	return data, nextFlag, nil
}

func db_read_students_ByPK(ctx context.Context, id string) (Table_students_ResponsePK, error) {
	item := Table_students_ResponsePK{}

	args := []any{id}

	query := `SELECT "Branch_Id_branches"."Branch_Id", "Branch_Id_branches"."Branch_Name", "Course_Id_courses"."Course_Id", "Course_Id_courses"."Course_Name", "students"."Student_Father", "students"."Student_Id", "students"."Student_Name", "students"."added_by", "college_id_college"."college_id", "college_id_college"."college_name" FROM "students" LEFT JOIN "branches" AS "Branch_Id_branches" ON "students"."Branch_Id" = "Branch_Id_branches"."Branch_Id" LEFT JOIN "courses" AS "Course_Id_courses" ON "students"."Course_Id" = "Course_Id_courses"."Course_Id" LEFT JOIN "college" AS "college_id_college" ON "students"."college_id" = "college_id_college"."college_id" WHERE "students"."Student_Id" = $1`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."Branch_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."Branch_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."Course_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."Course_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."college_id" IS NULL`)
		}
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_Branch_Id.Column_Branch_Id, &item.Column_Branch_Id.Column_Branch_Name, &item.Column_Course_Id.Column_Course_Id, &item.Column_Course_Id.Column_Course_Name, &item.Column_Student_Father, &item.Column_Student_Id, &item.Column_Student_Name, &item.Column_added_by, &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_students(ctx context.Context, id string, item *Table_students) error {
	args := []any{item.Column_Branch_Id, item.Column_Course_Id, item.Column_Student_Father, item.Column_Student_Id, item.Column_Student_Name, item.Column_added_by, item.Column_college_id, id}

	query := `UPDATE "students" SET "Branch_Id" = $1, "Course_Id" = $2, "Student_Father" = $3, "Student_Id" = $4, "Student_Name" = $5, "added_by" = $6, "college_id" = $7 WHERE "Student_Id" = $8`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."Branch_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."Branch_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."Course_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."Course_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."college_id" IS NULL`)
		}
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

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."Branch_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."Branch_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."Course_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."Course_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "students"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "students"."college_id" IS NULL`)
		}
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
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "subjects" ("Branch_Id", "Subject_Id", "Subject_Name", "added_by", "college_id", "course_id") VALUES ($1, $2, $3, $4, $5, $6)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Branch_Id, item.Column_Subject_Id, item.Column_Subject_Name, item.Column_added_by, item.Column_college_id, item.Column_course_id)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_subjects(ctx context.Context, clause string, args []any) ([]Table_subjects_ResponseAll, bool, error) {
	recordsCount := args[len(args)-1].(int)
	data := make([]Table_subjects_ResponseAll, 0, recordsCount-1)

	query := `SELECT "Branch_Id_branches"."Branch_Id", "Branch_Id_branches"."Branch_Name", "subjects"."Subject_Id", "subjects"."Subject_Name", "subjects"."added_by", "college_id_college"."college_id", "college_id_college"."college_name", "course_id_courses"."Course_Id", "course_id_courses"."Course_Name" FROM "subjects" LEFT JOIN "branches" AS "Branch_Id_branches" ON "subjects"."Branch_Id" = "Branch_Id_branches"."Branch_Id" LEFT JOIN "college" AS "college_id_college" ON "subjects"."college_id" = "college_id_college"."college_id" LEFT JOIN "courses" AS "course_id_courses" ON "subjects"."course_id" = "course_id_courses"."Course_Id"`

	var tokenClauses []string

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."added_by" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."added_by" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."Branch_Id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."Branch_Id" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."college_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."college_id" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."course_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"subjects"."course_id" IS NULL`))
		}
	}

	if len(tokenClauses) > 0 {
		if strings.HasPrefix(clause, " WHERE") {
			clause = strings.TrimPrefix(clause, " WHERE")
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + " AND " + clause
		} else {
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + clause
		}
	}

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, false, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, false, err
	}

	defer rows.Close()
	nextFlag := false

	for rows.Next() {
		recordsCount--
		if recordsCount == 0 {
			nextFlag = true
			break
		}

		item := Table_subjects_ResponseAll{}

		rows.Scan(&item.Column_Branch_Id.Column_Branch_Id, &item.Column_Branch_Id.Column_Branch_Name, &item.Column_Subject_Id, &item.Column_Subject_Name, &item.Column_added_by, &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name, &item.Column_course_id.Column_Course_Id, &item.Column_course_id.Column_Course_Name)

		data = append(data, item)
	}

	return data, nextFlag, nil
}

func db_read_subjects_ByPK(ctx context.Context, id string) (Table_subjects_ResponsePK, error) {
	item := Table_subjects_ResponsePK{}

	args := []any{id}

	query := `SELECT "Branch_Id_branches"."Branch_Id", "Branch_Id_branches"."Branch_Name", "subjects"."Subject_Id", "subjects"."Subject_Name", "subjects"."added_by", "college_id_college"."college_id", "college_id_college"."college_name", "course_id_courses"."Course_Id", "course_id_courses"."Course_Name" FROM "subjects" LEFT JOIN "branches" AS "Branch_Id_branches" ON "subjects"."Branch_Id" = "Branch_Id_branches"."Branch_Id" LEFT JOIN "college" AS "college_id_college" ON "subjects"."college_id" = "college_id_college"."college_id" LEFT JOIN "courses" AS "course_id_courses" ON "subjects"."course_id" = "course_id_courses"."Course_Id" WHERE "subjects"."Subject_Id" = $1`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."Branch_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."Branch_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."college_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."course_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."course_id" IS NULL`)
		}
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_Branch_Id.Column_Branch_Id, &item.Column_Branch_Id.Column_Branch_Name, &item.Column_Subject_Id, &item.Column_Subject_Name, &item.Column_added_by, &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name, &item.Column_course_id.Column_Course_Id, &item.Column_course_id.Column_Course_Name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_subjects(ctx context.Context, id string, item *Table_subjects) error {
	args := []any{item.Column_Branch_Id, item.Column_Subject_Id, item.Column_Subject_Name, item.Column_added_by, item.Column_college_id, item.Column_course_id, id}

	query := `UPDATE "subjects" SET "Branch_Id" = $1, "Subject_Id" = $2, "Subject_Name" = $3, "added_by" = $4, "college_id" = $5, "course_id" = $6 WHERE "Subject_Id" = $7`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."Branch_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."Branch_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."college_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."course_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."course_id" IS NULL`)
		}
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

func db_delete_subjects(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "subjects" WHERE "Subject_Id" = $1`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Branch_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."Branch_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."Branch_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."college_id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("course_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "subjects"."course_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "subjects"."course_id" IS NULL`)
		}
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

func db_readAll_TypeTest(ctx context.Context, clause string, args []any) ([]Table_TypeTest_ResponseAll, bool, error) {
	recordsCount := args[len(args)-1].(int)
	data := make([]Table_TypeTest_ResponseAll, 0, recordsCount-1)

	query := `SELECT "TypeTest"."Bool", "TypeTest"."Bool_Arr", "TypeTest"."Date", "TypeTest"."DateTime", "TypeTest"."Date_arr", "TypeTest"."Datetime_Arr", "TypeTest"."Float", "TypeTest"."Float_arr", "TypeTest"."Int", "TypeTest"."Int_Arr", "TypeTest"."Str_Arr", "TypeTest"."String", "TypeTest"."Time", "TypeTest"."Time_Arr" FROM "TypeTest"`

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, false, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, false, err
	}

	defer rows.Close()
	nextFlag := false

	for rows.Next() {
		recordsCount--
		if recordsCount == 0 {
			nextFlag = true
			break
		}

		item := Table_TypeTest_ResponseAll{}

		rows.Scan(&item.Column_Bool, pq.Array(&item.Column_Bool_Arr), &item.Column_Date, &item.Column_DateTime, pq.Array(&item.Column_Date_arr), pq.Array(&item.Column_Datetime_Arr), &item.Column_Float, pq.Array(&item.Column_Float_arr), &item.Column_Int, pq.Array(&item.Column_Int_Arr), pq.Array(&item.Column_Str_Arr), &item.Column_String, &item.Column_Time, pq.Array(&item.Column_Time_Arr))

		data = append(data, item)
	}

	return data, nextFlag, nil
}

func db_read_TypeTest_ByPK(ctx context.Context, id string) (Table_TypeTest_ResponsePK, error) {
	item := Table_TypeTest_ResponsePK{}

	args := []any{id}

	query := `SELECT "TypeTest"."Bool", "TypeTest"."Bool_Arr", "TypeTest"."Date", "TypeTest"."DateTime", "TypeTest"."Date_arr", "TypeTest"."Datetime_Arr", "TypeTest"."Float", "TypeTest"."Float_arr", "TypeTest"."Int", "TypeTest"."Int_Arr", "TypeTest"."Str_Arr", "TypeTest"."String", "TypeTest"."Time", "TypeTest"."Time_Arr" FROM "TypeTest" WHERE "TypeTest"."__ID" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_Bool, pq.Array(&item.Column_Bool_Arr), &item.Column_Date, &item.Column_DateTime, pq.Array(&item.Column_Date_arr), pq.Array(&item.Column_Datetime_Arr), &item.Column_Float, pq.Array(&item.Column_Float_arr), &item.Column_Int, pq.Array(&item.Column_Int_Arr), pq.Array(&item.Column_Str_Arr), &item.Column_String, &item.Column_Time, pq.Array(&item.Column_Time_Arr)); err != nil {
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
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "branches" ("Branch_Id", "Branch_Name", "Course_Id", "HoD", "Teachers", "added_by", "college_id") VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Branch_Id, item.Column_Branch_Name, item.Column_Course_Id, item.Column_HoD, pq.Array(item.Column_Teachers), item.Column_added_by, item.Column_college_id)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_branches(ctx context.Context, clause string, args []any) ([]Table_branches_ResponseAll, bool, error) {
	recordsCount := args[len(args)-1].(int)
	data := make([]Table_branches_ResponseAll, 0, recordsCount-1)

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "Course_Id_courses"."Course_Id", "Course_Id_courses"."Course_Name", "branches"."HoD", "branches"."Teachers", "college_id_college"."college_id", "college_id_college"."college_name" FROM "branches" LEFT JOIN "courses" AS "Course_Id_courses" ON "branches"."Course_Id" = "Course_Id_courses"."Course_Id" LEFT JOIN "college" AS "college_id_college" ON "branches"."college_id" = "college_id_college"."college_id"`

	var tokenClauses []string

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"branches"."added_by" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"branches"."added_by" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"branches"."Course_Id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"branches"."Course_Id" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"branches"."college_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"branches"."college_id" IS NULL`))
		}
	}

	if len(tokenClauses) > 0 {
		if strings.HasPrefix(clause, " WHERE") {
			clause = strings.TrimPrefix(clause, " WHERE")
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + " AND " + clause
		} else {
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + clause
		}
	}

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, false, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, false, err
	}

	defer rows.Close()
	nextFlag := false

	for rows.Next() {
		recordsCount--
		if recordsCount == 0 {
			nextFlag = true
			break
		}

		item := Table_branches_ResponseAll{}

		rows.Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Column_Course_Id.Column_Course_Id, &item.Column_Course_Id.Column_Course_Name, &item.Column_HoD, pq.Array(&item.Column_Teachers), &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name)

		data = append(data, item)
	}

	return data, nextFlag, nil
}

func db_read_branches_ByPK(ctx context.Context, id string) (Table_branches_ResponsePK, error) {
	item := Table_branches_ResponsePK{}

	args := []any{id}

	query := `SELECT "branches"."Branch_Id", "branches"."Branch_Name", "Course_Id_courses"."Course_Id", "Course_Id_courses"."Course_Name", "branches"."HoD", "branches"."Teachers", "branches"."added_by", "college_id_college"."college_id", "college_id_college"."college_name" FROM "branches" LEFT JOIN "courses" AS "Course_Id_courses" ON "branches"."Course_Id" = "Course_Id_courses"."Course_Id" LEFT JOIN "college" AS "college_id_college" ON "branches"."college_id" = "college_id_college"."college_id" WHERE "branches"."Branch_Id" = $1`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."Course_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."Course_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."college_id" IS NULL`)
		}
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_Branch_Id, &item.Column_Branch_Name, &item.Column_Course_Id.Column_Course_Id, &item.Column_Course_Id.Column_Course_Name, &item.Column_HoD, pq.Array(&item.Column_Teachers), &item.Column_added_by, &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_branches(ctx context.Context, id string, item *Table_branches) error {
	args := []any{item.Column_Branch_Id, item.Column_Branch_Name, item.Column_Course_Id, item.Column_HoD, pq.Array(item.Column_Teachers), item.Column_added_by, item.Column_college_id, id}

	query := `UPDATE "branches" SET "Branch_Id" = $1, "Branch_Name" = $2, "Course_Id" = $3, "HoD" = $4, "Teachers" = $5, "added_by" = $6, "college_id" = $7 WHERE "Branch_Id" = $8`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."Course_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."Course_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."college_id" IS NULL`)
		}
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

func db_delete_branches(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "branches" WHERE "Branch_Id" = $1`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("Course_Id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."Course_Id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."Course_Id" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "branches"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "branches"."college_id" IS NULL`)
		}
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

// college CRUD
func db_insert_college(ctx context.Context, item *Table_college) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "college" ("college_id", "college_name", "principal_id") VALUES ($1, $2, $3)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_college_id, item.Column_college_name, item.Column_principal_id)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_college(ctx context.Context, clause string, args []any) ([]Table_college_ResponseAll, bool, error) {
	recordsCount := args[len(args)-1].(int)
	data := make([]Table_college_ResponseAll, 0, recordsCount-1)

	query := `SELECT "college"."college_id", "college"."college_name", "college"."principal_id" FROM "college"`

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, false, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, false, err
	}

	defer rows.Close()
	nextFlag := false

	for rows.Next() {
		recordsCount--
		if recordsCount == 0 {
			nextFlag = true
			break
		}

		item := Table_college_ResponseAll{}

		rows.Scan(&item.Column_college_id, &item.Column_college_name, &item.Column_principal_id)

		data = append(data, item)
	}

	return data, nextFlag, nil
}

func db_read_college_ByPK(ctx context.Context, id string) (Table_college_ResponsePK, error) {
	item := Table_college_ResponsePK{}

	args := []any{id}

	query := `SELECT "college"."college_id", "college"."college_name", "college"."principal_id" FROM "college" WHERE "college"."college_id" = $1`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_college_id, &item.Column_college_name, &item.Column_principal_id); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_college(ctx context.Context, id string, item *Table_college) error {
	args := []any{item.Column_college_id, item.Column_college_name, item.Column_principal_id, id}

	query := `UPDATE "college" SET "college_id" = $1, "college_name" = $2, "principal_id" = $3 WHERE "college_id" = $4`

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

func db_delete_college(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "college" WHERE "college_id" = $1`

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
	stmt, err := db.PrepareContext(ctx, `INSERT INTO "courses" ("Course_Id", "Course_Name", "Lateral_Allowed", "added_by", "college_id") VALUES ($1, $2, $3, $4, $5)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, item.Column_Course_Id, item.Column_Course_Name, item.Column_Lateral_Allowed, item.Column_added_by, item.Column_college_id)

	if err != nil {
		return err
	}

	return nil
}

func db_readAll_courses(ctx context.Context, clause string, args []any) ([]Table_courses_ResponseAll, bool, error) {
	recordsCount := args[len(args)-1].(int)
	data := make([]Table_courses_ResponseAll, 0, recordsCount-1)

	query := `SELECT "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed", "courses"."added_by", "college_id_college"."college_id", "college_id_college"."college_name" FROM "courses" LEFT JOIN "college" AS "college_id_college" ON "courses"."college_id" = "college_id_college"."college_id"`

	var tokenClauses []string

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"courses"."added_by" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"courses"."added_by" IS NULL`))
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"courses"."college_id" = $%d`, len(args)))
		} else {
			tokenClauses = append(tokenClauses, fmt.Sprintf(`"courses"."college_id" IS NULL`))
		}
	}

	if len(tokenClauses) > 0 {
		if strings.HasPrefix(clause, " WHERE") {
			clause = strings.TrimPrefix(clause, " WHERE")
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + " AND " + clause
		} else {
			clause = " WHERE " + strings.Join(tokenClauses, " AND ") + clause
		}
	}

	query += clause

	preparedQuery, err := db.PrepareContext(ctx, query)

	if err != nil {
		return data, false, err
	}

	defer preparedQuery.Close()

	rows, err := preparedQuery.QueryContext(ctx, args...)

	if err != nil {
		return data, false, err
	}

	defer rows.Close()
	nextFlag := false

	for rows.Next() {
		recordsCount--
		if recordsCount == 0 {
			nextFlag = true
			break
		}

		item := Table_courses_ResponseAll{}

		rows.Scan(&item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed, &item.Column_added_by, &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name)

		data = append(data, item)
	}

	return data, nextFlag, nil
}

func db_read_courses_ByPK(ctx context.Context, id string) (Table_courses_ResponsePK, error) {
	item := Table_courses_ResponsePK{}

	args := []any{id}

	query := `SELECT "courses"."Course_Id", "courses"."Course_Name", "courses"."Lateral_Allowed", "courses"."added_by", "college_id_college"."college_id", "college_id_college"."college_name", "college_id_college"."principal_id" FROM "courses" LEFT JOIN "college" AS "college_id_college" ON "courses"."college_id" = "college_id_college"."college_id" WHERE "courses"."Course_Id" = $1`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "courses"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "courses"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "courses"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "courses"."college_id" IS NULL`)
		}
	}

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return item, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(&item.Column_Course_Id, &item.Column_Course_Name, &item.Column_Lateral_Allowed, &item.Column_added_by, &item.Column_college_id.Column_college_id, &item.Column_college_id.Column_college_name, &item.Column_college_id.Column_principal_id); err != nil {
		return item, err
	}

	return item, nil
}

func db_update_courses(ctx context.Context, id string, item *Table_courses) error {
	args := []any{item.Column_Course_Id, item.Column_Course_Name, item.Column_Lateral_Allowed, item.Column_added_by, item.Column_college_id, id}

	query := `UPDATE "courses" SET "Course_Id" = $1, "Course_Name" = $2, "Lateral_Allowed" = $3, "added_by" = $4, "college_id" = $5 WHERE "Course_Id" = $6`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "courses"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "courses"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "courses"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "courses"."college_id" IS NULL`)
		}
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

func db_delete_courses(ctx context.Context, id string) error {
	args := []any{id}

	query := `DELETE FROM "courses" WHERE "Course_Id" = $1`

	ctxVal := ctx.Value(ContextKey("added_by"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "courses"."added_by" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "courses"."added_by" IS NULL`)
		}
	}

	ctxVal = ctx.Value(ContextKey("college_id"))

	if ctxVal != nil {
		value := ctxVal.(string)
		if len(value) > 0 {
			args = append(args, value)
			query += fmt.Sprintf(` AND "courses"."college_id" = $%d`, len(args))
		} else {
			query += fmt.Sprintf(` AND "courses"."college_id" IS NULL`)
		}
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
