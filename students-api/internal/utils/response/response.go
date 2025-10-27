package response

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(msg)
}
