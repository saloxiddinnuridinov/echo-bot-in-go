package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	connection()
	botToken := "6166118537:AAFrnPz7lOcqqMPiIQA2h1fLmPGL6L0qR0c"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	cache := NewCache()
	//isWaitingForSecondQuestion := false

	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Error fetching updates:", err)
		}

		for _, update := range updates {
			if update.Message.Text == "/start" {

				id := update.Message.Chat.ChatId
				//fmt.Println(id)
				user := checkUser(db, id)
				if user == nil {
					saveUserToDB(db, update.Message.Chat)
					if err != nil {
						log.Println("Foydalanuvchini maʼlumotlar bazasiga saqlashda xatolik yuz berdi:", err)
					}
				}

				Question, err := getFirstQuestionFromDB(db)
				if err != nil {
					log.Println("Error getting the first Question:", err)
				}

				sendQuestion(botUrl, update.Message.Chat.ChatId, Question)

				cache.Set("update.Message.Chat.ChatId", Question.Position+1)

			} else {

				userAnswer := update.Message.Text
				value, exists := cache.Get("update.Message.Chat.ChatId")

				if exists {

					questionID, err := getQuestionIDByPosition(db, value-1)
					if err != nil {
						log.Println("Question position orqali Idsini olishda: ", err)
					}
					fmt.Println(questionID)
					//userID, err := getUserIDByChatID(db)
					err = saveMessageToDB(db, update.Message.From, update.Message.Chat, userAnswer, questionID)
					if err != nil {
						log.Println("Error saving message to the database:", err)
					}
					Question, err := getSecondQuestionFromDB(db, value)
					if err != nil {
						sendMessage(botUrl, update.Message.Chat.ChatId)
						cache.Delete("update.Message.Chat.ChatId")
						//fmt.Println("bazada savol tugadi")
						//log.Println("Error getting the second Question:", err)
					} else {
						cache.Delete("update.Message.Chat.ChatId")
						sendQuestion(botUrl, update.Message.Chat.ChatId, Question)
						cache.Set("update.Message.Chat.ChatId", Question.Position+1)
					}

				}
			}

			offset = update.UpdateId + 1
		}
	}
}
