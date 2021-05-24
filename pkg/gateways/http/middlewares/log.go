package middlewares

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type logRecord struct {
	http.ResponseWriter
	statusCode int
}

func (r *logRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

func (r *logRecord) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record := &logRecord{
			ResponseWriter: w,
		}

		log.WithFields(log.Fields{
			"method": r.Method,
			"route":  r.URL.String(),
		}).Info("request")

		next.ServeHTTP(record, r)

		log.WithFields(log.Fields{
			"method": r.Method,
			"route":  r.URL.String(),
			"status": record.statusCode,
		}).Info("response")
	})
}
