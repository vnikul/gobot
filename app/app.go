package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gobot/infrastructure"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Run() {
	_, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGKILL)

	conf, _ := infrastructure.LoadConfigFromEnv()

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Fatalf("There was an error %s", err.Error())
		return
	}

	bot.Debug = conf.Debug

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	startMessage := tgbotapi.NewMessage(conf.ChatId, "status ok")
	_, err = bot.Send(startMessage)
	if err != nil {
		log.Fatalf("There was an error %s", err.Error())
		return
	}

	for {
		select {
		case <-exit:
			sleepMsg := tgbotapi.NewMessage(conf.ChatId, "i sleep")
			_, err = bot.Send(sleepMsg)
			if err != nil {
				log.Fatalf("There was an error %s", err.Error())
				return
			}
			cancel()
		case update := <-updates:
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
				if update.Message.Text != "" && (strings.Contains(strings.ToLower(update.Message.Text), "pipa") || strings.Contains(strings.ToLower(update.Message.Text), "пипа")) {
					pipa := tgbotapi.NewMessage(update.Message.Chat.ID, "PIPA")
					pipa.ReplyToMessageID = update.Message.MessageID
					_, err = bot.Send(pipa)
					if err != nil {
						log.Fatalf("There was an error %s", err.Error())
						return
					}
				}

				if update.Message.Text != "" && (strings.Contains(strings.ToLower(strings.ReplaceAll(update.Message.Text, " ", "")), "haha") || strings.Contains(strings.ToLower(strings.ReplaceAll(update.Message.Text, " ", "")), "хаха")) {
					benis := tgbotapi.NewMessage(update.Message.Chat.ID, "BENIS")
					benis.ReplyToMessageID = update.Message.MessageID
					_, err = bot.Send(benis)
					if err != nil {
						log.Fatalf("There was an error %s", err.Error())
						return
					}
				}
			}(update)
		}
	}
}
