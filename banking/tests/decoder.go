package tests

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

type Decoder interface {
	Decode(t *testing.T, body *bytes.Buffer) interface{}
}

type ErrorMessageDecoder struct{}

func (ErrorMessageDecoder) Decode(t *testing.T, body *bytes.Buffer) interface{} {
	var errMessage rest.ErrorResponse

	err := json.NewDecoder(body).Decode(&errMessage)
	require.NoError(t, err)

	return errMessage
}
