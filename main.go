package main

import (
	"fmt"
	"net/http"

	"github.com/pressly/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Data NTM!"))
	})
	http.ListenAndServe(fmt.Sprintf(":%d", 8080), r)
}
