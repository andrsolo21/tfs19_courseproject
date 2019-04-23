package user

import (
	"courseproject/internal/users"
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"unicode"
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
		FirstName: us.FirstName,
		LastName:  us.LastName,
		ID:        us.ID,
	}
}

func ConvertDate(us users.UserInp) (us2 users.User, err error) {

	if us.Birthday == "" {
		//us2.Birthday, _ = time.Parse("2006-01-02 15:04:05-07:00", "0001-01-01 00:00:00+00:00")
		us2.Password = us.Password
		us2.Email = us.Email
		us2.LastName = us.LastName
		us2.FirstName = us.FirstName
		return us2, nil
	}

	if len(us.Birthday) == 10 {
		year, err := strconv.Atoi(us.Birthday[:4])
		if err != nil {
			return us2, err
		}

		month, err := strconv.Atoi(us.Birthday[5:7])

		if err != nil {
			return us2, err
		}
		day, err := strconv.Atoi(us.Birthday[8:10])

		if err != nil {
			return us2, err
		}
		us2.Birthday = string(year + month + day)

		us2.Password = us.Password
		us2.Email = us.Email
		us2.LastName = us.LastName
		us2.FirstName = us.FirstName
		//s[:4], s[5:7], s[8:10]
		return us2, nil
	}

	return us2, errors.New("bad date")
}

func CheckDate(date string) error {

	if date[4:5] != "-" && date[7:8] != "-"{
		return errors.New("delimetr is not -")
	}

	for _, el :=range(date[:4] + date[5:7] + date[8:10]){
		if !unicode.IsDigit(el){
			return errors.Errorf("it is not a digit %c", el)
		}
	}
	return nil
}
