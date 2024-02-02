package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatId int64

type Tokens struct {
	Telegram string
	Gpt      string
}

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

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
	waitMsg.ReplyToMessageID = update.Message.MessageID
	bot.Send(waitMsg)
	msg := tgbotapi.NewMessage(chatId, getGptAnswer(update))
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func getGptAnswer(update *tgbotapi.Update) string {
	apiKey := getGPTToken()
	client := resty.New()

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": update.Message.Text}},
			"max_tokens": 500,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "Error while decoding JSON response:" + err.Error()
	}

	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content
}

func getWaitAnswer() string {
	return "Wait your answer is prepearing"
}

func getTelegramToken() string {
	tokens := readJsonTokens()
	return tokens.Telegram
}

func getGPTToken() string {
	tokens := readJsonTokens()
	return tokens.Gpt
}

func readJsonTokens() Tokens {
	// read file
	jsonData, err := os.ReadFile("./token/tokens.json")
	if err != nil {
		fmt.Print(err)
	}
	var payload Tokens

	err = json.Unmarshal(jsonData, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
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
