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
	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
)

func TestIntegration_CreateAccount(t *testing.T) {
	type requestBody struct {
		Name   string `json:"name"`
		CPF    string `json:"cpf"`
		Secret string `json:"secret"`
	}

	testCases := []struct {
		name           string
		bodySetup      func() requestBody
		expectedStatus int
	}{
		{
			name: "should return status code 201 if account was created successfully",
			bodySetup: func() requestBody {
				body := requestBody{
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
			bodySetup: func() requestBody {
				body := requestBody{
					Name:   "",
					CPF:    testdata.GetValidCPF().String(),
					Secret: testdata.GetValidSecret().String(),
				}

				return body
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			bodySetup: func() requestBody {
				body := requestBody{
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
			bodySetup: func() requestBody {
				body := requestBody{
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
			bodySetup: func() requestBody {
				body := requestBody{
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
			bodySetup: func() requestBody {
				body := requestBody{
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
			bodySetup: func() requestBody {
				acc := createAccount(t, testdata.GetValidCPF(), testdata.GetValidSecret())

				body := requestBody{
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
			request := fakes.FakeRequest(http.MethodPost, server.URL+"/api/v1/accounts", tt.bodySetup())
			resp, err := http.DefaultClient.Do(request)
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
