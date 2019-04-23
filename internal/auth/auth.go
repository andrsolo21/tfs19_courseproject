package auth

import (
	"courseproject/internal/lots"
	"courseproject/internal/session"
	"courseproject/internal/storages"
	"courseproject/internal/user"
	"courseproject/internal/users"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

func AddUser(add users.User, data storages.DataB) (ierr int, err error) {

	ierr, err = checkUser(add, data)

	if err == nil {
		add.CreatedAt = time.Now()
		add.UpdatedAt = time.Now()

		err = data.AddUser(add)

		return ierr, err
	}
	return ierr, err
}

func checkUser(add users.User, db storages.DataB) (ierr int, err error) {

	if add.Email == "" {
		return 400, errors.New("email is empty")
	}

	if add.Password == "" {
		return 400, errors.New("password is empty")
	}

	if db.CheckEmail(add.Email) {
		return 409, errors.New("email already exist")
	}

	return 201, nil
}

func ChangeUser(id int, upd users.User, db storages.DataB) users.User {

	return db.ChangeUser(upd, id)
}

func CreateSession(log string, pas string, data storages.DataB) (string, error) {

	us, err := data.GetUserByEmPass(log, pas)

	if err != nil {
		return "tokenNotSafety", errors.New("invalid email or password")
	}

	token := genToken()
	_, err = data.GetSesByToken(token)

	for err == nil {
		token = genToken()
		_, err = data.GetSesByToken(token)
	}

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

func MassLotsToJSON(lts []lots.Lot, db storages.DataB) ([]byte, error) {
	var out []lots.LotForJSON
	for _, l := range lts {
		out = append(out, ToJSONLot(l, db))
	}
	return json.Marshal(out)
}

/*
func (a Auth) LenLots() int {
	return len(a.lots)
}*/

func ToJSONLot(l lots.Lot, db storages.DataB) lots.LotForJSON {
	user1, _ := db.GetUserByID(l.CreatorID)
	user2, _ := db.GetUserByID(l.BuyerID)

	return lots.LotForJSON{
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
		DeletedAt:   l.DeletedAt,
	}
}
