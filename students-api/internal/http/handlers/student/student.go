package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/yash-codes/students-api/internal/storage"
	"github.com/yash-codes/students-api/internal/types"
	"github.com/yash-codes/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("Welcome to students-api"))
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			slog.Error("Error while reading the request body", "Error", err)
			//response.WriteJson(w, http.StatusBadRequest, err.Error())
			//response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(err))
			response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			slog.Error("Error while reading the request body", "Error", err)
			response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(err))
			return
		}

		// Validate the request received
		if err := validator.New().Struct(student); err != nil {
			slog.Error("Error while validating request", "Error", err)
			validationErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidateAndCreateResponse(validationErr))
			return
		}

		lastInsertId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			slog.Error("something went wrong", "Error", err)
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		slog.Info("user created successfully")

		response.WriteJson(w, http.StatusCreated, map[string]int64{"ID": lastInsertId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student by", "ID", id)

		if id == "" {
			slog.Error("found no id to be deleted")
			response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(fmt.Errorf("id not passed")))
			return
		}

		student, err := storage.GetStudentById(id)
		if err != nil {
			slog.Error("something went wrong", "Error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.CreateErrorResponse(err))
			return
		}

		slog.Info("found student successfully", "student", student)
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		students, err := storage.GetAllStudentsList()
		if err != nil {
			slog.Error("something went wrong", "Error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.CreateErrorResponse(err))
			return
		}

		slog.Info("found all students successfully")
		response.WriteJson(w, http.StatusOK, students)
	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			slog.Error("found no id to be deleted")
			response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(fmt.Errorf("found no valid id to be deleted")))
			return
		}

		affectedId, err := storage.DeleteStudentById(id)
		if err != nil {
			slog.Error("something went wrong", "Error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.CreateErrorResponse(err))
			return
		}

		slog.Info("deleted student successfully")
		response.WriteJson(w, http.StatusOK, map[string]int64{"ID": affectedId})
	}
}
