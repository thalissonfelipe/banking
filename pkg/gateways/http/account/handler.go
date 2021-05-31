package account

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
)

type Handler struct {
	usecase account.Usecase
}

func NewHandler(r *mux.Router, usecase account.Usecase) *Handler {
	handler := Handler{usecase: usecase}

	r.HandleFunc("/accounts", handler.ListAccounts).Methods(http.MethodGet)
	r.HandleFunc("/accounts", handler.CreateAccount).Methods(http.MethodPost)
	r.HandleFunc("/accounts/{id}/balance", handler.GetAccountBalance).Methods(http.MethodGet)

	return &handler
}
