package main

import (
	"courseproject/internal/auth"
	"courseproject/internal/lot"
	"courseproject/internal/lots"
	"courseproject/internal/storages"
	tmpl "courseproject/internal/templates"
	"courseproject/internal/user"
	"courseproject/internal/users"
	"courseproject/pkg/log"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type rout struct {
	db     storages.INTT
	logger log.Logger
	templates map[string]*template.Template
}

func (dbr rout) signup(w http.ResponseWriter, r *http.Request) {

	var resp users.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//fmt.Fprintf(w, "err %q\n", err, err.Error())
		http.Error(w, "", http.StatusBadRequest)
		dbr.logger.Errorf("can't read message: %v", body)

		mapVar, _ := json.Marshal(map[string]string{"error": "can't readAll"})
		_, err = w.Write(mapVar)
		if err != nil {
			dbr.logger.Errorf("can't send error: %s", err.Error())
		}
		return
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		dbr.logger.Errorf("can't unmarshal message: %v", body)
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
		_, err = w.Write(mapVar)
		if err != nil {
			dbr.logger.Errorf("can't send error: %s", err.Error())
		}
		return
	}

	err = auth.AddUser(resp, dbr.db.Db())

	if err != nil {

		dbr.logger.Debugf("error in signup: %s, %v", err.Error(), resp)

		http.Error(w, "", http.StatusConflict)
		mapVar, _ := json.Marshal(map[string]string{"error": err.Error()})

		/*if err != nil{
			dbr.logger.Errorf("can't marshal error: %s", err.Error())
		}*/

		_, err = w.Write(mapVar)
		if err != nil {
			dbr.logger.Errorf("can't send message: %s", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (dbr rout) signin(w http.ResponseWriter, r *http.Request) {

	var token string

	token, err := auth.CreateSession(r.PostFormValue("email"), r.PostFormValue("password"), dbr.db.Db())

	if err == nil {

		mapVar, _ := json.Marshal(map[string]string{"access_token": token, "token_type": "bearer"})
		_, err = w.Write(mapVar)

		dbr.logger.Infof("User signin: %s", token)

		if err != nil {
			dbr.logger.Errorf("can't send message: %s", err.Error())
		}
		return
	}

	http.Error(w, "", http.StatusUnauthorized)
	mapVar, err := json.Marshal(map[string]string{"error": err.Error()})

	if err != nil {
		dbr.logger.Errorf("can't marshal error: %s", err.Error())
	}

	_, err = w.Write(mapVar)
	if err != nil {
		dbr.logger.Errorf("can't send error: %s", err.Error())
	}

}

func (dbr rout) userGet(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	ses, err := dbr.db.GetSesByToken(token)
	if err != nil {
		dbr.logger.Debugf("Unauthorized request token: %s", token)

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbr.logger.Infof("Authorized request token: %s", token)

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't read user's ID", http.StatusBadRequest)
		dbr.logger.Debugf("can't read user's ID: %s", token)
		return
	}

	if dbr.db.CountUsers() < userID {
		http.Error(w, "user for the given ID is not found", http.StatusNotFound)
		return
	}

	userID2 := userID
	if userID == 0 {
		userID = ses.UserID
	}

	us, err := dbr.db.GetUserByID(userID)

	if err != nil {
		http.Error(w, "user for the given ID not found", http.StatusNotFound)
		return
	}

	if userID2 == 0 {
		_, err = w.Write(user.ToJSON(true, us))

		if err != nil {
			dbr.logger.Errorf("can't send message: %s", err.Error())
		}
	} else {
		_, err = w.Write(user.ToJSON(false, us))

		if err != nil {
			dbr.logger.Errorf("can't send message: %s", err.Error())
		}
	}
}

func (dbr rout) userPut(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	sesio, err := dbr.db.GetSesByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var upd users.User

	switch r.Header.Get("Content-Type") {
	case "multipart/form-data":
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

	case "application/json":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			//fmt.Fprintf(w, "err %q\n", err, err.Error())
			http.Error(w, "", http.StatusBadRequest)
			mapVar, _ := json.Marshal(map[string]string{"error": "can't readAll"})
			_, _ = w.Write(mapVar)
			return
		}
		err = json.Unmarshal(body, &upd)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
			_, _ = w.Write(mapVar)
			return
		}

	}

	us := auth.ChangeUser(sesio.UserID, upd, dbr.db.Db())

	_, _ = w.Write(user.ToJSON(true, us))
}

func (dbr rout) getUsersLots(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	ses, err := dbr.db.GetSesByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't read user's ID", http.StatusBadRequest)
		return
	}

	typ := chi.URLParam(r, "type")

	if userID == 0 {
		userID = ses.UserID
	}
	lts := dbr.db.GetUsersLots(userID, typ)

	jLots, err := auth.MassLotsToJSON(lts, dbr.db.Db())

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = w.Write(jLots)

	if err != nil {
		dbr.logger.Errorf("Can't send message: %s", err.Error())
	}

}

func (dbr rout) getLots(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	_, err := dbr.db.GetSesByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//typ := chi.URLParam(r, "type")
	//var typ string
	//switch r.Header.Get("Content-Type") {
	//case "multipart/form-data":
	//typ = r.PostFormValue("status")
	//}
	typ := r.URL.Query().Get("status")

	jLots, err := auth.MassLotsToJSON(dbr.db.GetLots(typ))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, _ = w.Write(jLots)
}

func (dbr rout) addLot(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	ses, err := dbr.db.GetSesByToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var resp lots.LotTCU

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't readALL"})
		_, _ = w.Write(mapVar)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		dbr.logger.Debugf("can't unmarshall message:%+v, %s", body, err.Error())
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
		_, _ = w.Write(mapVar)
		return
	}
	lts, err := lot.Generate(resp)

	if err != nil {
		dbr.logger.Debugf("can't add lot: %s, error: %s", resp.Title, err.Error())
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't add lot " + err.Error()})
		_, _ = w.Write(mapVar)
		return
	}

	lts.CreatorID = ses.UserID

	lts, err = dbr.db.AddLot(lts)

	if err != nil {
		http.Error(w, "", http.StatusConflict)
		dbr.logger.Debugf("can't add lot:%d, %s", lts.ID, err.Error())

		mapVar, _ := json.Marshal(map[string]string{"error": err.Error()})
		_, _ = w.Write(mapVar)
		return
	}

	dbr.logger.Infof("lot was created, lot: %d", lts.ID)

	l, err := json.Marshal(auth.ToJSONLot(lts, dbr.db.Db()))
	if err != nil {
		http.Error(w, "problem with marshalling lots", http.StatusUnauthorized)
		return
	}
	_, _ = w.Write(l)

}

