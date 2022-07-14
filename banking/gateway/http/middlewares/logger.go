package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/instrumentation/log"
)

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.NewContext(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
