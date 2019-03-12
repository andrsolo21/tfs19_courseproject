package main

import (
	"courseproject/internal/auth"
	//"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

var data auth.Auth

func main() {

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", signup)
		r.Post("/signin", signin)
		r.Put("/users/0", users0)
	})
	http.ListenAndServe(":5000", r)
}
