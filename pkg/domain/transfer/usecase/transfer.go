package usecase

import (
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
)

type Transfer struct {
	repository     transfer.Repository
	accountUseCase account.UseCase
}

func NewTransferUsecase(repo transfer.Repository, accUseCase account.UseCase) *Transfer {
	return &Transfer{
		repository:     repo,
		accountUseCase: accUseCase,
	}
}
