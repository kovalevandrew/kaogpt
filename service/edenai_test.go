package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnswerEdenai(t *testing.T) {
	token := "edenai mock token"
	// Create a mock RESTful service
	mockRestyService := &MockRestyService{
		ResponseBody: []byte(`{"replicate":{"items":[{"image":"mock_image_url","image_resource_url":"mock_resource_url"}],"status":"mock_status"}}`),
		Err:          nil,
	}

	// Create an EdenaiService instance with the mock RESTful service
	edenaiService := NewEdenaiService(mockRestyService)

	// Test case: Mock response received successfully
	expectedURL := "mock_resource_url"
	message := "Test message"
	imageURL, err := edenaiService.GetAnswer(message, token)
	assert.NoError(t, err, "GetAnswer returned an unexpected error")
	assert.Equal(t, expectedURL, imageURL, "GetAnswer returned an unexpected imageURL")

	// Test case: Error returned from mock service
	mockRestyService.Err = fmt.Errorf("Mock error")
	_, err = edenaiService.GetAnswer(message, token)
	assert.Error(t, err, "Expected an error from GetAnswer")
}
