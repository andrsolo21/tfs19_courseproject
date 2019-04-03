package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

func init(){
	user := "gorm_user"
	dataB := "gorm_db"
	pas := "gormDB"

	dsn := "postgres://" + user + ":" + pas + "@localhost:5432/" + dataB +
		"?sslmode=disable&fallback_application_name=fintech-app"
	DB, err := gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("can't connect to db: %s", err)
	}
}
