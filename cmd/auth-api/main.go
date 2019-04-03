package main

import (
	"courseproject/internal/auth"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"log"

)


var data auth.Auth

func main() {

	user := "gorm_user"
	dataB := "gorm_db"
	pas := "gormDB"

	dsn := "postgres://" + user + ":" + pas + "@localhost:5432/" + dataB +
		"?sslmode=disable&fallback_application_name=fintech-app"
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("can't connect to db: %s", err)
	}
	defer db.Close()

	strartListening()

}
