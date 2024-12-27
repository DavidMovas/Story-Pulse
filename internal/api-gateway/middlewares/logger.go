package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewLoggerMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			logger.Infow("Request", "Method", r.Method, "Url", r.RequestURI)
			next.ServeHTTP(w, r)
			logger.Infow("Request", "Time", time.Since(now))
		})
	}
}
