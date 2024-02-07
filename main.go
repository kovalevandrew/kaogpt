package main

import (
	"kaogpt/api"
	"kaogpt/config"
	"kaogpt/di"
	"kaogpt/service"
	"log"
)

func main() {
	// Read configuration
	tokens := config.ReadJsonTokens()

	di.NewRestyClient()

	// Initialize Telegram bot
	telegramBot, err := di.NewTelegramBotClient(tokens.Telegram)
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}
	service.InitTelegramBot(telegramBot)

	// Start API server
	go api.StartBot()

	// Run indefinitely
	select {}

	// Run tests
	// tests.RunTests()
}
