package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yalagtyarzh/leaf_bot/config"
	"github.com/yalagtyarzh/leaf_bot/pocket"
	"github.com/yalagtyarzh/leaf_bot/repository"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectUrl     string
	messages        config.Messages
}

//NewBotcreates a new bot object
func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repository.TokenRepository, redirectUrl string, messages config.Messages) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, tokenRepository: tr, redirectUrl: redirectUrl, messages: messages}
}

//Start method starts the bot in telegram, Initializes the telegram bot update channel and starts handle updates
func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)
	return nil
}

//handleUpdates method is an infinity cycle, which handle updates from telegram chats
func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}

			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

//initUpdatesChannel method creates telegram bot update channel
func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
