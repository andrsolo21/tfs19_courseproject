package session

import (
	"time"

	"gitlab.com/andrsolo21/courseproject/internal/sessions"
)

func CreateSession(token string, id int) sessions.Session {

	valid := 5 * time.Hour

	return sessions.Session{
		SessionID:  sessions.S(token),
		UserID:     id,
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().Add(valid),
	}
}
