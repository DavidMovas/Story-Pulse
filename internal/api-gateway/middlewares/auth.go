package middlewares

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"net/http"
)

func AuthHTTPMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("AUTH MIDDLEWARE !")
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
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
