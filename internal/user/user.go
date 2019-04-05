package user

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time //`json:"created_at"`
	UpdatedAt time.Time //`json:"updated_at"`
}

type ShortUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        int    `json:"id"`
}

var em string = "email already exists"

func AddUser(us User, bd *gorm.DB){

}

func (us1 User) CheckUser(us2 User) (string, bool) {
	if us1.Email == us2.Email {
		return em, true
	}
	return "", false
}

func (us1 User) AuthUser(log string, pas string) bool {
	if us1.Email == log && us1.Password == pas {
		return true
	}
	return false
}

func (us User) ToJson(seq bool) ([]byte) {
	var mapVar []byte
	if seq {
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.FirstName,
			"last_name":  us.LastName,
			"birthday":   us.Birthday.Format(time.ANSIC),
			"email":      us.Email,
			"created_at": us.CreatedAt.Format(time.ANSIC),})
	} else {
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.FirstName,
			"last_name":  us.LastName,
			"birthday":   us.Birthday.Format(time.ANSIC),
			"email":      "*",
			"created_at": us.CreatedAt.Format(time.ANSIC),})
	}
	return mapVar
}

func (us User) ToShort() ShortUser {
	return ShortUser{
		FirstName: us.FirstName,
		LastName: us.LastName,
		ID: us.ID,
	}
}
