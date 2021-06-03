package account

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
)

type Handler struct {
	usecase account.Usecase
}

func NewHandler(r *chi.Mux, usecase account.Usecase) *Handler {
	handler := Handler{usecase: usecase}

	r.Route("/api/v1/accounts", func(r chi.Router) {
		r.Get("/", handler.ListAccounts)
		r.Post("/", handler.CreateAccount)
		r.Get("/{accountID}/balance", handler.GetAccountBalance)
	})

	return &handler
}
