package service

type MockRestyService struct {
	ResponseBody []byte
	Err          error
}

func (m *MockRestyService) Post(url string, apiKey string, payload interface{}) ([]byte, error) {
	return m.ResponseBody, m.Err
}
