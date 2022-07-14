package middlewares

import (
	"net/http"
	"strings"

	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
)

func AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const partsSize = 2
		authHeader := r.Header.Get("Authorization")

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != partsSize {
			rest.SendJSON(w, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		token := parts[1]
		if err := jwt.IsTokenValid(token); err != nil {
			rest.SendJSON(w, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
