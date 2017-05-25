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
		log.Fatal(err)
	}
	defer service.Close()

	r := chi.NewRouter()

	r.Use(initMiddleware(service))
	r.Mount("/ntm-api", registerRoutes())

	fmt.Println("Running on port:", service.Config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", service.Config.Port), r))
}
