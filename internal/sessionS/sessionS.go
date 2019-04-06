package sessionS

import "time"

type Session struct {
	Session_id  string    `json:"id" gorm:"AUTO_INCREMENT; PRIMARY_KEY"`
	User_id     int       `json:"user_id" `
	Created_at  time.Time `json:"created_at"`
	Valid_until time.Time `json:"valid_until"`
}
