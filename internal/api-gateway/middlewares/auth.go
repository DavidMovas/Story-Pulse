package middlewares

import (
	"google.golang.org/grpc/metadata"
	"net/http"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			pairs := metadata.Pairs()
			if token != "" {
				pairs.Append("token", token)
			}

			ctx := metadata.NewIncomingContext(r.Context(), pairs)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthAndIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			userId := r.URL.Query().Get("id")

			pairs := metadata.Pairs()
			if token != "" {
				pairs.Append("token", token)
			}

			if userId != "" {
				pairs.Append("userId", userId)
			}

			ctx := metadata.NewIncomingContext(r.Context(), pairs)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
