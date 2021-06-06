package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestIntegration_Login(t *testing.T) {
	type requestBody struct {
		CPF    string `json:"cpf"`
		Secret string `json:"secret"`
	}

	type responseBody struct {
		Token string `json:"token"`
	}

	secret := testdata.GetValidSecret()
	acc := createAccount(t, testdata.GetValidCPF(), secret, 0)

	reqBody := requestBody{CPF: acc.CPF.String(), Secret: secret.String()}

	request := fakes.FakeRequest(http.MethodPost, server.URL+"/api/v1/login", reqBody)
	resp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	var respBody responseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)

	t.Logf("%v", respBody)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	dockertest.TruncateTables(context.Background(), pgDocker.DB)
}
