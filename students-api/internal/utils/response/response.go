package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(msg)
}

func CreateErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidateAndCreateResponse(errs validator.ValidationErrors) ErrorResponse {

	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return ErrorResponse{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
