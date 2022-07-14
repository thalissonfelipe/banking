package integration

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/tests/dockertest"
	"github.com/thalissonfelipe/banking/banking/tests/fakes"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	"github.com/thalissonfelipe/banking/banking/tests/testenv"
)

func TestIntegration_CreateAccount(t *testing.T) {
	testCases := []struct {
		name           string
		bodySetup      func() schema.CreateAccountInput
		expectedStatus int
	}{
		{
			name: "should return status code 201 if account was created successfully",
			bodySetup: func() schema.CreateAccountInput {
				body := schema.CreateAccountInput{
					Name:   "Felipe",
					CPF:    testdata.CPF().String(),
					Secret: testdata.Secret().String(),
				}

				return body
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "should return status code 400 if name was not provided",
			bodySetup: func() schema.CreateAccountInput {
				body := schema.CreateAccountInput{
					Name:   "",
					CPF:    testdata.CPF().String(),
					Secret: testdata.Secret().String(),
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			bodySetup: func() schema.CreateAccountInput {
				body := schema.CreateAccountInput{
					Name:   "Felipe",
					CPF:    "",
					Secret: testdata.Secret().String(),
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if secret was not provided",
			bodySetup: func() schema.CreateAccountInput {
				body := schema.CreateAccountInput{
					Name:   "Felipe",
					CPF:    testdata.CPF().String(),
					Secret: "",
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if cpf provided was not valid",
			bodySetup: func() schema.CreateAccountInput {
				body := schema.CreateAccountInput{
					Name:   "Felipe",
					CPF:    "123.456.789-00",
					Secret: testdata.Secret().String(),
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 400 if secret provided was not valid",
			bodySetup: func() schema.CreateAccountInput {
				body := schema.CreateAccountInput{
					Name:   "Felipe",
					CPF:    testdata.CPF().String(),
					Secret: "12345678",
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should return status code 409 if account already exists",
			bodySetup: func() schema.CreateAccountInput {
				acc := createAccount(t, testdata.CPF(), testdata.Secret(), 0)

				body := schema.CreateAccountInput{
					Name:   "Felipe",
					CPF:    acc.CPF.String(),
					Secret: testdata.Secret().String(),
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

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			dockertest.TruncateTables(context.Background(), testenv.DB)
		})
	}
}
