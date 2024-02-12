package service

import (
	"fmt"
	"testing"
)

func TestGetAnswerGpt(t *testing.T) {
	token := "moked token for gpt"
	// Create a mock RESTful service
	mockRestyService := &MockRestyService{
		ResponseBody: []byte(`{"choices":[{"message":{"content":"Mock response"}}]}`),
		Err:          nil,
	}

	// Create a ChatGPTService instance with the mock RESTful service
	chatGPTService := NewChatGPTService(mockRestyService)

	// Test case: Mock response received successfully
	expectedAnswer := "Mock response"
	message := "Test message"
	answer, err := chatGPTService.GetAnswer(message, token)
	if err != nil {
		t.Errorf("GetAnswer returned error: %v", err)
	}
	if answer != expectedAnswer {
		t.Errorf("GetAnswer returned %q, expected %q", answer, expectedAnswer)
	}

	// Test case: Error returned from mock service
	mockRestyService.Err = fmt.Errorf("Mock error")
	_, err = chatGPTService.GetAnswer(message, token)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
