package testdata

import "github.com/thalissonfelipe/banking/banking/domain/vos"

func GetValidSecret() vos.Secret {
	secret, _ := vos.NewSecret("aZ1234Ds")
	return secret
}
