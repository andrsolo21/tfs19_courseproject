package auth

import (
	"courseproject/internal/database"
	"courseproject/internal/lot"
	"courseproject/internal/session"
	"courseproject/internal/sessionS"
	"courseproject/internal/user"
	"courseproject/internal/userS"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Auth struct {
	data []userS.User
	ses  []sessionS.Session
	lots []lot.Lot
}

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

/*func (auth Auth) LenUsers() int {
	//TODO будет с SQL вообще не обязательно
	return len(auth.data)
}*/

func (auth Auth) ChangeUser(id int, upd userS.User) userS.User {
	//TODO будет с SQL
	for i := range auth.data {
		if auth.data[i].ID == id {
			auth.data[i].UpdatedAt = time.Now()
			auth.data[i].Birthday = upd.Birthday
			auth.data[i].LastName = upd.LastName
			auth.data[i].FirstName = upd.FirstName
			return auth.data[i]
		}
	}
	return userS.User{}
}

func authUser(log string, pas string, db storages.DataB) (userS.User, bool) {

	var el userS.User

	db.DB.Where("email = ? AND password = ?", log, pas).First(&el)

	if el.ID != 0 {
		return el, true
	}
	return userS.User{}, false
}

//Sessions

func CreateSession(log string, pas string, data storages.DataB) (string, error) {

	us, flag := authUser(log, pas, data)

	if flag == false {
		return "tokenNotSafety", errors.New("invalid email or password")
	}

	token := genToken()

	sesio := session.CreateSession(token, us.ID)

	err := data.AddSession(sesio)

	return token, err
}

//Lots

func (a Auth) GetMassLots(id int) []lot.Lot {
	var lts []lot.Lot

	for _, l := range a.lots {
		if l.CreatorID == id {
			lts = append(lts, l)
		}
	}

	return lts
}

func (a Auth) GetAllLots() []lot.Lot {
	return a.lots
}

func MassLotsToJSON(lots []lot.Lot, db storages.DataB) ([]byte, error) {
	var out []lot.LotForJSON
	for _, l := range lots {
		out = append(out, ToJsonLot(l, db))
	}
	return json.Marshal(out)
}

/*func (a Auth) LenLots() int {
	//TODO будет с SQL вообще не обязательно
	return len(a.lots)
}*/

func ToJsonLot(l lot.Lot, db storages.DataB) lot.LotForJSON {
	user1 := db.GetUserByID(l.CreatorID)
	user2 := db.GetUserByID(l.BuyerID)

	return lot.LotForJSON{
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
