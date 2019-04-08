package auth

import (
	"courseproject/internal/database"
	"courseproject/internal/lotS"
	"courseproject/internal/session"
	"courseproject/internal/user"
	"courseproject/internal/userS"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func genToken() string {
	size := 10

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)
	return rs
}

func AddUser(add userS.User, data storages.DataB) (err error) {

	err = checkUser(add, data)

	if err == nil {
		add.CreatedAt = time.Now()
		add.UpdatedAt = time.Now()

		err = data.AddUser(add)

		return err
	}
	return err
}

func checkUser(add userS.User, db storages.DataB) (err error) {

	if db.CheckEmail(add.Email) {
		return errors.New("email already exist")
	}

	return nil
}

func ChangeUser(id int, upd userS.User, db storages.DataB) userS.User {

	return db.ChangeUser(upd, id)
}

func CreateSession(log string, pas string, data storages.DataB) (string, error) {

	us, err := data.GetUserByEmPass(log, pas)

	if err != nil {
		return "tokenNotSafety", errors.New("invalid email or password")
	}

	token := genToken()

	sesio := session.CreateSession(token, us.ID)

	err = data.AddSession(sesio)

	return token, err
}

/*func (a Auth) GetMassLots(id int) []lot.Lot {
	var lts []lot.Lot

	for _, l := range a.lots {
		if l.CreatorID == id {
			lts = append(lts, l)
		}
	}

	return lts
}*/

/*
func (a Auth) GetAllLots() []lot.Lot {
	return a.lots
}
*/


func MassLotsToJSON(lots []lotS.Lot, db storages.DataB) ([]byte, error) {
	var out []lotS.LotForJSON
	for _, l := range lots {
		out = append(out, ToJsonLot(l, db))
	}
	return json.Marshal(out)
}

/*
func (a Auth) LenLots() int {
	//TODO будет с SQL вообще не обязательно
	return len(a.lots)
}*/

func ToJsonLot(l lotS.Lot, db storages.DataB) lotS.LotForJSON {
	user1, _ := db.GetUserByID(l.CreatorID)
	user2, _ := db.GetUserByID(l.BuyerID)

	return lotS.LotForJSON{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		BuyPrice:    l.BuyPrice,
		MinPrice:    l.MinPrice,
		PriceStep:   l.PriceStep,
		Status:      l.Status,
		EndAt:       l.EndAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
		CreatorID:   user.ToShort(user1),
		BuyerID:     user.ToShort(user2),
	}
}
