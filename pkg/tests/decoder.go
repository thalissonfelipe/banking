package tests

import (
	"bytes"
	"encoding/json"

	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
)

type Decoder interface {
	Decode(body *bytes.Buffer) interface{}
}

type ErrorMessageDecoder struct{}

func (ErrorMessageDecoder) Decode(body *bytes.Buffer) interface{} {
	var errMessage responses.ErrorResponse
	json.NewDecoder(body).Decode(&errMessage)
	return errMessage
}
