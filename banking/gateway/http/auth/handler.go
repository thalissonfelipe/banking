package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

//go:generate moq -pkg auth -out usecase_mock.gen.go ../../../domain/usecases Auth:UsecaseMock

type Handler struct {
	usecase usecases.Auth
}

func NewHandler(r chi.Router, usecase usecases.Auth) *Handler {
	handler := &Handler{
		usecase: usecase,
	}

	r.Post("/login", rest.Wrap(handler.Login))

	return handler
}
