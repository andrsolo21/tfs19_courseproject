package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/andrsolo21/courseproject/internal/storages"
	tmpl "gitlab.com/andrsolo21/courseproject/internal/templates"
	"gitlab.com/andrsolo21/courseproject/pkg/log"

	"github.com/go-chi/chi"
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

	data2.logger = log.New()

	data2.templates = tmpl.NewTempl()

	go data2.db.KillBadLots(data2.logger)

	r.Route("/v1/auction", func(r chi.Router) {
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

		r.Get("/lots/html", data2.getLotsHTML)
		r.Get("/lot/{id}/html", data2.lotDescrHTML)
		r.Get("/users/{id}/lots/html", data2.getUsersLotsHTML)

	})

	//http.HandleFunc("/ws", data2.UpdateLots)

	//serv := server.New()
	//serv.Start()

	//_ = http.ListenAndServe(":5000", r)
	//_ = http.ListenAndServe(":5001", nil)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "swagger")
	FileServer(r, "/swagger", http.Dir(filesDir))
	if err := http.ListenAndServe(":5000", r); err != nil {
		logger.Fatalf("server error:%s", err)
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusTemporaryRedirect).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

}
