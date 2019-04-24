package auth

import (
	"gitlab.com/andrsolo21/courseproject/internal/lots"
	"gitlab.com/andrsolo21/courseproject/internal/sessions"
	"gitlab.com/andrsolo21/courseproject/internal/storages"
	"gitlab.com/andrsolo21/courseproject/internal/users"
	"gitlab.com/andrsolo21/courseproject/pkg/log"

	"testing"
)

type testpair struct {
	value users.User
	ierr  int
	m     MOK
	out   bool
}

func TestCheckUser(t *testing.T) {

	var tests = []testpair{
		{
			users.User{
				Email:    "123",
				Password: "456",
			}, 201, MOK{"123"}, true,
		},
		{
			users.User{
				Email:    "",
				Password: "456",
			}, 400, MOK{"123"}, false,
		},
		{
			users.User{
				Email:    "123",
				Password: "",
			}, 400, MOK{"123"}, false,
		},
		{
			users.User{
				Email:    "123",
				Password: "456",
			}, 409, MOK{"789"}, false,
		},
	}
	for i, pair := range tests {
		ierr, err := CheckUser(pair.value, pair.m)
		if (err == nil) != pair.out && (ierr == pair.ierr) {
			t.Error(
				//"For", pair.value,
				"expected", pair.out,
				"got", pair.out,
				"test ID: ", i,
			)
		}
	}
}

type MOK struct {
	email string
}

func (n MOK) CheckEmail(em string) bool {
	return n.email == em
}
func (n MOK) CreateTables() storages.DataB {
	return storages.DataB{}
}
func (n MOK) AddUser(use users.User) error {
	return nil
}
func (n MOK) AddSession(use sessions.Session) error {
	return nil
}
func (n MOK) GetUserByID(id int) (users.User, error) {
	return users.User{}, nil
}
func (n MOK) GetSesByToken(token string) (sessions.Session, error) {
	return sessions.Session{}, nil
}
func (n MOK) GetUserByEmPass(em string, pass string) (el users.User, err error) {
	return el, err
}
func (n MOK) ChangeUser(us users.User, id int) users.User {
	return users.User{}
}
func (n MOK) CountUsers() int {
	return 0
}
func (n MOK) AddLot(l lots.Lot) (lots.Lot, error) {
	return lots.Lot{}, nil
}
func (n MOK) GetLots(typ string) (lots []lots.Lot, d storages.DataB) {
	return lots, d
}
func (n MOK) GetUsersLots(id int, role string) (lots []lots.Lot) {
	return lots
}
func (n MOK) Db() storages.DataB {
	return storages.DataB{}
}
func (n MOK) GetLotByID(id int) (el lots.Lot, err error) {
	return el, err
}
func (n MOK) SMBBuyIT(userID int, l lots.Lot, price float64) (el lots.Lot, err error) {
	return el, err
}
func (n MOK) UpdateLot(l lots.Lot, id int) lots.Lot {
	return lots.Lot{}
}
func (n MOK) KillBadLots(logger log.Logger) {

}
func (n MOK) DeleteCrLor(id int) {
}
func (n MOK) CheckDelLot(id int) bool {
	return true
}
