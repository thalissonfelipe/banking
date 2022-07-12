package encrypter

type Encrypter interface {
	Hash(secret string) ([]byte, error)
	CompareHashAndSecret(hashedSecret, secret []byte) error
}
