package fakes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateway/http/auth/schemes"
	"github.com/thalissonfelipe/banking/banking/tests/testenv"
)

func FakeRequest(method, path string, requestBody interface{}) *http.Request {
	var body io.Reader

	if requestBody != nil {
		b, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatalf("could not marshal request body: %v", err)
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, path, body)
	if err != nil {
		log.Fatalf("could not create a new request: %v", err)
	}

	return req
}

func FakeAuthorizedRequest(t *testing.T, method, path, cpf, secret string, requestBody interface{}) *http.Request {
	loginBody := schemes.LoginInput{CPF: cpf, Secret: secret}
	request := FakeRequest(http.MethodPost, testenv.ServerURL+"/api/v1/login", loginBody)
	resp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	defer resp.Body.Close()

	var respBody schemes.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)

	request = FakeRequest(method, path, requestBody)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", respBody.Token))

	return request
}
