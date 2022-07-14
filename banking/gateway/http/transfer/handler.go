package transfer

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/middlewares"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

//go:generate moq -pkg transfer -skip-ensure -out usecase_mock.gen.go ../../../domain/usecases Transfer:UsecaseMock

type Handler struct {
	usecase usecases.Transfer
}

func NewHandler(r chi.Router, usecase usecases.Transfer) *Handler {
	handler := Handler{usecase: usecase}

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthorizeMiddleware)
		r.Route("/transfers", func(r chi.Router) {
			r.Get("/", rest.Wrap(handler.ListTransfers))
			r.Post("/", rest.Wrap(handler.PerformTransfer))
		})
	})

	return &handler
}
