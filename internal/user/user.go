package user

import (
	"time"
)

type User struct {
	ID         int       `json:"id"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Birthday   time.Time `json:"birthday"` //
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time //`json:"created_at"`
	Updated_at time.Time //`json:"updated_at"`
}

func (us1 User) CheckUser(us2 User)bool{
	if us1.Email == us2.Email {
		return false
	}
	return true
}

func (us1 User) AuthUser(log string, pas string)bool{
	if us1.Email == log && us1.Password == pas{
		return true
	}
	return false
}

type Lot struct {
	ID          int       `json:"id"`
	Creator_id  int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"` //
	Min_price   float64   `json:"min_price"`
	Price_step  float64   `json:"price_step"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}
