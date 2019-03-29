package main

import (
	"courseproject/internal/auth"
	"courseproject/vendor/github.com/jinzhu/gorm/gorm"
	_ "courseproject/vendor/github.com/jinzhu/gorm/dialects/postgres"
	"log"
)


var data auth.Auth

func main() {
	dsn := "postgres://user:passwd@localhost:5432/fintech" +
		"?sslmode=disable&fallback_application_name=fintech-app"
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("can't connect to db: %s", err)
	}
	defer db.Close()

	strartListening()

}
