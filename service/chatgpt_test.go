package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnswerGpt(t *testing.T) {
	token := "mocked token for gpt"
	mockRestyService := &MockRestyService{
		ResponseBody: []byte(`{"choices":[{"message":{"content":"Mock response"}}]}`),
		Err:          nil,
	}

	chatGPTService := NewChatGPTService(mockRestyService)

	expectedAnswer := "Mock response"
	message := "Test message"
	answer, err := chatGPTService.GetAnswer(message, token)
	assert.NoError(t, err, "GetAnswer returned an unexpected error")
	assert.Equal(t, expectedAnswer, answer, "GetAnswer returned an unexpected answer")

	mockRestyService.Err = fmt.Errorf("Mock error")
	_, err = chatGPTService.GetAnswer(message, token)
	assert.Error(t, err, "Expected an error from GetAnswer")
}
