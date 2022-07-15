package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/instrumentation/log"
)

type requestIDKey struct{}

func contextWithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestIDKey{}, uuid.NewString())
}

func requestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		log.FromContext(ctx).Warn("missing request id in context")
	}

	return requestID
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := contextWithRequestID(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDToLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := requestIDFromContext(ctx)
		logger := log.FromContext(ctx)
		logger = logger.With(zap.String("request_id", requestID))
		ctx = log.NewContext(ctx, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
