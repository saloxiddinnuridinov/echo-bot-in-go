package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

const (
	DBUsername = "root"
	DBPassword = ""
	DBName     = "go_bot"
)

func main() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", DBUsername, DBPassword, DBName))
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()
	db.AutoMigrate(&User{})

	botToken := "6166118537:AAFrnPz7lOcqqMPiIQA2h1fLmPGL6L0qR0c"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0

	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Error fetching updates:", err)
		}

		for _, update := range updates {
			if update.Message.Text == "/start" {

				id := update.Message.Chat.ChatId
				fmt.Println(id)

				// Check if the user is already in the database
				user := checkUser(db, id)

				fmt.Println(updates)
				if user == nil {
					// User is not in the database, save their data
					saveUserToDB(db, update.Message.Chat)
					if err != nil {
						log.Println("Foydalanuvchini maʼlumotlar bazasiga saqlashda xatolik yuz berdi:", err)
					}
				}

				saveMessageToDB(db, update.Message.From, update.Message.Chat, update.Message.Text)

				err := respond(botUrl, update)
				if err != nil {
					log.Println("Error sending response:", err)
				}
				// Handle the /start command
				//err := handleStartCommand(botUrl, update)
				//if err != nil {
				//	log.Println("Error handling the /start command:", err)
				//}
			} else {
				// Handle other messages
				user := saveMessageToDB(db, update.Message.From, update.Message.Chat, update.Message.Text)
				if user == nil {
					if err != nil {
						log.Println("Error saving user to the database:", err)
					}
				}
				err := respond(botUrl, update)
				if err != nil {
					log.Println("Error sending response:", err)
				}
			}

			offset = update.UpdateId + 1
		}
	}
}

//func handleStartCommand(url string, update Update) interface{} {
//	return nil
//}

func checkUser(db *gorm.DB, chatID int) *User {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)
	if user.ID == 0 {
		return nil // User not found
	}
	return &user
}

func saveUserToDB(db *gorm.DB, chat Chat) error {
	user := User{
		ChatId:    chat.ChatId,
		Username:  chat.Username,
		FirstName: chat.FirstName,
		LastName:  chat.LastName,
	}
	return db.Create(&user).Error
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err // Xatolik sodir bo'ldi
	}
	return restResponse.Result, nil
}

func respond(botUrl string, update Update) error {
	if update.Message.Text != "" {
		botMessage := BotMessage{
			ChatId: update.Message.Chat.ChatId,
			Text:   "Sizning xabaringiz: " + update.Message.Text,
		}
		buf, err := json.Marshal(botMessage)
		if err != nil {
			return err // Xatolik sodir bo'ldi
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err // Xatolik sodir bo'ldi
		}
	}
	return nil
}

func saveMessageToDB(db *gorm.DB, fromMessage From, chatID Chat, messageText string) error {
	message := Message{
		From: fromMessage,
		Chat: chatID,
		Text: messageText,
	}

	// Ma'lumotlar bazasiga saqlash
	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}
