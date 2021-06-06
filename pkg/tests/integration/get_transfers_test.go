package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/tests/dockertest"
	"github.com/thalissonfelipe/banking/pkg/tests/fakes"
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestIntegration_GetTransfers(t *testing.T) {
	testCases := []struct {
		name           string
		requestSetup   func() *http.Request
		expectedStatus int
	}{
		{
			name: "should return list of transfers sucessfully",
			requestSetup: func() *http.Request {
				// TODO: refactor login to avoid duplicated code
				type requestBody struct {
					CPF    string `json:"cpf"`
					Secret string `json:"secret"`
				}

				type responseBody struct {
					Token string `json:"token"`
				}

				secret := testdata.GetValidSecret()
				acc := createAccount(t, testdata.GetValidCPF(), secret)

				reqBody := requestBody{CPF: acc.CPF.String(), Secret: secret.String()}
				request := fakes.FakeRequest(http.MethodPost, server.URL+"/api/v1/login", reqBody)
				resp, err := http.DefaultClient.Do(request)
				require.NoError(t, err)

				var respBody responseBody
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)

				request = fakes.FakeRequest(http.MethodGet, server.URL+"/api/v1/transfers", nil)
				request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", respBody.Token))

				return request
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "should return status code 401 when user is not authorized",
			requestSetup: func() *http.Request {
				request := fakes.FakeRequest(http.MethodGet, server.URL+"/api/v1/transfers", nil)

				return request
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.requestSetup())
			require.NoError(t, err)

			var body bytes.Buffer

			_, err = io.Copy(&body, resp.Body)
			require.NoError(t, err)

			t.Log(body.String())

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			dockertest.TruncateTables(context.Background(), pgDocker.DB)
		})
	}
}
