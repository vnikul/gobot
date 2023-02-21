package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gobot/infrastructure"
	"log"
)

func Run() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, _ := infrastructure.LoadConfigFromEnv()

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Fatalf("There was an error %s", err.Error())
		return
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	startMessage := tgbotapi.NewMessage(conf.ChatId, "status ok")
	_, err = bot.Send(startMessage)
	if err != nil {
		log.Fatalf("There was an error %s", err.Error())
		return
	}

	go func(cancel context.CancelFunc) {
		cancel()
		return
	}(cancel)

	for update := range updates {
		go func(update tgbotapi.Update) {
			if update.Message == nil {
				return
			}
			if update.Message.NewChatMembers != nil {

				sticker := tgbotapi.FileID(conf.WelcomeSticker)

				msg := tgbotapi.NewSticker(update.Message.Chat.ID, sticker)
				msg.ReplyToMessageID = update.Message.MessageID

				if _, err := bot.Send(msg); err != nil {
					log.Fatalf("There was an error %s", err.Error())
					return
				}
			}
		}(update)
	}
}
