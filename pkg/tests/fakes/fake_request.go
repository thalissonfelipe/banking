package fakes

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
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
