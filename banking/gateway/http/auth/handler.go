package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/banking/services"
)

//go:generate moq -pkg auth -out auth_mock.gen.go ../../../services Auth

type Handler struct {
	authService services.Auth
}

func NewHandler(r chi.Router, authService services.Auth) *Handler {
	handler := &Handler{
		authService: authService,
	}

	r.Post("/login", handler.Login)

	return handler
}
