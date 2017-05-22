package main

import (
	"context"
	"net/http"
)

type contextKey string

const (
	mgoSession    contextKey = "MGOsession"
	configuration contextKey = "Configuration"
)

func initMiddleware(service *service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			copySession := service.Session.Copy()
			defer copySession.Close()

			ctx := context.WithValue(r.Context(), mgoSession, copySession)
			ctx = context.WithValue(ctx, configuration, service.Config)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
