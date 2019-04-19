package main

import (
	"courseproject/internal/storages"
	"courseproject/pkg/log"
	"github.com/go-chi/chi"

	"net/http"
)

func main() {

	var err error

	r := chi.NewRouter()

	logger := log.New()
	//r.Use(NewStructuredLogger(logger))

	var data2 rout

	data2.db, err = storages.NewDataB()
	if err != nil {
		logger.Fatal(err.Error())
	}

	//logger.Error("test")
	//logger.Debug("test debud")

	//logger.Info("test info")

	//logger.Warn("test warn")


	defer data2.db.Db().DB.Close()

	data2.db = data2.db.CreateTables()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", data2.signup)
		r.Post("/signin", data2.signin)

		r.Get("/users/{id}", data2.userGet)
		r.Put("/users/{id}", data2.userPut)
		r.Get("/users/{id}/lots", data2.getUsersLots)

		r.Get("/lots", data2.getLots)
		r.Post("/lots", data2.addLot)
	})

	_ = http.ListenAndServe(":5000", r)
}
