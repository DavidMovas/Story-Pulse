package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/gommon/log"
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
				var resp map[string]any
				if err := json.Unmarshal(res.body.Bytes(), &resp); err == nil {
					if refreshToken, ok := resp["refreshToken"].(string); ok {
						http.SetCookie(w, &http.Cookie{
							Name:     "refresh_token",
							Value:    refreshToken,
							HttpOnly: true,
							Secure:   true,
							Path:     "/",
						})

						delete(resp, "refreshToken")
					}
				}

				respData, err := json.Marshal(resp)
				if err != nil {
					http.Error(res.ResponseWriter, err.Error(), http.StatusInternalServerError)
				}

				bytesWritten, err := w.Write(respData)
				if err != nil {
					http.Error(res.ResponseWriter, err.Error(), http.StatusInternalServerError)
				}

				log.Infof("BYTES: %d", bytesWritten)
			}

			_, _ = w.Write(res.body.Bytes())
		})
	}
}

/*func RefreshTokenToCookieMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}
			next.ServeHTTP(res, r)

			if res.statusCode == http.StatusOK || res.statusCode == 0 {
				var resp map[string]any
				if err := json.Unmarshal(res.body.Bytes(), &resp); err == nil {
					if refreshToken, ok := resp["refreshToken"].(string); ok {
						http.SetCookie(w, &http.Cookie{
							Name:     "refresh_token",
							Value:    refreshToken,
							HttpOnly: true,
							Secure:   true,
							Path:     "/",
						})

						delete(resp, "refreshToken")
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)

					if err := json.NewEncoder(w).Encode(resp); err != nil {
						http.Error(w, "Failed to encode response", http.StatusInternalServerError)
						return
					}
					return
				}
			}

			w.WriteHeader(res.statusCode)
			_, _ = w.Write(res.body.Bytes())
		})
	}
}*/

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
