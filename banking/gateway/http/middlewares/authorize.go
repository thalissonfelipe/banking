package middlewares

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
	"github.com/thalissonfelipe/banking/banking/instrumentation/log"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const partsSize = 2
		authHeader := r.Header.Get("Authorization")

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != partsSize {
			if err := rest.SendJSON(w, http.StatusUnauthorized, rest.ErrUnauthorized); err != nil {
				log.FromContext(r.Context()).Error("failed to send json", zap.Error(err))
			}

			return
		}

		token := parts[1]
		if err := jwt.IsTokenValid(token); err != nil {
			if err := rest.SendJSON(w, http.StatusUnauthorized, rest.ErrUnauthorized); err != nil {
				log.FromContext(r.Context()).Error("failed to send json", zap.Error(err))
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}
