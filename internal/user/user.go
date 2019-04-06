package user

import (
	"courseproject/internal/userS"
	"encoding/json"
	"strconv"
	"time"
)

func ToJson(seq bool, us userS.User) []byte {
	var mapVar []byte
	if seq {
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.FirstName,
			"last_name":  us.LastName,
			"birthday":   us.Birthday.Format(time.ANSIC),
			"email":      us.Email,
			"created_at": us.CreatedAt.Format(time.ANSIC)})
	} else {
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.FirstName,
			"last_name":  us.LastName,
			"birthday":   us.Birthday.Format(time.ANSIC),
			"email":      "*",
			"created_at": us.CreatedAt.Format(time.ANSIC)})
	}
	return mapVar
}

func ToShort(us userS.User) userS.ShortUser {
	return userS.ShortUser{
		FirstName: us.FirstName,
		LastName:  us.LastName,
		ID:        us.ID,
	}
}
