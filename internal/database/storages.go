package storages

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type INTT interface {
	AddUser(use interface{})(error)
}

type DataB struct {
	DB *gorm.DB
}

func NewDataB() (d DataB, err error) {
	user := "gorm_user"
	dataB := "gorm_db"
	pas := "gormDB"

	dsn := "postgres://" + user + ":" + pas + "@localhost:5432/" + dataB +
		"?sslmode=disable&fallback_application_name=fintech-app"
	database, err := gorm.Open("postgres", dsn)
	/*if err != nil {
		fmt.Printf("can't connect to db: %s", err)
	}*/
	d = DataB{
		DB: database,
	}
	return d, err
}

func (db DataB) AddUser(use interface{})(error){
	db.DB = db.DB.Create(use)
	if err := db.DB.Error; err!= nil{
		return err
	}
	return nil
}
