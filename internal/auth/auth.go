package auth

import (
	"courseproject/internal/lot"
	"courseproject/internal/session"
	"courseproject/internal/user"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Auth struct {
	data []user.User
	ses  []session.Session
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

func (a Auth) AddUser(add user.User) (Auth, string, bool) {
	str, flag := a.checkUser(add)
	if flag {
		add.ID = a.LenUsers() + 1
		a.data = append(a.data, add)
		return a, "ok", true
	}
	return a, str, false
}

func (auth Auth) checkUser(add user.User) (string, bool) {
	//TODO будет с SQL
	var (
		str string
		fl  bool
	)
	for _, el := range auth.data {
		str, fl = add.CheckUser(el)
		if fl {
			return str, false
		}
	}
	return str, true
}

func (auth Auth) LenUsers() int {
	//TODO будет с SQL вообще не обязательно
	return len(auth.data)
}

func (auth Auth) ChangeUser(id int, upd user.User) user.User {
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
	return user.User{}
}

func (auth Auth) authUser(log string, pas string) (user.User, bool) {
	//TODO будет с SQL
	for _, el := range auth.data {
		if el.AuthUser(log, pas) {
			return el, true
		}
	}
	return user.User{}, false
}

func (auth Auth) GetUserById(ID int) (user.User, bool) {
	//TODO будет с SQL
	for _, el := range auth.data {
		if el.ID == ID {
			return el, true
		}
	}
	return user.User{}, false
}

//Sessions

func (auth Auth) GetSession(token string) (session.Session, bool) {
	//TODO будет с SQL
	for _, el := range auth.ses {
		if token == el.Session_id {
			if time.Now().After(el.Valid_until) {
				return session.Session{}, false
			}
			return el, true
		}
	}
	return session.Session{}, false
}

func (auth Auth) CreateSession(log string, pas string) (Auth, string, string) {

	us, flag := auth.authUser(log, pas)

	if flag == false {
		return auth, "tokenNotSafety", "invalid email or password"
	}

	valid := time.Duration(5 * time.Hour)

	token := genToken()

	sesio := session.Session{
		Session_id:  token,
		User_id:     us.ID,
		Created_at:  time.Now(),
		Valid_until: time.Now().Add(valid),
	}

	auth.ses = append(auth.ses, sesio)

	return auth, token, ""
}

//Lots

func (a Auth) GetMassLots(id int) []lot.Lot {
	var lts []lot.Lot

	for _, l := range a.lots {
		if l.CreatorID == id{
			lts = append(lts, l)
		}
	}

	return lts
}

func (a Auth) GetAllLots() []lot.Lot {
	return a.lots
}

func (a Auth) MassLotsToJSON(lots []lot.Lot) ([]byte, error){
	var out []lot.LotForJSON
	for _, l := range(lots){
		out = append(out, ToJsonLot(a,l))
	}
	return json.Marshal(out)
}

func (a Auth) LenLots() int {
	//TODO будет с SQL вообще не обязательно
	return len(a.lots)
}

func ToJsonLot(data Auth, l lot.Lot) lot.LotForJSON {
	user1, _ := data.GetUserById(l.CreatorID)
	user2, _ := data.GetUserById(l.BuyerID)

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
		CreatorID:   user1.ToShort(),
		BuyerID:     user2.ToShort(),
	}
}