package di

import (
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewRestyClient() *resty.Client {
	return resty.New()
}

func NewTelegramBotClient(token string) (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(token)
}
