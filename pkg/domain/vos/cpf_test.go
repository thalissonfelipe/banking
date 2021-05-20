package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPFIsValid(t *testing.T) {
	testCases := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"cpf must be valid #1", "648.446.967-93", true},
		{"cpf must be valid #2", "626.413.228-46", true},
		{"cpf must be valid #3", "871.957.260-37", true},
		{"cpf must be valid #4", "64844696793", true},
		{"cpf must be valid #5", "62641322846", true},
		{"cpf must be invalid #2", "000.000.000-00", false},
		{"cpf must be invalid #1", "111.111.111-11", false},
		{"cpf must be invalid #3", "222.222.222-33", false},
		{"cpf must be invalid #4", "123.456.789-01", false},
		{"cpf must be invalid #4", "00000000000", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cpf := NewCPF(tt.cpf)
			result := cpf.IsValid()

			assert.Equal(t, tt.expected, result)
		})
	}
}
