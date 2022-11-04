package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeRequestBody(t *testing.T) {
	t.Parallel()

	type requestBody struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name        string
		body        []byte
		wantMessage string
		wantErr     bool
	}{
		{
			name:        "should decode request body successfully",
			body:        json.RawMessage(`{"message": "test"}`),
			wantMessage: "test",
			wantErr:     false,
		},
		{
			name:        "should return an error if body is invalid",
			body:        json.RawMessage(`{"message": "test"`),
			wantMessage: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(tt.body))

			var body requestBody

			err := DecodeRequestBody(r, &body)
			assert.Equal(t, tt.wantErr, err != nil)

			assert.Equal(t, tt.wantMessage, body.Message)
		})
	}
}
