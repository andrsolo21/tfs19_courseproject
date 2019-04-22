package main

import (
	"courseproject/internal/storages"
	"courseproject/pkg/log"
	"net/http"

	"github.com/go-chi/chi"
)

/*func gRPCService(serviceAddr string)  {
lis , err := net.Listen("tcp", serviceAddr)
if err != nil{
	//fatal
}



s := grpc.NewServer()

//session

*/

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

	data2.logger = log.New()

	go data2.db.KillBadLots(data2.logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", data2.signup)
		r.Post("/signin", data2.signin)

		r.Get("/users/{id}", data2.userGet)
		r.Put("/users/{id}", data2.userPut)
		r.Get("/users/{id}/lots", data2.getUsersLots)

		r.Get("/lots", data2.getLots)
		r.Post("/lots", data2.addLot)
		r.Put("/lots/{id}/buy", data2.buyLot)
		r.Get("/lots/{id}", data2.getLot)
		r.Put("/lots/{id}", data2.updateLot)
		r.Delete("/lots/{id}", data2.deleteLot)
	})

	//serv := server.New()
	//serv.Start()

	_ = http.ListenAndServe(":5000", r)
}
