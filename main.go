package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yalagtyarzh/leaf_bot/pocket"
	"github.com/yalagtyarzh/leaf_bot/telegram"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("YourAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("YourAwesomePocketToken")
	if err != nil {
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, pocketClient)
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
