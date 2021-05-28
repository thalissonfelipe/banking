package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCPF(t *testing.T) {
	testCases := []struct {
		name        string
		cpf         string
		expectedCPF string
		expectedErr error
	}{
		{"cpf must be valid #1", "648.446.967-93", "648.446.967-93", nil},
		{"cpf must be valid #2", "626.413.228-46", "626.413.228-46", nil},
		{"cpf must be valid #3", "871.957.260-37", "871.957.260-37", nil},
		{"cpf must be valid #4", "64844696793", "64844696793", nil},
		{"cpf must be valid #5", "62641322846", "62641322846", nil},
		{"cpf must be invalid #2", "000.000.000-00", "", errInvalidCPF},
		{"cpf must be invalid #1", "111.111.111-11", "", errInvalidCPF},
		{"cpf must be invalid #3", "222.222.222-33", "", errInvalidCPF},
		{"cpf must be invalid #4", "123.456.789-01", "", errInvalidCPF},
		{"cpf must be invalid #4", "00000000000", "", errInvalidCPF},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cpf, err := NewCPF(tt.cpf)

			assert.Equal(t, tt.expectedCPF, cpf.String())
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
