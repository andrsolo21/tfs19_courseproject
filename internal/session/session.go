package session

import "time"

type Session struct {
	Session_id  string    `json:"id"`
	User_id     int       `json:"user_id"`
	Created_at  time.Time `json:"created_at"`
	Valid_until time.Time `json:"valid_until"`
}
