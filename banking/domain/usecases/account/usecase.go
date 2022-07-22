package account

import (
	"github.com/thalissonfelipe/banking/banking/domain/encrypter"
	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg account -out repository_mock.gen.go ../../entity AccountRepository:RepositoryMock
//go:generate moq -pkg account -out encrypter_mock.gen.go ../../encrypter Encrypter

var _ usecases.Account = (*Usecase)(nil)

type Usecase struct {
	repository entity.AccountRepository
	encrypter  encrypter.Encrypter
}

func NewUsecase(repo entity.AccountRepository, enc encrypter.Encrypter) *Usecase {
	return &Usecase{
		repository: repo,
		encrypter:  enc,
	}
}
