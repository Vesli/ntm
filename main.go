package main

/*
   Remember to structure via data.
   ex:
   	Is there an email to send at subscription (email package)
	Is there any permission on event creation, subscription? (permissions package)

	USE Chi Render for responses and Content-Type JSON
*/

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

// Config API
type Config struct {
	Port string
}

// Service api
type Service struct {
	Conf   *Config
	Router *chi.Mux
}

func loadEnvironment() *Config {
	c := &Config{
		Port: os.Getenv("port"),
	}

	return c
}

func routes(s *Service) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
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

func newService() *Service {
	s := &Service{
		Conf:   loadEnvironment(),
		Router: chi.NewRouter(),
	}

	s.Router.Mount("/api", routes(s))
	return s
}

func main() {
	environment := flag.String("env", "dev", "Which env to run the API")
	flag.Parse()

	_ = environment

	s := newService()
	_ = s
}
