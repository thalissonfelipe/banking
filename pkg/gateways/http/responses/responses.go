package responses

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string
}

func SendError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func SendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
