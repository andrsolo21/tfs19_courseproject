package users

import "time"

type User struct {
	ID        int       `json:"id" gorm:"AUTO_INCREMENT; PRIMARY_KEY"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Birthday  string    `json:"birthday,omitempty"`
	Email     string    `json:"email" gorm:"not null; unique"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type ShortUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	ID        int    `json:"id,omitempty"`
}


func U(a string) string{
	return a
}