package rest

import (
	"encoding/json"
	"net/http"
)

func Wrap(handler func(r *http.Request) Response) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := handler(r)
		if resp.Error != nil {
			// TODO: log errors
		}

		err := SendJSON(w, resp.Status, resp.Payload)
		if err != nil {
			// TODO: log errors
		}
	})
}

// TODO: add tests.
func SendJSON(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}
