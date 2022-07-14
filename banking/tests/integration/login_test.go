package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateway/http/auth/schema"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	"github.com/thalissonfelipe/banking/banking/tests/testenv"
)

func TestIntegration_Login(t *testing.T) {
	testCases := []struct {
		name           string
		requestSetup   func() *http.Request
		expectedStatus int
	}{
		{
			name: "should log in successfully",
			requestSetup: func() *http.Request {
				secret := testdata.Secret()
				acc := createAccount(t, testdata.CPF(), secret, 0)

				reqBody := schema.LoginInput{CPF: acc.CPF.String(), Secret: secret.String()}
				request := fakes.FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", reqBody)

				return request
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "should return status code 400 if account does not exist",
			requestSetup: func() *http.Request {
				reqBody := schema.LoginInput{
					CPF:    testdata.CPF().String(),
					Secret: testdata.Secret().String(),
				}
				request := fakes.FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if secret is not valid",
			requestSetup: func() *http.Request {
				acc := createAccount(t, testdata.CPF(), testdata.Secret(), 0)

				reqBody := schema.LoginInput{CPF: acc.CPF.String(), Secret: "12345678"}
				request := fakes.FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", reqBody)

				return request
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(tt.requestSetup())
			require.NoError(t, err)

			defer resp.Body.Close()

			var body bytes.Buffer

			_, err = io.Copy(&body, resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			dockertest.TruncateTables(context.Background(), testenv.DB)
		})
	}
}
