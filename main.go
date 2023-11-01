package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	DBUsername = "root"
	DBPassword = ""
	DBName     = "go_bot"
)

func main() {
	// Initialize the database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", DBUsername, DBPassword, DBName))
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

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
			// Save the message to the database
			if update.Message.Text != "" {
				err := saveMessageToDB(db, update.Message.Chat.ChatId, update.Message.Text)
				if err != nil {
					log.Println("Error saving message to the database:", err)
				}
			}

			err := respond(botUrl, update)
			if err != nil {
				log.Println("Error sending response:", err)
			}

			offset = update.UpdateId + 1
		}
	}
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
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(botUrl string, update Update) error {
	if update.Message.Text != "" {
		botMessage := BotMessage{
			ChatId: update.Message.Chat.ChatId,
			Text:   "You said: " + update.Message.Text,
		}
		buf, err := json.Marshal(botMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
	}
	return nil
}

func saveMessageToDB(db *sql.DB, chatID int, messageText string) error {
	// Prepare the SQL statement to insert the message into the database
	stmt, err := db.Prepare("INSERT INTO messages (chat_id, text) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the chat ID and message text
	_, err = stmt.Exec(chatID, messageText)
	return err
}
