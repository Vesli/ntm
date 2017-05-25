package main

/*
	Basic middleware file.
*/

import (
	"context"
	"net/http"
)

type contextKey string

const (
	psqlDB        contextKey = "psqlDB"
	configuration contextKey = "Configuration"
)

func initMiddleware(service *service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), psqlDB, service.DB)
			ctx = context.WithValue(ctx, configuration, service.Config)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
