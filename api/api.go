// api/handlers.go
package api

import (
	"kaogpt/service"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot() {
	bot := service.GetTelegramBot()
	updateConfig := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			chatID := update.Message.Chat.ID
			SendMessage("If you want to get an image, type \"generate\" firstly, or ask your question and AI will respond to you.", chatID)
		}
		if isMessageRequestForEdenai(&update) {
			SendPicture(&update)
			return
		}
		if isMessageRequestForGPT(&update) {
			SendAnswer(&update)
		}
	}
}

func SendMessage(msg string, chatID int64) {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	service.GetTelegramBot().Send(msgConfig)
}

func SendAnswer(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	waitMsg := tgbotapi.NewMessage(chatID, service.GetWaitAnswer())
	waitMsg.ReplyToMessageID = update.Message.MessageID
	service.GetTelegramBot().Send(waitMsg)

	msg := tgbotapi.NewMessage(chatID, service.GetGptAnswer(update))
	msg.ReplyToMessageID = update.Message.MessageID
	service.GetTelegramBot().Send(msg)
}

func SendPicture(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	waitMsg := tgbotapi.NewMessage(chatID, service.GetImageWaitAnswer())
	waitMsg.ReplyToMessageID = update.Message.MessageID
	service.GetTelegramBot().Send(waitMsg)

	msg := tgbotapi.NewMessage(chatID, service.GetEdenaiImage(update))
	msg.ReplyToMessageID = update.Message.MessageID
	service.GetTelegramBot().Send(msg)
}

func isMessageRequestForEdenai(update *tgbotapi.Update) bool {
	if update.Message != nil && strings.Contains(update.Message.Text, "generat") {
		return true
	}
	return false
}

func isMessageRequestForGPT(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text != "/start"
}
