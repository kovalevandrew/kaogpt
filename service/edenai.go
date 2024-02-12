package service

import (
	"encoding/json"
	"fmt"
	"kaogpt/config"
	"strings"
)

type EdenaiService struct {
	endpoint    string
	restService RESTService
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

func NewEdenaiService(restService RESTService) *EdenaiService {
	return &EdenaiService{
		endpoint:    "https://api.edenai.run/v2/image/generation",
		restService: restService,
	}
}

func (es *EdenaiService) GetAnswer(message string) (string, error) {
	apiKey := config.GetEdenaiToken()
	payload := strings.NewReader(fmt.Sprintf(`{"response_as_dict":true,"attributes_as_list":false,"show_original_response":false,"resolution":"1024x1024","num_images":1,"text":"%s","providers":"replicate"}`, message))
	response, err := es.restService.Post(es.endpoint, apiKey, payload)
	if err != nil {
		return "", fmt.Errorf("error sending request to Edenai API: %v", err)
	}

	var result ImageEdenai
	if err := json.Unmarshal(response, &result); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	if len(result.Replicate.Items) == 0 {
		return "", fmt.Errorf("no items found in the response")
	}

	return result.Replicate.Items[0].ImageResourceUrl, nil
}
