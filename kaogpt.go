package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatId int64

type Tokens struct {
	Telegram string
	Gpt      string
	Edenai   string
}

type ImageEdenai struct {
	Replicate struct {
		Items []struct {
			Image            string `json:"image"`
			ImageResourceUrl string `json:"image_resource_url"`
		} `json:"items"`
		Status string `json:"status"`
	} `json:"replicate"`
}

const (
	apiEndpoint    = "https://api.openai.com/v1/chat/completions"
	edenaiEndpoint = "https://api.edenai.run/v2/image/generation"
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

func sendPicture(update *tgbotapi.Update) {
	waitMsg := tgbotapi.NewMessage(chatId, getImageWaitAnswer())
	waitMsg.ReplyToMessageID = update.Message.MessageID
	bot.Send(waitMsg)
	msg := tgbotapi.NewMessage(chatId, getEdenaiImage(update))
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

func getEdenaiImage(update *tgbotapi.Update) string {
	payload := strings.NewReader("{\"response_as_dict\":true,\"attributes_as_list\":false,\"show_original_response\":false,\"resolution\":\"1024x1024\",\"num_images\":1,\"text\":\"" + update.Message.Text + "\",\"providers\":\"replicate\"}")
	req, _ := http.NewRequest("POST", edenaiEndpoint, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+getEdenaiToken()+"")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var result ImageEdenai
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result.Replicate.Items[0].ImageResourceUrl
}

func getWaitAnswer() string {
	return "Wait your answer is prepearing"
}

func getImageWaitAnswer() string {
	return "Wait your image is generating"
}

func getTelegramToken() string {
	tokens := readJsonTokens()
	return tokens.Telegram
}

func getGPTToken() string {
	tokens := readJsonTokens()
	return tokens.Gpt
}

func getEdenaiToken() string {
	tokens := readJsonTokens()
	return tokens.Edenai
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
func isMessageRequestForEdenai(update *tgbotapi.Update) bool {
	if strings.Contains(update.Message.Text, "generat") {
		return true
	}

	return false
}

func main() {
	connectWithTelegram()

	updateConfig := tgbotapi.NewUpdate(0)
	for update := range bot.GetUpdatesChan(updateConfig) {
		if update.Message != nil && update.Message.Text == "/start" {
			chatId = update.Message.Chat.ID
			sendMessage("If you want get image type \"generate\" firstly, or ask your question and ai will respons to you.")
		}
		if isMessageRequestForEdenai(&update) {
			sendAnswer(&update)
		}
		if isMessageRequestForGPT(&update) {
			sendAnswer(&update)
		}
	}
}
