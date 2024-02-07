package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"kaogpt/config"
)

type EdenaiService struct {
	endpoint string
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

func NewEdenaiService() *EdenaiService {
	return &EdenaiService{endpoint: "https://api.edenai.run/v2/image/generation"}
}

func (es *EdenaiService) GetEdenaiImage(message string) (string, error) {
	payload := strings.NewReader(fmt.Sprintf(`{"response_as_dict":true,"attributes_as_list":false,"show_original_response":false,"resolution":"1024x1024","num_images":1,"text":"%s","providers":"replicate"}`, message))
	req, err := http.NewRequest("POST", es.endpoint, payload)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+config.GetEdenaiToken())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to Edenai API: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	var result ImageEdenai
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON response: %v", err)
	}
	if len(result.Replicate.Items) == 0 {
		return "", fmt.Errorf("no items found in the response")
	}
	return result.Replicate.Items[0].ImageResourceUrl, nil
}
