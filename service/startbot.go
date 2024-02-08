package service

import (
	"kaogpt/config"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot            *tgbotapi.BotAPI
	chatGPTService *ChatGPTService
	edenaiService  *EdenaiService
}

func NewTelegramService(chatGPTService *ChatGPTService, edenaiService *EdenaiService) *TelegramService {
	telegramBot, err := NewTelegramBotClient()
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}
	return &TelegramService{
		bot:            telegramBot,
		chatGPTService: chatGPTService,
		edenaiService:  edenaiService,
	}
}

func NewTelegramBotClient() (*tgbotapi.BotAPI, error) {
	tokens := config.ReadJsonTokens()
	return tgbotapi.NewBotAPI(tokens.Telegram)
}

func (t *TelegramService) Start() {
	updates := t.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			chatID := update.Message.Chat.ID
			if err := t.SendMessage("If you want to get an image, type \"generate\" firstly, or ask your question and AI will respond to you.", chatID); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
		if isMessageRequestForEdenai(&update) {
			if err := t.sendPicture(&update); err != nil {
				log.Printf("Error sending picture: %v", err)
			}
		} else if isMessageRequestForGPT(&update) {
			if err := t.sendAnswer(&update); err != nil {
				log.Printf("Error sending answer: %v", err)
			}
		}
	}
}

func (t *TelegramService) SendMessage(msg string, chatID int64) error {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	_, err := t.bot.Send(msgConfig)
	return err
}

func (t *TelegramService) sendAnswer(update *tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	if err := t.SendMessage(GetWaitAnswer(), chatID); err != nil {
		return err
	}
	answer, err := t.chatGPTService.GetGptAnswer(update.Message.Text)
	if err != nil {
		return err
	}
	if err := t.SendMessage(answer, chatID); err != nil {
		return err
	}
	return nil
}

func (t *TelegramService) sendPicture(update *tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	if err := t.SendMessage(GetImageWaitAnswer(), chatID); err != nil {
		return err
	}
	imageURL, err := t.edenaiService.GetEdenaiImage(update.Message.Text)
	if err != nil {
		return err
	}
	if err := t.SendMessage(imageURL, chatID); err != nil {
		return err
	}
	return nil
}

func isMessageRequestForEdenai(update *tgbotapi.Update) bool {
	return update.Message != nil && strings.Contains(update.Message.Text, "generate")
}

func isMessageRequestForGPT(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text != "/start"
}

func GetWaitAnswer() string {
	return "Wait, your answer is preparing."
}

func GetImageWaitAnswer() string {
	return "Wait, your image is generating."
}
