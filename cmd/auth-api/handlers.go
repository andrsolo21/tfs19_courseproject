package main

import (
	"courseproject/internal/user"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func strartListening() {

	r := chi.NewRouter()

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	//r.Use(NewStructuredLogger(logger))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", signup)
		r.Post("/signin", signin)
		r.Put("/users/0", users0)
	})

	http.ListenAndServe(":5000", r)
}

func signup(w http.ResponseWriter, r *http.Request) {

	var resp user.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//fmt.Fprintf(w, "err %q\n", err, err.Error())
		http.Error(w, "can't readAll", http.StatusBadRequest)
		return
	} else {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			http.Error(w, "can't unmarshal", http.StatusBadRequest)
			return
		}

		resp.Created_at = time.Now()
		resp.Updated_at = time.Now()

		var (
			flag   bool
			errStr string
		)
		data, errStr, flag = data.AddUser(resp)

		if !flag {
			http.Error(w, errStr, http.StatusConflict)
			return
		}
	}

}

func signin(w http.ResponseWriter, r *http.Request) {

	var token , errStr string

	data, token, errStr = data.CreateSession(r.PostFormValue("email"), r.PostFormValue("password"))

	if errStr == "" {
		mapVar2, _ := json.Marshal(map[string]string{"access_token": token, "token_type": "bearer"})
		w.Write([]byte(mapVar2))

	} else {

		http.Error(w, errStr, http.StatusUnauthorized)
		//w.Write([]byte())
		return
	}
}

func users0(w http.ResponseWriter, r *http.Request) {

	var resp user.User
	var err error

	resp.First_name = r.PostFormValue("first_name")
	resp.Last_name = r.PostFormValue("last_name")
	resp.Birthday, err = time.Parse("2006-01-02T15:04:05-07:00", r.PostFormValue("Birthday"))
	if err != nil {
		fmt.Println(w, "can't parse: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	token := r.Header.Get("Authorization")

	var flag bool
	data, flag = data.ProfileUpdate(token, resp)

	if flag {
		w.Header().Add("X-MY-added", "true")
	} else {
		w.Header().Add("X-MY-added", "false")
	}

	//fmt.Fprintln(w, "otv", flag)

}
