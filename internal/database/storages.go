package storages

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type INTT interface {
	AddUser(use interface{}) (error)
	AddSession(use interface{}) (error)
}

type DataB struct {
	DB *gorm.DB
}

type UserDB struct {
	gorm.Model

	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
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

	database.CreateTable(&UserDB{})

	return DataB{DB: database}, err
}

func (db DataB) AddUser(use interface{}) (error) {
	db.DB = db.DB.Create(&use)
	if err := db.DB.Error; err != nil {
		return err
	}
	return nil
}

func (db DataB) AddSession(use interface{}) (error) {
	db.DB = db.DB.Create(&use)
	if err := db.DB.Error; err != nil {
		return err
	}
	return nil
}
