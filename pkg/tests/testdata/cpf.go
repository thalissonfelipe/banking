package testdata

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
)

func GetValidCPF() vos.CPF {
	cpf, _ := vos.NewCPF(generateValidCPF())
	return cpf
}

// https://github.com/fnando/cpf/blob/master/src/index.ts
func generateValidCPF() string {
	var numbers string

	for i := 0; i < 9; i++ {
		numbers += fmt.Sprint(math.Floor(rand.Float64() * 9))
	}

	d := strconv.Itoa(verifierDigit(numbers))
	numbers += d

	d = strconv.Itoa(verifierDigit(numbers))
	numbers += d

	return numbers
}

func verifierDigit(digits string) int {
	numbersList := strings.Split(digits, "")

	var numbers []int

	for _, n := range numbersList {
		nn, _ := strconv.Atoi(n)
		numbers = append(numbers, nn)
	}

	modulus := len(numbersList) + 1

	var multiplied []int

	for i, n := range numbers {
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
