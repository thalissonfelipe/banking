package fakes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
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
