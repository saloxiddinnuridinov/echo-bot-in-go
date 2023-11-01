package main

import "github.com/jinzhu/gorm"

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	gorm.Model
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type From struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ChatId    int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type User struct {
	gorm.Model
	Phone     int
	Username  string
	FirstName string
	LastName  string
	ChatId    int

	// Add more fields as needed
}
