package user

import (
	"encoding/json"
	"strconv"
	"time"
	"unicode"

	"gitlab.com/andrsolo21/courseproject/internal/users"

	"github.com/pkg/errors"
)

func ToJSON(seq bool, us users.User) []byte {
	var mapVar []byte
	if seq {
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.FirstName,
			"last_name":  us.LastName,
			"birthday":   us.Birthday, //.Format(time.ANSIC),
			"email":      us.Email,
			"created_at": us.CreatedAt.Format(time.ANSIC)})
	} else {
		mapVar, _ = json.Marshal(map[string]string{
			"id":         strconv.Itoa(us.ID),
			"first_name": us.FirstName,
			"last_name":  us.LastName,
			"birthday":   us.Birthday, //.Format(time.ANSIC),
			"email":      "*",
			"created_at": us.CreatedAt.Format(time.ANSIC)})
	}
	return mapVar
}

func ToShort(us users.User) users.ShortUser {
	return users.ShortUser{
		FirstName: users.U(us.FirstName),
		LastName:  us.LastName,
		ID:        us.ID,
	}
}

func CheckDate(date string) error {

	if date == "" {
		return nil
	}

	if len(date) != 10 {
		return errors.New("can't parse date")
	}
	if date[4:5] != "-" && date[7:8] != "-" {
		return errors.New("delimetr is not -")
	}

	for _, el := range date[:4] + date[5:7] + date[8:10] {
		if !unicode.IsDigit(el) {
			return errors.Errorf("it is not a digit %c", el)
		}
	}

	date += " 15:04:05-07:00"

	_, err := time.Parse("2006-01-02 15:04:05-07:00", date)

	if err != nil {
		return errors.New("this month or day doesn't exist")
	}

	return nil
}
