package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.com/andrsolo21/courseproject/internal/auth"
	"gitlab.com/andrsolo21/courseproject/internal/lot"
	"gitlab.com/andrsolo21/courseproject/internal/lots"
	"gitlab.com/andrsolo21/courseproject/internal/storages"
	tmpl "gitlab.com/andrsolo21/courseproject/internal/templates"
	"gitlab.com/andrsolo21/courseproject/internal/user"
	"gitlab.com/andrsolo21/courseproject/internal/users"
	"gitlab.com/andrsolo21/courseproject/pkg/log"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type rout struct {
	db        storages.INTT
	logger    log.Logger
	templates map[string]*template.Template
}

func (dbr rout) signup(w http.ResponseWriter, r *http.Request) {

	var resp users.User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//fmt.Fprintf(w, "err %q\n", err, err.Error())
		http.Error(w, "", 400)
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
		http.Error(w, "", 400)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
		_, err = w.Write(mapVar)
		if err != nil {
			dbr.logger.Errorf("can't send error: %s", err.Error())
		}
		return
	}

	//resp, err := user.ConvertDate(respInp)

	/*if err != nil {
		dbr.returnError(w , "bad : %s, %v" ,err , 400 , respInp)
	}
	*/
	err = user.CheckDate(resp.Birthday)

	if err != nil {
		dbr.returnError(w, "Problem with birthday: %s%s", err, 400, "")
	}
	ierr, err := auth.AddUser(resp, dbr.db.Db())

	if err != nil {

		dbr.logger.Debugf("error in signup: %s", err.Error())

		http.Error(w, "", ierr)
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

	//_, err = w.Write([]byte("Пользователь зарегистрирован"))

}

func (dbr rout) signin(w http.ResponseWriter, r *http.Request) {

	var (
		token string
	)

	m := make(map[string]string)
	switch r.Header.Get("Content-Type") {
	case "multipart/form-data":
		m["email"] = r.PostFormValue("email")
		m["password"] = r.PostFormValue("password")
	case "application/json":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			//fmt.Fprintf(w, "err %q\n", err, err.Error())
			http.Error(w, "", http.StatusBadRequest)
			mapVar, _ := json.Marshal(map[string]string{"error": "can't readAll"})
			_, _ = w.Write(mapVar)
			return
		}
		err = json.Unmarshal(body, &m)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
			_, _ = w.Write(mapVar)
			return
		}
	default:
		dbr.returnError(w, "this body data doesn't support%s%s", errors.New(""), 404, "")
	}

	token, err := auth.CreateSession(m["email"], m["password"], dbr.db.Db())

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
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbr.logger.Infof("Authorized request token: %s", token)

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//http.Error(w, "can't read user's ID", http.StatusBadRequest)
		dbr.returnError(w, "can't read user's ID: %s, %+v", err, 400, "")
		dbr.logger.Debugf("can't read user's ID: %s", token)
		return
	}

	if dbr.db.CountUsers() < userID {
		dbr.returnError(w, "user for the given ID is not found: %s, %+v", errors.New("user for the given ID is not found:"+chi.URLParam(r, "id")), 404, "")
		//http.Error(w, "user for the given ID is not found", http.StatusNotFound)
		return
	}

	if userID == ses.UserID {
		userID = 0
	}

	userID2 := userID
	if userID == 0 {
		userID = ses.UserID
	}

	us, err := dbr.db.GetUserByID(userID)

	if err != nil {
		//http.Error(w, "user for the given ID not found", http.StatusNotFound)
		dbr.returnError(w, "user for the given ID is not found: %s, %+v", errors.New("user for the given ID is not found:"+chi.URLParam(r, "id")), 404, "")
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
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
		return
	}

	var upd users.User

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
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

	err = user.CheckDate(upd.Birthday)
	if err != nil {
		dbr.returnError(w, "Problem with birthday: %s%s", err, 400, "")
		return
	}

	if upd.FirstName == "" || upd.LastName == "" {
		//http.Error(w, "empty names", http.StatusUnauthorized)
		dbr.returnError(w, "first name or last name is empty%s%s", errors.New(""), 400, "")
		return
	}

	if upd.Email != "" || upd.Password != "" {
		dbr.returnError(w, "can't update email and/or password%s%s", errors.New(""), 400, "")
		return
	}

	us := auth.ChangeUser(sesio.UserID, upd, dbr.db.Db())

	dbr.logger.Infof("User %d was changed", sesio.UserID)
	_, _ = w.Write(user.ToJSON(true, us))
}

