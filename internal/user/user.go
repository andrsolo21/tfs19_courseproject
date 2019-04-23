package user

import (
	"courseproject/internal/users"
	"encoding/json"
	"strconv"
	"time"
)

func ToJSON(seq bool, us users.User) []byte {
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

func ToShort(us users.User) users.ShortUser {
	return users.ShortUser{
		FirstName: us.FirstName,
		LastName:  us.LastName,
		ID:        us.ID,
	}
}

func ConvertDate(us users.UserInp) (us2 users.User){

	if us.Birthday == ""{
		us2.Birthday, _ = time.Parse("2006-01-02 15:04:05-07:00","0001-01-01 00:00:00+00:00")
	}

	return
}