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
	"github.com/thalissonfelipe/banking/pkg/tests/testenv"
)

func TestIntegration_GetAccounts(t *testing.T) {
	defer dockertest.TruncateTables(context.Background(), testenv.DB)

	request := fakes.FakeRequest(http.MethodGet, testenv.ServerURL+"/api/v1/accounts", nil)
	resp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	defer resp.Body.Close()

	var body bytes.Buffer

	_, err = io.Copy(&body, resp.Body)
	require.NoError(t, err)

	t.Log(body.String())

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
