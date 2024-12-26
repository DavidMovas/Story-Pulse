package middlewares

import (
	"context"
	"errors"
	"net/http"
)

func RequiredCookieMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("refresh_token")
			if err != nil && errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Unauthorized"))
				return
			}

			if cookie == nil || cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Unauthorized"))
				return
			}

			ctx := context.WithValue(r.Context(), "refresh_token", cookie.Value)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
