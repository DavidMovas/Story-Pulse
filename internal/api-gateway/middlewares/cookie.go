package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"google.golang.org/grpc/metadata"
	"net/http"
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

func RefreshTokenToCookieMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}
			next.ServeHTTP(res, r)

			if res.statusCode == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(res.body.Bytes(), &resp); err != nil {
					if refreshToken, ok := resp["refresh_token"]; ok {
						http.SetCookie(res, &http.Cookie{
							Name:     "refresh_token",
							Value:    refreshToken.(string),
							HttpOnly: true,
							Path:     "/",
						})
					}
				}
			}

			w.WriteHeader(res.statusCode)
			_, _ = w.Write(res.body.Bytes())
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	return r.body.Write(data)
}
