package main

import (
	"courseproject/internal/auth"
	storages "courseproject/internal/database"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var data auth.Auth

func main() {


	var err error

	r := chi.NewRouter()

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	//r.Use(NewStructuredLogger(logger))

	var data2 rout

	data2.db, err = storages.NewDataB()
	if err != nil {
		log.Fatal(err)
	}

	defer data2.db.DB.Close()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", data2.signup)
		r.Post("/signin", data2.signin)
		r.Put("/users/{id}", data2.userPut)
		r.Get("/users/{id}", data2.userGet)
		r.Get("/lots", data2.getLots)
		r.Post("/lots", data2.addLot)
	})

	_ = http.ListenAndServe(":5000", r)
}
