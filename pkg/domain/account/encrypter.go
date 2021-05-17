package account

type Encrypter interface {
	Hash(secret string) ([]byte, error)
}
