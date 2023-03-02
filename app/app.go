package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gobot/infrastructure"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	_, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(exit, os.Interrupt, syscall.SIGKILL)

	conf, _ := infrastructure.LoadConfigFromEnv()

	bot, err := NewGrammBot(conf)
	if err != nil {
		log.Fatalf("There was an error %s", err.Error())
		return
	}

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
		case update := <-bot.updates:
			go bot.ProcessUpdate(update)
		case update := <-bot.retryChannel:
			go bot.ProcessMessage(update)
		}
	}
}
