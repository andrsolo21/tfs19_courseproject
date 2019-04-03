package main

import (
	"courseproject/internal/auth"
	"courseproject/internal/lot"
	"courseproject/internal/user"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
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
		r.Put("/users/{id}", userPut)
		r.Get("/users/{id}", userGet)
		r.Get("/lots", getLots)
		r.Post("/lots", addLot)
	})

	_ = http.ListenAndServe(":5000", r)
}

func signup(w http.ResponseWriter, r *http.Request) {

	var resp user.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//fmt.Fprintf(w, "err %q\n", err, err.Error())
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't readAll"})
		_, _ = w.Write([]byte(mapVar))
		return
	} else {
		err = json.Unmarshal(body, &resp)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
			_, _ = w.Write([]byte(mapVar))
			return
		}

		resp.CreatedAt = time.Now()
		resp.UpdatedAt = time.Now()

		var (
			flag   bool
			errStr string
		)
		data, errStr, flag = data.AddUser(resp)

		if !flag {
			http.Error(w, "", http.StatusConflict)
			mapVar, _ := json.Marshal(map[string]string{"error": errStr})
			_, _ = w.Write([]byte(mapVar))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

}

func signin(w http.ResponseWriter, r *http.Request) {

	var token, errStr string

	data, token, errStr = data.CreateSession(r.PostFormValue("email"), r.PostFormValue("password"))

	if errStr == "" {
		mapVar, _ := json.Marshal(map[string]string{"access_token": token, "token_type": "bearer"})
		_, _ = w.Write([]byte(mapVar))

	} else {

		http.Error(w, "", http.StatusUnauthorized)
		mapVar, _ := json.Marshal(map[string]string{"error": errStr})
		_, _ = w.Write([]byte(mapVar))
	}
}

func userPut(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	sesio, flag := data.GetSession(token)
	if !flag {
		http.Error(w, "problem with authorization", http.StatusUnauthorized)
		return
	}

	var upd user.User
	var err error

	upd.FirstName = r.PostFormValue("first_name")
	upd.LastName = r.PostFormValue("last_name")
	upd.Birthday, err = time.Parse("2006-01-02T15:04:05-07:00", r.PostFormValue("Birthday"))
	if err != nil {
		http.Error(w, "can't parse time", http.StatusUnauthorized)
		return
	}

	if upd.FirstName == "" || upd.LastName == "" {
		http.Error(w, "empty names", http.StatusUnauthorized)
		return
	}

	us := data.ChangeUser(sesio.User_id, upd)

	_, _ = w.Write(us.ToJson(true))
}

func userGet(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	ses, flag := data.GetSession(token)

	if !flag {
		http.Error(w, "problem with authorization", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't find user by ID", http.StatusBadRequest)
		return
	}

	if data.LenUsers() < userID {
		http.Error(w, "user for the given ID not found", http.StatusNotFound)
		return
	}

	userID2 := userID
	if userID == 0 {
		userID = ses.User_id
	}

	us, flag := data.GetUserById(userID)

	if !flag {
		http.Error(w, "user for the given ID not found", http.StatusNotFound)
		return
	}

	if userID2 == 0 {
		_, _ = w.Write(us.ToJson(true))
	} else {
		_, _ = w.Write(us.ToJson(false))
	}

}

func getLots(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	_, flag := data.GetSession(token)

	if !flag {
		http.Error(w, "problem with authorization", http.StatusUnauthorized)
		return
	}

	jLots, err := data.MassLotsToJSON(data.GetAllLots())

	if err != nil {
		http.Error(w, "problem with marshalling lots", http.StatusUnauthorized)
		return
	}

	_, _ = w.Write(jLots)
}

func addLot(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	ses, flag := data.GetSession(token)

	if !flag {
		http.Error(w, "problem with authorization", http.StatusUnauthorized)
		return
	}

	var resp lot.LotTCU

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't readALL"})
		_, _ = w.Write([]byte(mapVar))
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
		_, _ = w.Write([]byte(mapVar))
		return
	}
	lts := resp.Generate()

	lts.CreatorID = ses.User_id
	lts.ID = data.LenLots() + 1

	l, err := json.Marshal(auth.ToJsonLot(data, lts))
	if err != nil {
		http.Error(w, "problem with marshalling lots", http.StatusUnauthorized)
		return
	}
	_, _ = w.Write(l)

}
