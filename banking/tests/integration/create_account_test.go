package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schemes"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	"github.com/thalissonfelipe/banking/banking/tests/testenv"
)

func TestIntegration_CreateAccount(t *testing.T) {
	testCases := []struct {
		name           string
		bodySetup      func() schemes.CreateAccountInput
		expectedStatus int
	}{
		{
			name: "should return status code 201 if account was created successfully",
			bodySetup: func() schemes.CreateAccountInput {
				body := schemes.CreateAccountInput{
					Name:   "Felipe",
					CPF:    testdata.GetValidCPF().String(),
					Secret: testdata.GetValidSecret().String(),
				}

				return body
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "should return status code 400 if name was not provided",
			bodySetup: func() schemes.CreateAccountInput {
				body := schemes.CreateAccountInput{
					Name:   "",
					CPF:    testdata.GetValidCPF().String(),
					Secret: testdata.GetValidSecret().String(),
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			bodySetup: func() schemes.CreateAccountInput {
				body := schemes.CreateAccountInput{
					Name:   "Felipe",
					CPF:    "",
					Secret: testdata.GetValidSecret().String(),
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if secret was not provided",
			bodySetup: func() schemes.CreateAccountInput {
				body := schemes.CreateAccountInput{
					Name:   "Felipe",
					CPF:    testdata.GetValidCPF().String(),
					Secret: "",
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if cpf provided was not valid",
			bodySetup: func() schemes.CreateAccountInput {
				body := schemes.CreateAccountInput{
					Name:   "Felipe",
					CPF:    "123.456.789-00",
					Secret: testdata.GetValidSecret().String(),
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if secret provided was not valid",
			bodySetup: func() schemes.CreateAccountInput {
				body := schemes.CreateAccountInput{
					Name:   "Felipe",
					CPF:    testdata.GetValidCPF().String(),
					Secret: "12345678",
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 409 if account already exists",
			bodySetup: func() schemes.CreateAccountInput {
				acc := createAccount(t, testdata.GetValidCPF(), testdata.GetValidSecret(), 0)

				body := schemes.CreateAccountInput{
					Name:   "Felipe",
					CPF:    acc.CPF.String(),
					Secret: testdata.GetValidSecret().String(),
				}

				return body
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := fakes.FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/accounts", tt.bodySetup())
			resp, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			defer resp.Body.Close()

			var body bytes.Buffer

			_, err = io.Copy(&body, resp.Body)
			require.NoError(t, err)

			t.Log(body.String())

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			dockertest.TruncateTables(context.Background(), testenv.DB)
		})
	}
}
