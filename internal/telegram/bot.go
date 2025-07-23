package telegram

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Api *tgbotapi.BotAPI
}

type TelegramBotInterface interface {
	SendAuthCode(chatID int64, code string) error
}

func (bot *Bot) SendAuthCode(chatID int64, code string) error {
	message := fmt.Sprintf("Your code: %s", code)
	return bot.sendMessage(chatID, message)
}

func BotCreate(token string) (*Bot, error) {
	Api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("Failed bot create")
		return nil, err
	}

	bot := &Bot{
		Api: Api,
	}

	log.Printf("Authorized on account %s", bot.Api.Self.UserName)
	return bot, nil
}

func (bot *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.Api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			if update.Message.IsCommand() && update.Message.Command() == "start" {
				bot.handleStart(update.Message)
			}
		}
	}
}

func (bot *Bot) handleStart(msg *tgbotapi.Message) {
	bot.sendMessage(msg.Chat.ID, "This bot will send you codes")
}

func (bot *Bot) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Api.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	return nil
}
