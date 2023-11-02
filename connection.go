package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

var db *gorm.DB

const (
	DBUsername = "root"
	DBPassword = ""
	DBName     = "go_bot"
)

func connection() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", DBUsername, DBPassword, DBName))
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	//defer db.Close()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Question{})
	seedQuestions(db)
	db.AutoMigrate(&Message{})
}
