package testdata

import (
	"math/rand"
	"strconv"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func CPF() vos.CPF {
	cpf, _ := vos.NewCPF(generateCPF())
	return cpf
}

// https://github.com/fnando/cpf/blob/master/src/index.ts
func generateCPF() string {
	numbers := rand.Perm(9)
	numbers = append(numbers, verifierDigit(numbers))
	numbers = append(numbers, verifierDigit(numbers))

	var cpfString string

	for _, n := range numbers {
		cpfString += strconv.Itoa(n)
	}

	return cpfString
}

func verifierDigit(digits []int) int {
	modulus := len(digits) + 1

	var multiplied []int

	for i, n := range digits {
		multiplied = append(multiplied, (n * (modulus - i)))
	}

	sum := 0

	for _, n := range multiplied {
		sum += n
	}

	mod := sum % 11
	if mod < 2 {
		return 0
	}

	return 11 - mod
}
