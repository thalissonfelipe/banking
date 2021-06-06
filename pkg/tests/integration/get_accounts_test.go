package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
)

func TestIntegration_GetAccounts(t *testing.T) {
	defer dockertest.TruncateTables(context.Background(), pgDocker.DB)

	request := fakes.FakeRequest(http.MethodGet, server.URL+"/api/v1/accounts", nil)
	resp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	var body bytes.Buffer

	_, err = io.Copy(&body, resp.Body)
	require.NoError(t, err)

	t.Log(body.String())

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
