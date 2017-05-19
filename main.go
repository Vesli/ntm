package main

/*
   Decouper via data.
   ex:
       Y a un t-il un email a envoyer a inscription (package email)
       Y a t-il des permissions pour un event, pour inscription? (package permissions)
       Des medias?
*/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pressly/chi"
)

func main() {
	r := chi.NewRouter()

	r.Mount("/", registerRoutes())
	fmt.Println("Running on port:", 8080)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), r))
}
