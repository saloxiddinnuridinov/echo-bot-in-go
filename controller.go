package main

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

func seedQuestions(db *gorm.DB) {
	// Question ma'lumotlari
	Questions := []Question{
		{Title: "Savol 1", Position: 0},
		{Title: "Savol 2", Position: 1},
		{Title: "Savol 3", Position: 2},
		{Title: "Savol 4", Position: 3},
		{Title: "Savol 5", Position: 4},
		{Title: "Savol 6", Position: 5},
		{Title: "Savol 7", Position: 6},
		{Title: "Savol 8", Position: 7},
	}

	// Question jadvalini tozalash
	db.Exec("TRUNCATE TABLE Questions")

	// Question ma'lumotlarini jadvalga qo'shish
	for _, Question := range Questions {
		db.Create(&Question)
	}
}

func saveMessageToDB(db *gorm.DB, fromMessage From, chatID Chat, messageText string, QuestionId int) error {
	message := Message{
		From:       fromMessage,
		Chat:       chatID,
		Text:       messageText,
		QuestionId: QuestionId,
		ChatId:     chatID.ChatId,
	}

	//bazasiga saqlash
	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}

func getQuestionIDByPosition(db *gorm.DB, position int) (int, error) {
	var question Question
	if err := db.Where("position = ?", position).First(&question).Error; err != nil {
		return 0, err
	}
	return int(question.ID), nil
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

func sendQuestion(botUrl string, chatID int, Question *Question) error {
	botMessage := BotMessage{
		ChatId: chatID,
		Text:   Question.Title,
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(botUrl string, chatID int) error {
	botMessage := BotMessage{
		ChatId: chatID,
		Text:   "Rahmat! Siz barcha savollarga javob berdingiz",
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
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

func checkUser(db *gorm.DB, chatID int) *User {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)
	if user.ID == 0 {
		return nil // User not found
	}
	return &user
}

func getQuestionsFromDB(db *gorm.DB) ([]Question, error) {
	var Questions []Question
	if err := db.Order("position").Find(&Questions).Error; err != nil {
		return nil, err
	}
	return Questions, nil
}

func getFirstQuestionFromDB(db *gorm.DB) (*Question, error) {
	var Question Question
	if err := db.Where("Position = ?", 0).First(&Question).Error; err != nil {
		return nil, err
	}
	return &Question, nil
}

func getSecondQuestionFromDB(db *gorm.DB, value int) (*Question, error) {
	var Question Question
	if err := db.Where("Position = ?", value).First(&Question).Error; err != nil {
		return nil, err
	}
	return &Question, nil
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]int),
	}
}

func (c *Cache) Set(key string, value int) {
	c.mu.Lock()
	c.cache[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (int, bool) {
	c.mu.Lock()
	value, exists := c.cache[key]
	c.mu.Unlock()
	return value, exists
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}

//func deleteFromCache(cache map[string]int, questionID string) {
//	delete(cache, questionID)
//}
