package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func getJsonResponse(success bool, message string, data any) []byte {
	response := ApiResponse{Success: success, Message: message, Data: data}
	jsonData, err := json.Marshal(response)

	if err != nil {
		response := ApiResponse{Success: false, Message: "failed to jsonify response"}
		jsonData, _ := json.Marshal(response)
		return jsonData
	}

	return jsonData
}

func startServer() *http.Server {

	// students routes
	http.HandleFunc("POST /students", api_create_students)
	http.HandleFunc("GET /students", api_getAll_students)
	http.HandleFunc("GET /studentsByPK", api_getByPk_students)
	http.HandleFunc("PUT /students", api_update_students)
	http.HandleFunc("DELETE /students", api_delete_students)

	// subjects routes
	http.HandleFunc("POST /subjects", api_create_subjects)
	http.HandleFunc("GET /subjects", api_getAll_subjects)
	http.HandleFunc("GET /subjectsByPK", api_getByPk_subjects)
	http.HandleFunc("PUT /subjects", api_update_subjects)
	http.HandleFunc("DELETE /subjects", api_delete_subjects)

	// TypeTest routes
	http.HandleFunc("POST /TypeTest", api_create_TypeTest)
	http.HandleFunc("GET /TypeTest", api_getAll_TypeTest)
	http.HandleFunc("GET /TypeTestByPK", api_getByPk_TypeTest)
	http.HandleFunc("PUT /TypeTest", api_update_TypeTest)
	http.HandleFunc("DELETE /TypeTest", api_delete_TypeTest)

	// branches routes
	http.HandleFunc("POST /branches", api_create_branches)
	http.HandleFunc("GET /branches", api_getAll_branches)
	http.HandleFunc("GET /branchesByPK", api_getByPk_branches)
	http.HandleFunc("PUT /branches", api_update_branches)
	http.HandleFunc("DELETE /branches", api_delete_branches)

	// college routes
	http.HandleFunc("POST /college", api_create_college)
	http.HandleFunc("GET /college", api_getAll_college)
	http.HandleFunc("GET /collegeByPK", api_getByPk_college)
	http.HandleFunc("PUT /college", api_update_college)
	http.HandleFunc("DELETE /college", api_delete_college)

	// courses routes
	http.HandleFunc("POST /courses", api_create_courses)
	http.HandleFunc("GET /courses", api_getAll_courses)
	http.HandleFunc("GET /coursesByPK", api_getByPk_courses)
	http.HandleFunc("PUT /courses", api_update_courses)
	http.HandleFunc("DELETE /courses", api_delete_courses)

	// login routes
	http.HandleFunc("POST /login", api_create_login)
	http.HandleFunc("GET /login", api_getAll_login)
	http.HandleFunc("GET /loginByPK", api_getByPk_login)
	http.HandleFunc("PUT /login", api_update_login)
	http.HandleFunc("DELETE /login", api_delete_login)

	s := &http.Server{
		Addr: ":8080",
	}

	return s
}

// students handler functions

