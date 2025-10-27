package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/yash-codes/students-api/internal/types"
	"github.com/yash-codes/students-api/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("Welcome to students-api"))
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
