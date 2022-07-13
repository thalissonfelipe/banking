package transfer

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/gateway/http/middlewares"
)

//go:generate moq -pkg transfer -skip-ensure -out usecase_mock.gen.go ../../../domain/transfer Usecase

type Handler struct {
	usecase transfer.Usecase
}

func NewHandler(r chi.Router, usecase transfer.Usecase) *Handler {
	handler := Handler{usecase: usecase}

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthorizeMiddleware)
		r.Route("/transfers", func(r chi.Router) {
			r.Get("/", handler.ListTransfers)
			r.Post("/", handler.PerformTransfer)
		})
	})

	return &handler
}