func api_create_students(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_students

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_insert_students(ctx, &item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_students(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var clause string
	var args []any

	if len(queryValues) > 0 {
		clause, args, err = getQueryClauseArgs(queryValues, Map_students, "students")

		if err != nil {
			message := fmt.Sprintf("error while parsing request query: %v", err)
			log.Print(message)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(getJsonResponse(false, message, nil))
			return
		}
	}

	ctx := r.Context()
	data, err := db_readAll_students(ctx, clause, args)

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_students(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	data, err := db_read_students_ByPK(ctx, id)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_students(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var item Table_students

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_students(ctx, id, &item); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_students(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_delete_students(ctx, id); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}

// subjects handler functions

func api_create_subjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_subjects

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_insert_subjects(ctx, &item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_subjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var clause string
	var args []any

	if len(queryValues) > 0 {
		clause, args, err = getQueryClauseArgs(queryValues, Map_subjects, "subjects")

		if err != nil {
			message := fmt.Sprintf("error while parsing request query: %v", err)
			log.Print(message)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(getJsonResponse(false, message, nil))
			return
		}
	}

	ctx := r.Context()
	data, err := db_readAll_subjects(ctx, clause, args)

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_subjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	data, err := db_read_subjects_ByPK(ctx, id)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_subjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var item Table_subjects

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_subjects(ctx, id, &item); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_subjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_delete_subjects(ctx, id); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}

// TypeTest handler functions

func api_create_TypeTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_TypeTest

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_insert_TypeTest(ctx, &item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_TypeTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var clause string
	var args []any

	if len(queryValues) > 0 {
		clause, args, err = getQueryClauseArgs(queryValues, Map_TypeTest, "TypeTest")

		if err != nil {
			message := fmt.Sprintf("error while parsing request query: %v", err)
			log.Print(message)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(getJsonResponse(false, message, nil))
			return
		}
	}

	ctx := r.Context()
	data, err := db_readAll_TypeTest(ctx, clause, args)

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_TypeTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	data, err := db_read_TypeTest_ByPK(ctx, id)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_TypeTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var item Table_TypeTest

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_TypeTest(ctx, id, &item); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_TypeTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullInt")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_delete_TypeTest(ctx, id); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}

// branches handler functions

func api_create_branches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_branches

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_insert_branches(ctx, &item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_branches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var clause string
	var args []any

	if len(queryValues) > 0 {
		clause, args, err = getQueryClauseArgs(queryValues, Map_branches, "branches")

		if err != nil {
			message := fmt.Sprintf("error while parsing request query: %v", err)
			log.Print(message)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(getJsonResponse(false, message, nil))
			return
		}
	}

	ctx := r.Context()
	data, err := db_readAll_branches(ctx, clause, args)

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_branches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	data, err := db_read_branches_ByPK(ctx, id)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_branches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var item Table_branches

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_branches(ctx, id, &item); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_branches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_delete_branches(ctx, id); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}

// college handler functions

func api_create_college(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_college

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_insert_college(ctx, &item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_college(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var clause string
	var args []any

	if len(queryValues) > 0 {
		clause, args, err = getQueryClauseArgs(queryValues, Map_college, "college")

		if err != nil {
			message := fmt.Sprintf("error while parsing request query: %v", err)
			log.Print(message)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(getJsonResponse(false, message, nil))
			return
		}
	}

	ctx := r.Context()
	data, err := db_readAll_college(ctx, clause, args)

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_college(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	data, err := db_read_college_ByPK(ctx, id)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_college(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var item Table_college

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_college(ctx, id, &item); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_college(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_delete_college(ctx, id); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}

// courses handler functions

func api_create_courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_courses

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_insert_courses(ctx, &item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var clause string
	var args []any

	if len(queryValues) > 0 {
		clause, args, err = getQueryClauseArgs(queryValues, Map_courses, "courses")

		if err != nil {
			message := fmt.Sprintf("error while parsing request query: %v", err)
			log.Print(message)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(getJsonResponse(false, message, nil))
			return
		}
	}

	ctx := r.Context()
	data, err := db_readAll_courses(ctx, clause, args)

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	data, err := db_read_courses_ByPK(ctx, id)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var item Table_courses

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_courses(ctx, id, &item); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_delete_courses(ctx, id); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}

// login handler functions

func api_create_login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_login

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := hashData([]*CustomNullString{&item.Column_password}, nil); err != nil {
		message := fmt.Sprintf("error while hashing fields: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_insert_login(ctx, &item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var clause string
	var args []any

	if len(queryValues) > 0 {
		clause, args, err = getQueryClauseArgs(queryValues, Map_login, "login")

		if err != nil {
			message := fmt.Sprintf("error while parsing request query: %v", err)
			log.Print(message)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(getJsonResponse(false, message, nil))
			return
		}
	}

	ctx := r.Context()
	data, err := db_readAll_login(ctx, clause, args)

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	data, err := db_read_login_ByPK(ctx, id)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	var item Table_login

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := hashData([]*CustomNullString{&item.Column_password}, nil); err != nil {
		message := fmt.Sprintf("error while hashing fields: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_login(ctx, id, &item); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		message := fmt.Sprintf("error while parsing request query: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	id := getPkParam(queryValues, "CustomNullString")
	if len(id) == 0 {
		message := "missing id param in request query"
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	ctx := r.Context()

	if err := db_delete_login(ctx, id); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}
