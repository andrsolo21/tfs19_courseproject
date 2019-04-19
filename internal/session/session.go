package session

import (
	"courseproject/internal/sessions"
	"time"
)

func CreateSession(token string, id int) sessions.Session {

	valid := 5 * time.Hour

	return sessions.Session{
		SessionID:  token,
		UserID:     id,
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().Add(valid),
	}
}
