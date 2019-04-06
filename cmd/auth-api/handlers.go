package main

import (
	"courseproject/internal/auth"
	"courseproject/internal/database"
	"courseproject/internal/lot"
	"courseproject/internal/user"
	"courseproject/internal/userS"
	"encoding/json"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type rout struct {
	db storages.DataB
}

func (dbr rout) signup(w http.ResponseWriter, r *http.Request) {

	var resp userS.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//fmt.Fprintf(w, "err %q\n", err, err.Error())
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't readAll"})
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

	err = auth.AddUser(resp, dbr.db)

	if err != nil {
		http.Error(w, "", http.StatusConflict)
		mapVar, _ := json.Marshal(map[string]string{"error": err.Error()})
		_, _ = w.Write([]byte(mapVar))
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (dbr rout) signin(w http.ResponseWriter, r *http.Request) {

	var token string

	token, err := auth.CreateSession(r.PostFormValue("email"), r.PostFormValue("password"), dbr.db)

	if err == nil {
		mapVar, _ := json.Marshal(map[string]string{"access_token": token, "token_type": "bearer"})
		_, _ = w.Write([]byte(mapVar))
		return
	}

	http.Error(w, "", http.StatusUnauthorized)
	mapVar, _ := json.Marshal(map[string]string{"error": err.Error()})
	_, _ = w.Write([]byte(mapVar))

}

func (db rout) userPut(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	sesio, flag := data.GetSession(token)
	if !flag {
		http.Error(w, "problem with authorization", http.StatusUnauthorized)
		return
	}

	var upd userS.User
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

	_, _ = w.Write(user.ToJson(true, us))
}

func (db rout) userGet(w http.ResponseWriter, r *http.Request) {

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
		_, _ = w.Write(user.ToJson(true, us))
	} else {
		_, _ = w.Write(user.ToJson(false, us))
	}

}

func (db rout) getLots(w http.ResponseWriter, r *http.Request) {

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

func (db rout) addLot(w http.ResponseWriter, r *http.Request) {

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
