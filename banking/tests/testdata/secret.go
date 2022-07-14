package testdata

import "github.com/thalissonfelipe/banking/banking/domain/vos"

func Secret() vos.Secret {
	secret, _ := vos.NewSecret("aZ1234Ds")
	return secret
}
