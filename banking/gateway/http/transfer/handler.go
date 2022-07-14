package transfer

import (
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg transfer -skip-ensure -out usecase_mock.gen.go ../../../domain/usecases Transfer:UsecaseMock

type Handler struct {
	usecase usecases.Transfer
}

func NewHandler(usecase usecases.Transfer) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
