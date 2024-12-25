package middlewares

import "net/http"

func Wrap(first, second http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		first.ServeHTTP(w, r)
		second.ServeHTTP(w, r)
	})
}
