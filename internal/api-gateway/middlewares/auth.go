package middlewares

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"net/http"
)

func NewAuthMiddleware() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			token := r.Header.Get("Authorization")
			userId := r.URL.Query().Get("userId")

			pairs := metadata.Pairs()
			if token != "" {
				pairs.Append("token", token)
			}

			if userId != "" {
				pairs.Append("userId", userId)
			}

			ctx := metadata.NewIncomingContext(r.Context(), pairs)
			next(w, r.WithContext(ctx), pathParams)
		}
	}
}
