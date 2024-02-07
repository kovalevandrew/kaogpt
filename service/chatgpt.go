package service

import (
	"encoding/json"
	"fmt"
	"kaogpt/config"

	"github.com/go-resty/resty/v2"
)

type GPTResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
}

type ChatGPTService struct {
	apiEndpoint string
	restyClient *resty.Client
}

func NewChatGPTService(restyClient *resty.Client) *ChatGPTService {
	return &ChatGPTService{
		apiEndpoint: "https://api.openai.com/v1/chat/completions",
		restyClient: restyClient,
	}
}

func (cs *ChatGPTService) GetGptAnswer(message string) (string, error) {
	apiKey := config.GetGPTToken()
	client := resty.New()
	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(struct {
			Model     string        `json:"model"`
			Messages  []interface{} `json:"messages"`
			MaxTokens int           `json:"max_tokens"`
		}{
			Model:     "gpt-3.5-turbo",
			Messages:  []interface{}{map[string]interface{}{"role": "system", "content": message}},
			MaxTokens: 500,
		}).
		Post(cs.apiEndpoint)
	if err != nil {
		return "", fmt.Errorf("error while sending the request: %v", err)
	}
	var gptResponse GPTResponse
	if err := json.Unmarshal(response.Body(), &gptResponse); err != nil {
		return "", fmt.Errorf("error decoding JSON response: %v", err)
	}
	if len(gptResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices found in the response")
	}
	return gptResponse.Choices[0].Message.Content, nil
}
