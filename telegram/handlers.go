package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const commandStart = "start"

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "You entered unknown command")

	switch message.Command() {
	case commandStart:
		msg.Text = "You entered /start command"
		_, err := b.bot.Send(msg)
		return err
	default:
		_, err := b.bot.Send(msg)
		return err
	}

}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}
