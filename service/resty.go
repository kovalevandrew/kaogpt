package service

import "github.com/go-resty/resty/v2"

type RESTService interface {
	Post(url string, apiKey string, payload interface{}) ([]byte, error)
}

type RestyService struct {
	client *resty.Client
}

func NewRestyService(client *resty.Client) *RestyService {
	return &RestyService{
		client: client,
	}
}

func (rs *RestyService) Post(url string, apiKey string, payload interface{}) ([]byte, error) {
	response, err := rs.client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(apiKey).
		SetBody(payload).
		Post(url)

	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}
