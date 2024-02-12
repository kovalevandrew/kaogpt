package main

import (
	"kaogpt/service"

	"github.com/go-resty/resty/v2"
)

func main() {
	restyClient := resty.New()
	restService := service.NewRestyService(restyClient)

	chatGPTService := service.NewChatGPTService(restService)
	edenaiService := service.NewEdenaiService(restService)
	telegramService := service.NewTelegramService(chatGPTService, edenaiService)

	telegramService.StartBot()

	// Run indefinitely
	select {}
}
