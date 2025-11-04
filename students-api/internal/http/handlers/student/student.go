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
			//response.WriteJson(w, http.StatusBadRequest, err.Error())
			//response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(err))
			response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.CreateErrorResponse(err))
			return
		}

		// Validate the request received
		if err := validator.New().Struct(student); err != nil {
			validationErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidateAndCreateResponse(validationErr))
			return
		}

		lastInsertId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		slog.Info("user created successfully")

		response.WriteJson(w, http.StatusCreated, map[string]int64{"ID": lastInsertId})
	}
}
