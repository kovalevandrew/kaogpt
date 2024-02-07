package tests

import (
	"kaogpt/service"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestGetGptAnswer(t *testing.T) {
	// Provide a sample update
	update := &tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
		},
	}

	// Call the function under test
	answer := service.GetGptAnswer(update)

	// Check if the answer is not empty
	if answer == "" {
		t.Error("Expected non-empty answer, got empty")
	}
}
