package storages

import (
	"courseproject/internal/sessionS"
	"courseproject/internal/userS"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type INTT interface {
	AddUser(use interface{}) error
	AddSession(use interface{}) error
	CreateTables() DataB
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

func (db DataB) GetUserByID(id int) (userS.User) {
	var el userS.User

	db.DB.Where("ID = ?", id).First(&el)

	return el
}

func (db DataB) GetSesByToken(id int) (sessionS.Session, error) {
	var el sessionS.Session

	db.DB.Where("session_id = ?", id).First(&el)

	if el.User_id == 0{
		return el, errors.New("this sesion dosen't exist")
	}

	return el, nil
}
