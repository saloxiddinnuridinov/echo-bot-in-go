package main

import "github.com/jinzhu/gorm"

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	gorm.Model
	Chat   Chat   `json:"chat"`
	Text   string `json:"text"`
	ChatID int
}

type Chat struct {
	ChatId    int `json:"id"`
	Username  string
	FirstName string
	LastName  string
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
	TelegramUserID int
	Username       string
	FirstName      string
	LastName       string
	ChatId         int

	// Add more fields as needed
}
