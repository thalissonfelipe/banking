package testdata

import "github.com/thalissonfelipe/banking/pkg/domain/vos"

func GetValidSecret() vos.Secret {
	secret, _ := vos.NewSecret("aZ1234Ds")
	return secret
}
