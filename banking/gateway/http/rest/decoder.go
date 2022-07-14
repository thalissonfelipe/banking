package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: add tests.
func DecodeRequestBody(r *http.Request, dest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	v, ok := dest.(interface{ IsValid() error })
	if !ok {
		return nil
	}

	if err := v.IsValid(); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}

	return nil
}
