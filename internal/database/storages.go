package storages

import (
	"courseproject/internal/sessionS"
	"courseproject/internal/userS"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type INTT interface {
	AddUser(use userS.User) error
	AddSession(use sessionS.Session) error
	CreateTables() DataB
	CheckEmail(email string) bool
	GetUserByID(id int) (userS.User, error)
	GetSesByToken(id int) (sessionS.Session, error)
	GetUserByEmPass(em string, pass string)(userS.User, error)
	ChangeUser(us userS.User, id int)userS.User
	CountUsers()int
}

type DataB struct {
	DB *gorm.DB
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

	db.DB = db.DB.AutoMigrate(&userS.User{}, &sessionS.Session{})

	//db.DB = db.DB.CreateTable(&userS.ShortUser{})

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

	if el.ID == 0{
		return el, errors.New("this user dosen't exist")
	}

	return el, nil
}

func (db DataB) GetSesByToken(token string) (sessionS.Session, error) {
	var el sessionS.Session

	db.DB.Where("session_id = ?", token).First(&el)

	if el.User_id == 0{
		return el, errors.New("this sesion dosen't exist")
	}

	if el.Valid_until.Before(time.Now()){
		return el, errors.New("expiration date of token is out")
	}

	return el, nil
}

func (db DataB) GetUserByEmPass(em string, pass string)(userS.User, error) {
	var el userS.User

	db.DB.Where("email = ? AND password = ?", em, pass).First(&el)

	if el.ID != 0 {
		return el, nil
	}
	return userS.User{}, errors.New("this user is not exist")
}

func (db DataB) ChangeUser(us userS.User, id int)userS.User{

	var el userS.User

	db.DB.Where(userS.User{ID : id}).Assign(userS.User{
		FirstName:us.FirstName,
		LastName:us.LastName,
		Birthday:us.Birthday,
		UpdatedAt:time.Now(),
	}).FirstOrCreate(&el)

	return el
}

func (db DataB) CountUsers()int{

	var count int
	var users []userS.User
	db.DB.Find(&users).Count(&count)

	return count
}
