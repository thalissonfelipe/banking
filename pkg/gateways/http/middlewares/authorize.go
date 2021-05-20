package middlewares

import (
	"net/http"
	"strings"

	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

func AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			responses.SendError(w, http.StatusUnauthorized, errUnauthorized.Error())
			return
		}

		token := parts[1]
		if err := auth.IsValidToken(token); err != nil {
			responses.SendError(w, http.StatusUnauthorized, errUnauthorized.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
