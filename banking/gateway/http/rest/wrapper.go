package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/instrumentation/log"
)

func Wrap(handler func(r *http.Request) Response) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := log.FromContext(ctx)

		resp := handler(r)
		if resp.Error != nil {
			logger.Error("failed to handle request", zap.Error(resp.Error))
		}

		err := SendJSON(w, resp.Status, resp.Payload)
		if err != nil {
			logger.Error("failed to send response", zap.Error(err))
		}
	})
}

// TODO: add tests.
func SendJSON(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("encoding payload: %w", err)
	}

	return nil
}
