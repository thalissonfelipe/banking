package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
)

func DecodeRequestBody(r *http.Request, dest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		return responses.ErrInvalidJSON
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