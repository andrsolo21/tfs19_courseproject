package auth

import (
	"courseproject/internal/session"
	"courseproject/internal/user"
	"time"
)

type Auth struct {
	data []user.User
	ses  []session.Session
}

/*func Init(){

	var auth Auth

	auth


}*/

func genToken() string {
	return "tokenSafety"
}

func (auth Auth) AddUser(add user.User) (Auth, bool) {
	if auth.checkUser(add) {
		add.ID = auth.LenA()
		auth.data = append(auth.data, add)
		return auth, true
	}
	return auth, false
}

func (auth Auth) checkUser(add user.User) bool {
	//TODO будет с SQL
	for _, el := range auth.data {
		if add.CheckUser(el) {
			return false
		}
	}
	return true
}

func (auth Auth) LenA() int {
	//TODO будет с SQL вообще не обязательно
	return len(auth.data)
}

func (auth Auth) CreateSession(log string, pas string) (Auth, string, bool) {

	us, flag := auth.authUser(log, pas)

	if flag == false {
		return auth, "tokenNotSafety", false
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

	return auth, token, true
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

/*func (auth Auth) getUserById(ID int) (user.User, bool) {
	//TODO будет с SQL
	for _, el := range auth.data {
		if el.ID == ID {
			return el, true
		}
	}
	return user.User{}, false
}*/

func (auth Auth) ProfileUpdate(token string, upd user.User)(Auth, bool) {

	sesio, flag := auth.getSession(token)
	if flag == false{
		return auth, false
	}

	auth.changeUser(sesio.User_id, upd)
	return auth, true
}

func (auth Auth) getSession(token string) (session.Session, bool) {
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

func (auth Auth) changeUser(id int, upd user.User){
	//TODO будет с SQL
	for i := range auth.data{
		if auth.data[i].ID == id{
			auth.data[i].Updated_at = time.Now()
			auth.data[i].Birthday = upd.Birthday
			auth.data[i].Last_name = upd.Last_name
			auth.data[i].First_name = upd.First_name
			return
		}
	}
}
