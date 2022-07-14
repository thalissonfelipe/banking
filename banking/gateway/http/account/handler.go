package account

import (
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg account -skip-ensure -out usecase_mock.gen.go ../../../domain/usecases Account:UsecaseMock

type Handler struct {
	usecase usecases.Account
}

func NewHandler(usecase usecases.Account) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
