package service

import (
	"encoding/json"
	"fmt"
	"kaogpt/config"
	"strings"

	"github.com/go-resty/resty/v2"
)

type EdenaiService struct {
	endpoint    string
	restyClient *resty.Client
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

func NewEdenaiService(restyClient *resty.Client) *EdenaiService {
	return &EdenaiService{
		endpoint:    "https://api.edenai.run/v2/image/generation",
		restyClient: restyClient,
	}
}

func (es *EdenaiService) GetEdenaiImage(message string) (string, error) {
	payload := strings.NewReader(fmt.Sprintf(`{"response_as_dict":true,"attributes_as_list":false,"show_original_response":false,"resolution":"1024x1024","num_images":1,"text":"%s","providers":"replicate"}`, message))
	response, err := es.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+config.GetEdenaiToken()).
		SetBody(payload).
		Post(es.endpoint)

	if err != nil {
		return "", fmt.Errorf("error sending request to Edenai API: %v", err)
	}

	var result ImageEdenai
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	if len(result.Replicate.Items) == 0 {
		return "", fmt.Errorf("no items found in the response")
	}

	return result.Replicate.Items[0].ImageResourceUrl, nil
}
