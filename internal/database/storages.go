package storages

import (
	"courseproject/internal/lotS"
	"courseproject/internal/sessionS"
	"courseproject/internal/userS"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type DataB struct {
	DB *gorm.DB
}

type INTT interface {
	CreateTables() DataB
	AddUser(use userS.User) error
	AddSession(use sessionS.Session) error
	CheckEmail(email string) bool
	GetUserByID(id int) (userS.User, error)
	GetSesByToken(token string) (sessionS.Session, error)
	GetUserByEmPass(em string, pass string) (el userS.User, err error)
	ChangeUser(us userS.User, id int) userS.User
	CountUsers() int
	AddLot(l lotS.Lot) error
	GetLots(typ string) (lots []lotS.Lot, d DataB)
	GetUsersLots(id int, role string)(lots []lotS.Lot)
	Db() DataB
	GetLotByID(id int)(el lotS.Lot, err error)
	SMBBuyIT(userID int, l lotS.Lot, price float64) (el lotS.Lot, err error)
}

func NewDataB() (d DataB, err error) {
	user := "gorm_user"
	gormDataBase := "gorm_db"
	pas := "gormDB"

	dsn := "postgres://" + user + ":" + pas + "@localhost:5432/" + gormDataBase +
		"?sslmode=disable&fallback_application_name=fintech-app"
	database, err := gorm.Open("postgres", dsn)

	/*if err != nil {
		fmt.Printf("can't connect to db: %s", err)
	}*/

	return DataB{DB: database}, err
}

func (db DataB) CreateTables() DataB {

	db.DB = db.DB.AutoMigrate(&userS.User{}, &sessionS.Session{}, &lotS.Lot{})

	//db.DB = db.DB.CreateTable(&userS.ShortUser{})

	return db
}

func (db DataB) Db() DataB{
	return db
}

func (db DataB) AddUser(use userS.User) error {
	db.DB = db.DB.Create(&use)

	if err := db.DB.Error; err != nil {
		return err
	}
	return nil
}

func (db DataB) AddSession(use sessionS.Session) error {
	db.DB = db.DB.Create(&use)
	if err := db.DB.Error; err != nil {
		return err
	}
	return nil
}

func (db DataB) CheckEmail(email string) bool {
	var el userS.User

	db.DB.Where("email = ?", email).First(&el)
	if el.ID == 0 {
		return false
	}

	return true
}

func (db DataB) GetUserByID(id int) (userS.User, error) {
	var el userS.User

	db.DB.Where("ID = ?", id).First(&el)

	if el.ID == 0 {
		return el, errors.New("this user dosen't exist")
	}

	return el, nil
}

func (db DataB) GetSesByToken(token string) (sessionS.Session, error) {
	var el sessionS.Session

	db.DB.Where("session_id = ?", token).First(&el)

	if el.User_id == 0 {
		return el, errors.New("this sesion dosen't exist")
	}

	if el.Valid_until.Before(time.Now()) {
		return el, errors.New("expiration date of token is out")
	}

	return el, nil
}

func (db DataB) GetUserByEmPass(em string, pass string) (el userS.User, err error) {

	db.DB.Where("email = ? AND password = ?", em, pass).First(&el)

	if el.ID != 0 {
		return el, nil
	}
	return el, errors.New("this user is not exist")
}

func (db DataB) ChangeUser(us userS.User, id int) userS.User {

	var el userS.User

	db.DB.Where(userS.User{ID: id}).Assign(userS.User{
		FirstName: us.FirstName,
		LastName:  us.LastName,
		Birthday:  us.Birthday,
		UpdatedAt: time.Now(),
	}).FirstOrCreate(&el)

	return el
}

func (db DataB) CountUsers() int {

	var count int
	var users []userS.User
	db.DB.Find(&users).Count(&count)

	return count
}

func (db DataB) AddLot(l lotS.Lot) error {
	db.DB = db.DB.Create(&l)

	if err := db.DB.Error; err != nil {
		return err
	}
	return nil
}

func (db DataB) GetLots(typ string) (lots []lotS.Lot, d DataB) {

	if typ == ""{
		db.DB.Find(&lots)

		return lots, db
	}
	db.DB.Where("status = ?", typ).Find(&lots)

	return lots, db
}

func (db DataB) GetUsersLots(id int, role string)(lots []lotS.Lot){

	switch role{
	case "":
		db.DB.Where("creator_id = ?", id).Or("buyer_id = ?", id).Find(&lots)

	case "own":
		db.DB.Where("creator_id = ?", id).Find(&lots)

	case "buyed":
		db.DB.Where("buyer_id = ?", id).Find(&lots)

	}
	//db.DB.Where("status = ?", typ).Find(&lots)

	return lots
}

func (db DataB) GetLotByID(id int)(el lotS.Lot, err error){
	db.DB.Where("ID = ?", id).First(&el)

	if el.ID == 0 {
		return el, errors.New("this lot dosen't exist")
	}

	return el, err
}

func (db DataB) SMBBuyIT(userID int, l lotS.Lot, price float64) (el lotS.Lot, err error){

	l.BuyPrice = price
	l.BuyerID = userID

	db.DB.Where(lotS.Lot{ID: l.ID}).Assign(l).FirstOrCreate(&el)

	return el, err
}