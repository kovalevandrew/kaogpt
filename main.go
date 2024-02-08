package main

import (
	"kaogpt/service"

	"github.com/go-resty/resty/v2"
)

func main() {
	restyClient := resty.New()

	chatGPTService := service.NewChatGPTService(restyClient)
	edenaiService := service.NewEdenaiService(restyClient)
	telegramService := service.NewTelegramService(chatGPTService, edenaiService)

	telegramService.Start()

	// Run indefinitely
	select {}
}
