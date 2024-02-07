package main

import (
	"kaogpt/config"
	"kaogpt/di"
	"kaogpt/service"
	"log"
)

func main() {
	// Read configuration
	tokens := config.ReadJsonTokens()

	// Initialize dependencies
	restyClient := di.NewRestyClient()
	telegramBot, err := di.NewTelegramBotClient(tokens.Telegram)
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}

	// Initialize services
	chatGPTService := service.NewChatGPTService(restyClient)
	edenaiService := service.NewEdenaiService()

	// Initialize Telegram service
	telegramService := service.NewTelegramService(telegramBot, chatGPTService, edenaiService)

	// Start Telegram service
	telegramService.Start()

	// Run indefinitely
	select {}
}
