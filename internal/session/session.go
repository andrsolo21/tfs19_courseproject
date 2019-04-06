package session

import (
	"courseproject/internal/sessionS"
	"time"
)

func CreateSession(token string, id int) sessionS.Session {

	valid := time.Duration(5 * time.Hour)

	return sessionS.Session{
		Session_id:  token,
		User_id:     id,
		Created_at:  time.Now(),
		Valid_until: time.Now().Add(valid),
	}
}
