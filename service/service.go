package service

import (
	"encoding/json"
	"fmt"
	"io"
	"kaogpt/config"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	telegramBot *tgbotapi.BotAPI
)

const (
	apiEndpoint    = "https://api.openai.com/v1/chat/completions"
	edenaiEndpoint = "https://api.edenai.run/v2/image/generation"
)

type ImageEdenai struct {
	Replicate struct {
		Items []struct {
			Image            string `json:"image"`
			ImageResourceUrl string `json:"image_resource_url"`
		} `json:"items"`
		Status string `json:"status"`
	} `json:"replicate"`
}

type TelegramBotService struct {
	Bot         *tgbotapi.BotAPI
	RestyClient *resty.Client
	TelegramKey string
	GptKey      string
	EdenaiKey   string
}

type GPTResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
}

func NewTelegramBotService(bot *tgbotapi.BotAPI, restyClient *resty.Client, telegramKey, gptKey, edenaiKey string) *TelegramBotService {
	return &TelegramBotService{
		Bot:         bot,
		RestyClient: restyClient,
		TelegramKey: telegramKey,
		GptKey:      gptKey,
		EdenaiKey:   edenaiKey,
	}
}

func GetWaitAnswer() string {
	return "Wait, your answer is preparing."
}

func GetImageWaitAnswer() string {
	return "Wait, your image is generating."
}

func GetGptAnswer(update *tgbotapi.Update) string {
	apiKey := config.GetGPTToken()
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
		log.Fatalf("Error while sending the request: %v", err)
	}

	// Check if the response is nil or empty
	if response == nil || response.Body() == nil {
		log.Fatalf("Received a nil response or body from the server")
	}

	var gptResponse GPTResponse
	if err := json.Unmarshal(response.Body(), &gptResponse); err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	// Check if choices are present and not empty
	if len(gptResponse.Choices) == 0 {
		log.Fatalf("No choices found in the response")
	}

	// Get the content from the first choice
	content := gptResponse.Choices[0].Message.Content
	return content
}

func GetEdenaiImage(update *tgbotapi.Update) string {
	payload := strings.NewReader("{\"response_as_dict\":true,\"attributes_as_list\":false,\"show_original_response\":false,\"resolution\":\"1024x1024\",\"num_images\":1,\"text\":\"" + update.Message.Text + "\",\"providers\":\"replicate\"}")
	req, _ := http.NewRequest("POST", edenaiEndpoint, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+config.GetEdenaiToken()+"")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to Edenai API: %v", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var result ImageEdenai
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	if len(result.Replicate.Items) == 0 {
		log.Fatalf("No items found in the response")
	}

	return result.Replicate.Items[0].ImageResourceUrl
}

func InitTelegramBot(bot *tgbotapi.BotAPI) {
	telegramBot = bot
}

func GetTelegramBot() *tgbotapi.BotAPI {
	return telegramBot
}
