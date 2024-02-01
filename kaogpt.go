package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatId int64

func connectWithTelegram() {
	var TOKEN = getTelegramToken()

	var err error
	if bot, err = tgbotapi.NewBotAPI(TOKEN); err != nil {
		panic("Cannot connect to Telegram API!")
	}
}

func sendMessage(msg string) {
	msgConfig := tgbotapi.NewMessage(chatId, msg)
	bot.Send(msgConfig)
}

func sendAnswer(update *tgbotapi.Update) {
	waitMsg := tgbotapi.NewMessage(chatId, getWaitAnswer())
	msg := tgbotapi.NewMessage(chatId, getGptAnswer())
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(waitMsg)
	time.Sleep(2 * time.Second)
	bot.Send(msg)
}

func getGptAnswer() string {
	return "Hello from GPT chat!"
}

func getWaitAnswer() string {
	return "Wait your answer is prepearing"
}

type TokenData struct {
	Token string
}

func getTelegramToken() string {
	// read file
	data, err := os.ReadFile("./token/tgtoken.json")
	if err != nil {
		fmt.Print(err)
	}

	// Now let's unmarshall the data into `payload`
	var payload TokenData
	err = json.Unmarshal(data, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return payload.Token
}

func isMessageRequestForGPT(update *tgbotapi.Update) bool {
	if update.Message.Text == "/start" || update.Message == nil || update.Message.Text == "" {
		return false
	}

	return true
}

func main() {
	connectWithTelegram()

	updateConfig := tgbotapi.NewUpdate(0)
	for update := range bot.GetUpdatesChan(updateConfig) {
		if update.Message != nil && update.Message.Text == "/start" {
			chatId = update.Message.Chat.ID
			sendMessage("Ask you question here and chatGPT will respond")
		}

		if isMessageRequestForGPT(&update) {

			sendAnswer(&update)
		}
	}
}
