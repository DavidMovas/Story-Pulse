package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"
)

func RequiredCookieMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("refresh_token")
			if err != nil && errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Refresh token not provided"))
				return
			}

			if cookie == nil || cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Refresh token not provided"))
				return
			}

			ctx := metadata.AppendToOutgoingContext(r.Context(), "refresh_token", cookie.Value)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RefreshTokenToCookieMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}
			next.ServeHTTP(res, r)

			if res.statusCode == http.StatusOK || res.statusCode == 0 {
				var resp map[string]interface{}
				if err := json.Unmarshal(res.body.Bytes(), &resp); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				if refreshToken, ok := resp["refreshToken"].(string); ok {
					http.SetCookie(w, &http.Cookie{
						Name:     "refresh_token",
						Value:    refreshToken,
						HttpOnly: true,
						Secure:   true,
						Path:     "/",
					})
					delete(resp, "refreshToken")
				} else {
					http.Error(w, "Refresh token not provided", http.StatusInternalServerError)
					return
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Length", strconv.Itoa(len(respData)))
				_, _ = w.Write(respData)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write(res.body.Bytes())
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) Write(p []byte) (int, error) {
	return r.body.Write(p)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}
