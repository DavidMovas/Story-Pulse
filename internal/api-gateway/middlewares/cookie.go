package middlewares

import (
	"bytes"
	"errors"
	"github.com/goccy/go-json"
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

func RefreshTokenToCookieMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}
			next.ServeHTTP(res, r)

			if res.statusCode == http.StatusOK || res.statusCode == 0 {
				var resp map[string]interface{}
				if err := json.Unmarshal(res.body.Bytes(), &resp); err == nil {
					if refreshToken, ok := resp["refreshToken"].(string); ok {
						http.SetCookie(res.ResponseWriter, &http.Cookie{
							Name:     "refresh_token",
							Value:    refreshToken,
							HttpOnly: true,
							Secure:   true,
							Path:     "/",
						})

						delete(resp, "refreshToken")
					}

					respData, err := json.Marshal(resp)
					if err != nil {
						http.Error(res.ResponseWriter, err.Error(), http.StatusInternalServerError)
					}

					_, err = res.ResponseWriter.Write(respData)
					if err != nil {
						http.Error(res.ResponseWriter, err.Error(), http.StatusInternalServerError)
					}

					return
				}
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
