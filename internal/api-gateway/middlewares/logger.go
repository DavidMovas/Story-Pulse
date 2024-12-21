package middlewares

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"net/http"
)

func NewLoggerMiddleware(logger *zap.SugaredLogger) runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			logger.Infow("Request", "Method", r.Method, "Url", r.RequestURI)
			next(w, r, pathParams)
		}
	}
}
