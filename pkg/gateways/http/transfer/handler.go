package transfer

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/middlewares"
)

type Handler struct {
	usecase transfer.UseCase
}

func NewHandler(r *mux.Router, usecase transfer.UseCase) *Handler {
	handler := Handler{usecase: usecase}

	r.Handle("/transfers", middlewares.AuthorizeMiddleware(http.HandlerFunc(handler.ListTransfers))).Methods(http.MethodGet)
	r.Handle("/transfers", middlewares.AuthorizeMiddleware(http.HandlerFunc(handler.CreateTransfer))).Methods(http.MethodPost)

	return &handler
}
