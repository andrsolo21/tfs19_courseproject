package auth

import (
	"courseproject/internal/session"
	"courseproject/internal/user"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

type Auth struct {
	data []user.User
	ses  []session.Session
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

func (auth Auth) ProfileUpdate(token string, upd user.User) (Auth, bool) {

	sesio, flag := auth.GetSession(token)
	if flag == false {
		return auth, false
	}

	auth.changeUser(sesio.User_id, upd)
	return auth, true
}

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

func (auth Auth) changeUser(id int, upd user.User) {
	//TODO будет с SQL
	for i := range auth.data {
		if auth.data[i].ID == id {
			auth.data[i].Updated_at = time.Now()
			auth.data[i].Birthday = upd.Birthday
			auth.data[i].Last_name = upd.Last_name
			auth.data[i].First_name = upd.First_name
			return
		}
	}
}
