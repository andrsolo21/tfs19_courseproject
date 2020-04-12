package sessions

import "time"

type Session struct {
	SessionID  string    `json:"id" gorm:"AUTO_INCREMENT; PRIMARY_KEY"`
	UserID     int       `json:"user_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	ValidUntil time.Time `json:"valid_until"`
}

func S(a string) string{
	return a
}
