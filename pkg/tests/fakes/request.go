package fakes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/pkg/gateways/http/auth/schemes"
	"github.com/thalissonfelipe/banking/pkg/tests/testenv"
)

func FakeRequest(method string, path string, requestBody interface{}) *http.Request {
	var body io.Reader
	if requestBody != nil {
		b, err := json.Marshal(requestBody)
		if err != nil {
			log.Println("could not marshal request body")
		}
		body = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, path, body)
	return req
}

func FakeAuthorizedRequest(t *testing.T, method, path, cpf, secret string, requestBody interface{}) *http.Request {
	loginBody := schemes.LoginInput{CPF: cpf, Secret: secret}
	request := FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", loginBody)
	resp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	var respBody schemes.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)

	request = FakeRequest(method, path, requestBody)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", respBody.Token))

	return request
}
