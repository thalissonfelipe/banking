package auth

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

type Handler struct {
	authService *auth.Auth
}

func NewHandler(r *mux.Router, authService *auth.Auth) *Handler {
	handler := &Handler{
		authService: authService,
	}

	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)

	return handler
}
