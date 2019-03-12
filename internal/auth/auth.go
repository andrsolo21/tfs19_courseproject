package auth

import (
	"courseproject/internal/session"
	"courseproject/internal/user"
	"time"
)

type Auth struct {
	data []user.User
	ses []session.Session
}

func (auth Auth) AddUser(add user.User)(Auth, bool) {
	if auth.checkUser(add) {
		auth.data = append(auth.data, add)
		return auth,true
	}
	return auth, false
}

func (auth Auth) checkUser(add user.User) bool {
	//TODO будет в SQL
	for _, el:= range auth.data {
		if auth.checkUser(el){
			return false
		}
	}
	return true
}

func (auth Auth) LenA()int{
	return len(auth.data)
}

func (auth Auth) CreateSession(log string, pas string)(Auth,string, bool){

	us, flag := auth.authUser(log,pas)

	if flag == false{
		return auth, "tokenNotSafety", false
	}

	valid := time.Duration(5 * time.Hour)

	sesio = session.Session{
		Session_id = genToken()
		User_id = us.ID
		Created_at = time.Now()
		Valid_until = time.Now().
	}

	auth.ses = append(auth.ses, sesio)

	return auth, "tokenSafety", true
}

func (auth Auth) authUser(log string, pas string)(user.User, bool){
	//TODO будет в SQL
	for _, el:= range auth.data {
		if el.AuthUser(log, pas){
			return el, true
		}
	}
	return auth.data[0], false
}

func genToken()string{
	return "tokenSafety"
}