package main

import (
	"net/http"

	"github.com/pressly/chi"
)

func registerRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the NTM API!"))
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/subscribe", registerAndLogginUser)
		r.Route("/profile", func(r chi.Router) {
			r.Get("/:id", getUserFromDB)
		})
	})

	r.Route("/event", func(r chi.Router) {
		r.Post("/create", createEvent)
	})
	return r
}
