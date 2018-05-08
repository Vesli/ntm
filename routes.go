package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func routes(s *Service) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/", s.RegisterOrLoggin)
		r.Route("/profile", func(r chi.Router) {
			r.Get("/:id", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("/profile/:id"))
			})
		})
	})

	r.Route("/event", func(r chi.Router) {
		r.Post("/create", s.AddEvent)
	})
	return r
}
