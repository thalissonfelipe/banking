package account

import (
	"github.com/gorilla/mux"
	"github.com/thalissonfelipe/banking/pkg/domain/account"
)

type Handler struct {
	usecase account.UseCase
}

func NewHandler(r *mux.Router, usecase account.UseCase) *Handler {
	handler := Handler{usecase: usecase}

	r.HandleFunc("/accounts/", handler.ListAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}/balance/", handler.GetAccountBalance).Methods("GET")

	return &handler
}
