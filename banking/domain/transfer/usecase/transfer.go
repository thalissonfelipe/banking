package usecase

import (
	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
)

//go:generate moq -pkg usecase -out repository_mock.gen.go .. Repository
//go:generate moq -pkg usecase -out acc_usecase_mock.gen.go ../../account Usecase

type Transfer struct {
	repository     transfer.Repository
	accountUsecase account.Usecase
}

func NewTransferUsecase(repo transfer.Repository, accUsecase account.Usecase) *Transfer {
	return &Transfer{
		repository:     repo,
		accountUsecase: accUsecase,
	}
}
