package transfer

import (
	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg transfer -out repository_mock.gen.go ../../entity TransferRepository:RepositoryMock
//go:generate moq -pkg transfer -out acc_usecase_mock.gen.go .. Account:UsecaseMock

var _ usecases.Transfer = (*Usecase)(nil)

type Usecase struct {
	repository     entity.TransferRepository
	accountUsecase usecases.Account
}

func NewUsecase(repo entity.TransferRepository, accUsecase usecases.Account) *Usecase {
	return &Usecase{
		repository:     repo,
		accountUsecase: accUsecase,
	}
}
