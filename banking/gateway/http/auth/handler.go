package auth

import (
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
)

//go:generate moq -pkg auth -out usecase_mock.gen.go ../../../domain/usecases Auth:UsecaseMock

type Handler struct {
	usecase usecases.Auth
}

func NewHandler(usecase usecases.Auth) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
