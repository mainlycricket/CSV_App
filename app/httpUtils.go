package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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

	// TypeTest routes
	http.HandleFunc("POST /TypeTest", api_create_TypeTest)
	http.HandleFunc("GET /TypeTest", api_getAll_TypeTest)
	http.HandleFunc("GET /TypeTest/{id}", api_getByPk_TypeTest)
	http.HandleFunc("PUT /TypeTest/{id}", api_update_TypeTest)
	http.HandleFunc("DELETE /TypeTest/{id}", api_delete_TypeTest)

	// branches routes
	http.HandleFunc("POST /branches", api_create_branches)
	http.HandleFunc("GET /branches", api_getAll_branches)
	http.HandleFunc("GET /branches/{id}", api_getByPk_branches)
	http.HandleFunc("PUT /branches/{id}", api_update_branches)
	http.HandleFunc("DELETE /branches/{id}", api_delete_branches)

	// courses routes
	http.HandleFunc("POST /courses", api_create_courses)
	http.HandleFunc("GET /courses", api_getAll_courses)
	http.HandleFunc("GET /courses/{id}", api_getByPk_courses)
	http.HandleFunc("PUT /courses/{id}", api_update_courses)
	http.HandleFunc("DELETE /courses/{id}", api_delete_courses)

	// empty routes
	http.HandleFunc("POST /empty", api_create_empty)
	http.HandleFunc("GET /empty", api_getAll_empty)
	http.HandleFunc("GET /empty/{id}", api_getByPk_empty)
	http.HandleFunc("PUT /empty/{id}", api_update_empty)
	http.HandleFunc("DELETE /empty/{id}", api_delete_empty)

	// students routes
	http.HandleFunc("POST /students", api_create_students)
	http.HandleFunc("GET /students", api_getAll_students)
	http.HandleFunc("GET /students/{id}", api_getByPk_students)
	http.HandleFunc("PUT /students/{id}", api_update_students)
	http.HandleFunc("DELETE /students/{id}", api_delete_students)

	// subjects routes
	http.HandleFunc("POST /subjects", api_create_subjects)
	http.HandleFunc("GET /subjects", api_getAll_subjects)
	http.HandleFunc("GET /subjects/{id}", api_getByPk_subjects)
	http.HandleFunc("PUT /subjects/{id}", api_update_subjects)
	http.HandleFunc("DELETE /subjects/{id}", api_delete_subjects)

	s := &http.Server{
		Addr: ":8080",
	}

	return s
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

	if err := db_insert_TypeTest(&item); err != nil {
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

	data, err := db_readAll_TypeTest()

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast id", nil))
		return
	}

	data, err := db_read_TypeTest_ByPK(pk)

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast", nil))
		return
	}

	var item Table_TypeTest

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_TypeTest(&item, pk); err != nil {
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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type caste", nil))
		return
	}

	if err := db_delete_TypeTest(pk); err != nil {
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

	if err := db_insert_branches(&item); err != nil {
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

	data, err := db_readAll_branches()

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast id", nil))
		return
	}

	data, err := db_read_branches_ByPK(pk)

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast", nil))
		return
	}

	var item Table_branches

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_branches(&item, pk); err != nil {
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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type caste", nil))
		return
	}

	if err := db_delete_branches(pk); err != nil {
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

	if err := db_insert_courses(&item); err != nil {
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

	data, err := db_readAll_courses()

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast id", nil))
		return
	}

	data, err := db_read_courses_ByPK(pk)

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast", nil))
		return
	}

	var item Table_courses

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_courses(&item, pk); err != nil {
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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type caste", nil))
		return
	}

	if err := db_delete_courses(pk); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}

// empty handler functions

func api_create_empty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Table_empty

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_insert_empty(&item); err != nil {
		message := fmt.Sprintf("error while creating: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(getJsonResponse(true, "created successfully", nil))
}

func api_getAll_empty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := db_readAll_empty()

	if err != nil {
		message := fmt.Sprintf("error while reading: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "data fetched successfully", data))
}

func api_getByPk_empty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast id", nil))
		return
	}

	data, err := db_read_empty_ByPK(pk)

	if err != nil {
		message := fmt.Sprintf("error while reading data: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "found data", data))
}

func api_update_empty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast", nil))
		return
	}

	var item Table_empty

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_empty(&item, pk); err != nil {
		message := fmt.Sprintf("error while updating : %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "updated successfully", nil))
}

func api_delete_empty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type caste", nil))
		return
	}

	if err := db_delete_empty(pk); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
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

	if err := db_insert_students(&item); err != nil {
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

	data, err := db_readAll_students()

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast id", nil))
		return
	}

	data, err := db_read_students_ByPK(pk)

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast", nil))
		return
	}

	var item Table_students

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_students(&item, pk); err != nil {
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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type caste", nil))
		return
	}

	if err := db_delete_students(pk); err != nil {
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

	if err := db_insert_subjects(&item); err != nil {
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

	data, err := db_readAll_subjects()

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast id", nil))
		return
	}

	data, err := db_read_subjects_ByPK(pk)

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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type cast", nil))
		return
	}

	var item Table_subjects

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("error while reading request body: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	if err := db_update_subjects(&item, pk); err != nil {
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

	id := r.PathValue("id")
	pk, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, "failed to type caste", nil))
		return
	}

	if err := db_delete_subjects(pk); err != nil {
		message := fmt.Sprintf("error while deleting: %v", err)
		log.Print(message)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(getJsonResponse(false, message, nil))
		return
	}

	w.Write(getJsonResponse(true, "deleted successfully", nil))
}
