package service

import (
	"encoding/json"
	"fmt"
	"kaogpt/config"

	"github.com/go-resty/resty/v2"
)

type ChatGPTService struct {
	apiEndpoint string
	restyClient *resty.Client
}

type GPTRequest struct {
	Model     string        `json:"model"`
	Messages  []MessageItem `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
}

type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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

func NewChatGPTService(restyClient *resty.Client) *ChatGPTService {
	return &ChatGPTService{
		apiEndpoint: "https://api.openai.com/v1/chat/completions",
		restyClient: restyClient,
	}
}

func (cs *ChatGPTService) GetGptAnswer(message string) (string, error) {
	apiKey := config.GetGPTToken()
	requestBody := GPTRequest{
		Model:     "gpt-3.5-turbo",
		Messages:  []MessageItem{{Role: "system", Content: message}},
		MaxTokens: 500,
	}

	response, err := cs.restyClient.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
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
