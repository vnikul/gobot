package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gobot/entities"
	"gobot/service"
	"log"
	"math/rand"
	"strings"
	"time"
)

type RetryChannel chan tgbotapi.MessageConfig

type GrammBot struct {
	*tgbotapi.BotAPI
	updates        tgbotapi.UpdatesChannel
	retryChannel   RetryChannel
	welcomeSticker string
}

func NewGrammBot(config entities.Config) (*GrammBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return &GrammBot{}, err
	}
	bot.Debug = config.Debug

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	retryChan := make(chan tgbotapi.MessageConfig)

	return &GrammBot{BotAPI: bot, welcomeSticker: config.WelcomeSticker, updates: updates, retryChannel: retryChan}, nil
}

func (bot *GrammBot) ProcessUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if update.Message.NewChatMembers != nil {
		sticker := tgbotapi.FileID(bot.welcomeSticker)

		msg := tgbotapi.NewSticker(update.Message.Chat.ID, sticker)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.BotAPI.Send(msg); err != nil {
			if strings.Contains(err.Error(), "replied message not found") {
				return
			}
			log.Printf("There was an error %s\n", err.Error())
			return
		}
	}
	updateMsg := ""
	if update.Message.Text != "" {
		updateMsg = update.Message.Text
	} else if update.Message.Caption != "" {
		updateMsg = update.Message.Caption
	}
	if updateMsg != "" && (strings.Contains(strings.ToLower(updateMsg), "pipa") || service.ContainsPipa(updateMsg)) {
		pipa := tgbotapi.NewMessage(update.Message.Chat.ID, "PIPA")
		pipa.ReplyToMessageID = update.Message.MessageID
		_, err := bot.Send(pipa)
		if err != nil {
			bot.processError(pipa, err)
		}
	}

	if updateMsg != "" && service.ContainsHeh(updateMsg) {
		benis := tgbotapi.NewMessage(update.Message.Chat.ID, "BENIS")
		benis.ReplyToMessageID = update.Message.MessageID
		_, err := bot.BotAPI.Send(benis)
		if err != nil {
			bot.processError(benis, err)
		}
	}

	if updateMsg != "" && service.ContainsYes(updateMsg) && rand.Intn(2) == 1 {
		badWord := tgbotapi.NewMessage(update.Message.Chat.ID, "ПИЗДА")
		badWord.ReplyToMessageID = update.Message.MessageID
		_, err := bot.BotAPI.Send(badWord)
		if err != nil {
			bot.processError(badWord, err)
		}
	}

	if updateMsg != "" && service.ContainsNo(updateMsg) && rand.Intn(2) == 1 {
		badWord := tgbotapi.NewMessage(update.Message.Chat.ID, "ГОМОСЕКСУАЛЬНЫЙ ОТВЕТ")
		badWord.ReplyToMessageID = update.Message.MessageID
		_, err := bot.BotAPI.Send(badWord)
		if err != nil {
			bot.processError(badWord, err)
		}
	}

	return
}

func (bot *GrammBot) ProcessMessage(msg tgbotapi.MessageConfig) {
	time.Sleep(10 * time.Second)
	_, err := bot.BotAPI.Send(msg)
	if err != nil {
		bot.processError(msg, err)
	}
}

func (bot *GrammBot) processError(msg tgbotapi.MessageConfig, err error) {
	if strings.Contains(err.Error(), "chat not found") {
		return
	}
	if strings.Contains(err.Error(), "too many requests") {
		bot.retryChannel <- msg
	}
	log.Fatalf("There was an error %s", err.Error())
	return

}
