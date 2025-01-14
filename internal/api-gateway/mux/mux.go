package mux

import (
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	. "story-pulse/internal/api-gateway/middlewares"
)

const (
	usersApiPrefix = "/v1/users"
	authApiPrefix  = "/v1/auth"
)

func Register(httpMux *chi.Mux, grpcMux *runtime.ServeMux) {
	httpMux.Route(authApiPrefix, func(r chi.Router) {
		r.With(RefreshTokenToCookieMiddleware()).Route("/register", func(r chi.Router) {
			r.Mount("/", grpcMux)
		})

		r.With(RefreshTokenToCookieMiddleware()).Route("/login", func(r chi.Router) {
			r.Mount("/", grpcMux)
		})

		r.With(RequiredCookieMiddleware()).Route("/refresh", func(r chi.Router) {
			r.Mount("/", grpcMux)
		})

		r.Mount("/", grpcMux)
	})

	httpMux.Route(usersApiPrefix, func(r chi.Router) {
		r.With(AuthAndIDMiddleware()).Route("/{id}", func(r chi.Router) {
			r.Mount("/", grpcMux)
		})

		r.Mount("/", grpcMux)
	})
}
