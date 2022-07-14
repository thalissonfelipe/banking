package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDKey struct{}

func contextWithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestIDKey{}, uuid.NewString())
}

func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey{}).(string)
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := contextWithRequestID(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
