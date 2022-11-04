package rest

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendJSON(t *testing.T) {
	t.Parallel()

	const (
		wantStatus      = http.StatusOK
		wantContentType = "application/json"
		wantPayload     = "payload"
	)

	rec := httptest.NewRecorder()

	err := SendJSON(rec, wantStatus, wantPayload)
	require.NoError(t, err)

	gotBody := regexp.MustCompile(`[\W]`).ReplaceAllString(rec.Body.String(), "")

	assert.Equal(t, wantStatus, rec.Code)
	assert.Equal(t, wantContentType, rec.Header().Get("Content-Type"))
	assert.Equal(t, wantPayload, gotBody)
}