func (dbr rout) getUsersLots(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	ses, err := dbr.db.GetSesByToken(token)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//http.Error(w, "can't read user's ID", http.StatusBadRequest)
		dbr.returnError(w, "can't read user's ID", nil, 400, "")
		return
	}

	typ := r.URL.Query().Get("type")

	if userID == 0 {
		userID = ses.UserID
	}
	lts := dbr.db.GetUsersLots(userID, typ)

	jLots, err := auth.MassLotsToJSON(lot.Separate(lts), dbr.db.Db())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		dbr.logger.Errorf("InternalServerError: %s", err.Error())
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
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
		return
	}

	//typ := chi.URLParam(r, "type")
	//var typ string
	//switch r.Header.Get("Content-Type") {
	//case "multipart/form-data":
	//typ = r.PostFormValue("status")
	//}
	typ := r.URL.Query().Get("status")
	lts, _ := dbr.db.GetLots(typ)

	jLots, err := auth.MassLotsToJSON(lot.Separate(lts), dbr.db)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = w.Write(jLots)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (dbr rout) addLot(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	ses, err := dbr.db.GetSesByToken(token)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
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
		dbr.logger.Debugf("can't unmarshall message: %s", err.Error())
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
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
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
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
		return
	}

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//http.Error(w, "can't read user's ID", http.StatusBadRequest)
		dbr.returnError(w, "can't read user's ID%s%s", errors.New(""), 400, "")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't readALL", 400)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't readALL"})
		_, _ = w.Write(mapVar)
		return
	}

	var resp lots.LotTCU
	err = json.Unmarshal(body, &resp)
	if err != nil {
		dbr.logger.Debugf("can't unmarshall message:%+v, %s", body, err.Error())
		http.Error(w, "", 400)
		mapVar, _ := json.Marshal(map[string]string{"error": "can't unmarshal"})
		_, _ = w.Write(mapVar)
		return
	}

	el, err := lot.UpdateLot(ses.UserID, resp, lotID, dbr.db)

	if err != nil {
		//http.Error(w, err.Error(), 404)
		dbr.returnError(w, "can't update lot: %s%s", err, 404, "")
		return
	}

	mess, err := json.Marshal(auth.ToJSONLot(el, dbr.db.Db()))
	if err != nil {
		http.Error(w, "problem with marshalling lot", 500)
		return
	}
	_, _ = w.Write(mess)
}

func (dbr rout) getLot(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	_, err := dbr.db.GetSesByToken(token)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s %+v", errors.New(token), 401, "")
		return
	}

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//http.Error(w, "can't read user's ID", http.StatusBadRequest)
		dbr.returnError(w, "can't read user's ID%s%s", errors.New(""), 400, "")
		return
	}

	el, err := dbr.db.GetLotByID(lotID)

	if err != nil {
		//http.Error(w, err.Error(), 404)
		dbr.returnError(w, "can't get lot: %s%s", err, 404, err)
		return
	}

	mess, err := json.Marshal(auth.ToJSONLot(el, dbr.db.Db()))
	if err != nil {
		http.Error(w, "problem with marshalling lot", 500)
		return
	}
	_, _ = w.Write(mess)
}

func (dbr rout) deleteLot(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	ses, err := dbr.db.GetSesByToken(token)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s%s", errors.New(token), 401, "")
		return
	}

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//http.Error(w, "can't read user's ID", http.StatusBadRequest)
		dbr.returnError(w, "can't read user's ID%s%s", errors.New(""), 400, "")
		return
	}

	//el, err := dbr.db.GetLotByID(lotID)
	err = lot.DeleteLot(ses.UserID, lotID, dbr.db)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusNotFound)
		dbr.returnError(w, "Can't delete error: %s ID: %d", err, 404, lotID)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dbr rout) getLotsHTML(w http.ResponseWriter, r *http.Request) {

	lts, _ := dbr.db.GetLots("")
	tmpl.RenderTemplate(w, "index", "base", lot.Separate(lts), dbr.templates)
}

func (dbr rout) lotDescrHTML(w http.ResponseWriter, r *http.Request) {

	lotID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "can't read user's ID", http.StatusBadRequest)
		return
	}

	lts, err := dbr.db.GetLotByID(lotID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	tmpl.RenderTemplate(w, "lotDescription", "base", lts, dbr.templates)
}

/*
func (dbr rout) UpdateLots(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("can't upgrade connection: %s\n", err)
		return
	}
	//defer conn.Close()

	//go webs.SendNewLots(conn, dbr.db, dbr.templates)

	for {
		lts, _ := dbr.db.GetLots("")
		msg, err := json.Marshal(lts)

		// Write message back to browser
		if err = conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {

			return
		}
	}


		for n := 0; n < 10; n++ {
			msg := "hello  " + string(n+48)
			//fmt.Printf("sending to client: %s\n", msg)
			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			_, reply, err := conn.ReadMessage()

			if err != nil {
				fmt.Printf("can't receive: %s\n", err)
			}
			fmt.Printf("received back from client: %s\n", string(reply[:]))
		}
}
*/
func (dbr rout) returnError(w http.ResponseWriter, format string, err error, ierr int, resp interface{}) {

	dbr.logger.Debugf(format, err.Error(), resp)

	http.Error(w, "", ierr)
	mapVar, _ := json.Marshal(map[string]string{"error": errors.Errorf(format, err.Error(), resp).Error()})

	/*if err != nil{
		dbr.logger.Errorf("can't marshal error: %s", err.Error())
	}*/

	_, err = w.Write(mapVar)
	if err != nil {
		dbr.logger.Errorf("can't send message: %s", err.Error())
	}
}

func (dbr rout) getUsersLotsHTML(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	ses, err := dbr.db.GetSesByToken(token)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		dbr.logger.Debugf("Unauthorized request token: %s", token)
		dbr.returnError(w, "Unauthorized request token: %s, %+v", errors.New(token), 401, "")
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//http.Error(w, "can't read user's ID", http.StatusBadRequest)
		dbr.returnError(w, "can't read user's ID%s%s", errors.New(""), 400, "")
		return
	}

	typ := r.URL.Query().Get("type")

	if userID == 0 {
		userID = ses.UserID
	}
	lts := dbr.db.GetUsersLots(userID, typ)

	tmpl.RenderTemplate(w, "index", "base", lot.Separate(lts), dbr.templates)
}
