package transfer

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
)

type Handler struct {
	usecase transfer.UseCase
}

func NewHandler(r *mux.Router, usecase transfer.UseCase) *Handler {
	handler := Handler{usecase: usecase}

	r.HandleFunc("/transfers", handler.ListTransfers).Methods(http.MethodGet)

	return &handler
}
