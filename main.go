package main

/*
   Decouper via data.
   ex:
       Y a un t-il un email a envoyer a inscription (package email)
       Y a t-il des permissions pour un event, pour inscription? (package permissions)
       Des medias?
*/

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pressly/chi"
)

func main() {
	pathToConfig := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	service, err := newService(*pathToConfig)
	if err != nil {
		panic(err)
	}
	defer service.Close()

	r := chi.NewRouter()
	initRet := initMiddleware(service)

	r.Use(initRet)
	r.Mount("/ntm-api", registerRoutes())

	fmt.Println("Running on port:", service.Config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", service.Config.Port), r))
}
