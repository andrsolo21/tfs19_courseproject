package session

import "time"

type Session struct {
	Session_id  string    `json:"id"`
	User_id     int       `json:"user_id"`
	Created_at  time.Time `json:"created_at"`
	Valid_until time.Time `json:"valid_until"`
}

func CreateSession(token string, id int )Session{

	valid := time.Duration(5 * time.Hour)

	return Session{
		Session_id:  token,
		User_id:     id,
		Created_at:  time.Now(),
		Valid_until: time.Now().Add(valid),
	}
}
