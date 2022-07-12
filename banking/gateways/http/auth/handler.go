package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/banking/services/auth"
)

type Handler struct {
	authService *auth.Auth
}

func NewHandler(r chi.Router, authService *auth.Auth) *Handler {
	handler := &Handler{
		authService: authService,
	}

	r.Post("/login", handler.Login)

	return handler
}
