package main

import (
	//"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/api/v1/signup", signup)
		r.Post("/api/v1/signin", signin)
		r.Get("/api/v1/signup", signup2)
	})
	http.ListenAndServe(":5000", r)
}
