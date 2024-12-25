package middlewares

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func NewCookieMiddleware() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			cookie, err := r.Cookie("refresh_token")
			if err != nil && errors.Is(err, http.ErrNoCookie) {
				next(w, r, pathParams)
			}

			ctx := r.Context()
			if cookie != nil && cookie.Value != "" {
				ctx = context.WithValue(ctx, "refresh_token", cookie.Value)
			}

			next(w, r.WithContext(ctx), pathParams)
		}
	}
}
