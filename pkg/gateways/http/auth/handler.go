package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

type Handler struct {
	authService *auth.Auth
}

func NewHandler(r *chi.Mux, authService *auth.Auth) *Handler {
	handler := &Handler{
		authService: authService,
	}

	r.Post("/api/v1/login", handler.Login)

	return handler
}