func (dbr rout) buyLot(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	ses, err := dbr.db.GetSesByToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't read user's ID", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't readAll"})
		_, _ = w.Write(mapVar)
		return
	}
	var upd lots.Price
	err = json.Unmarshal(body, &upd)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
		_, _ = w.Write(mapVar)
		return
	}

	el, err := lot.BuyLot(ses.UserID, lotID, upd.Price, dbr.db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	mess, err := json.Marshal(auth.ToJSONLot(el, dbr.db.Db()))
	if err != nil {
		http.Error(w, "problem with marshalling lots", http.StatusBadRequest)
		return
	}
	_, _ = w.Write(mess)
}

func (dbr rout) updateLot(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	ses, err := dbr.db.GetSesByToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't read user's ID", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't readALL"})
		_, _ = w.Write(mapVar)
		return
	}

	var resp lots.LotTCU
	err = json.Unmarshal(body, &resp)
	if err != nil {
		dbr.logger.Debugf("can't unmarshall message:%+v, %s", body, err.Error())
		http.Error(w, "", http.StatusBadRequest)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
		_, _ = w.Write(mapVar)
		return
	}

	el, err := lot.UpdateLot(ses.UserID, resp, lotID, dbr.db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	mess, err := json.Marshal(auth.ToJSONLot(el, dbr.db.Db()))
	if err != nil {
		http.Error(w, "problem with marshalling lot", http.StatusUnauthorized)
		return
	}
	_, _ = w.Write(mess)
}

func (dbr rout) getLot(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	_, err := dbr.db.GetSesByToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't read user's ID", http.StatusBadRequest)
		return
	}

	el, err := dbr.db.GetLotByID(lotID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	mess, err := json.Marshal(auth.ToJSONLot(el, dbr.db.Db()))
	if err != nil {
		http.Error(w, "problem with marshalling lot", http.StatusBadRequest)
		return
	}
	_, _ = w.Write(mess)
}

func (dbr rout) deleteLot(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	ses, err := dbr.db.GetSesByToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't read user's ID", http.StatusBadRequest)
		return
	}

	//el, err := dbr.db.GetLotByID(lotID)
	err = lot.DeleteLot(ses.UserID, lotID, dbr.db)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dbr rout) getLotsHTML(w http.ResponseWriter, r *http.Request){
	lts, _ := dbr.db.GetLots("finished")
	tmpl.RenderTemplate(w, "index", "base", lts , dbr.templates)
}
