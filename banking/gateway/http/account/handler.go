package account

import (
	"github.com/go-chi/chi/v5"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

//go:generate moq -pkg account -skip-ensure -out usecase_mock.gen.go ../../../domain/usecases Account:UsecaseMock

type Handler struct {
	usecase usecases.Account
}

func NewHandler(r chi.Router, usecase usecases.Account) *Handler {
	handler := Handler{usecase: usecase}

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", rest.Wrap(handler.ListAccounts))
		r.Post("/", rest.Wrap(handler.CreateAccount))
		r.Get("/{accountID}/balance", rest.Wrap(handler.GetAccountBalance))
	})

	return &handler
}
