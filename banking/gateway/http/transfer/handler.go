package transfer

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/gateway/http/middlewares"
)

type Handler struct {
	usecase transfer.UseCase
}

func NewHandler(r chi.Router, usecase transfer.UseCase) *Handler {
	handler := Handler{usecase: usecase}

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthorizeMiddleware)
		r.Route("/transfers", func(r chi.Router) {
			r.Get("/", handler.ListTransfers)
			r.Post("/", handler.CreateTransfer)
		})
	})

	return &handler
}
