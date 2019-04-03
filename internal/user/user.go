package user

import (
	"encoding/json"
	"strconv"
	"time"
)

type User struct {
	ID         int       `json:"id"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Birthday   time.Time `json:"birthday"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time //`json:"created_at"`
	Updated_at time.Time //`json:"updated_at"`
}

type ShortUser struct {
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Birthday   time.Time `json:"birthday"`
}

var em string = "email already exists"

func (us1 User) CheckUser(us2 User)(string , bool){
	if us1.Email == us2.Email {
		return em, true
	}
	return "", false
}

func (us1 User) AuthUser(log string, pas string)bool{
	if us1.Email == log && us1.Password == pas{
		return true
	}
	return false
}

func (us User) ToJson(seq bool)([]byte){
	var mapVar []byte
	if seq {
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.First_name,
			"last_name":  us.Last_name,
			"birthday":   us.Birthday.Format(time.ANSIC),
			"email":      us.Email,
			"created_at": us.Created_at.Format(time.ANSIC),})
	}else{
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.First_name,
			"last_name":  us.Last_name,
			"birthday":   us.Birthday.Format(time.ANSIC),
			"email":      "*",
			"created_at": us.Created_at.Format(time.ANSIC),})
	}
	return mapVar
}

