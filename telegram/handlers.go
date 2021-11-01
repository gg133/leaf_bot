package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Constants for command handling
const (
	commandStart = "start"
)

//Constants for message templates
const (
	replyStartTemplate     = "Hi! To save links to your Pocket account, first you need to give me access to it. To do this, follow the link:\n%s"
	replyAlreadyAuthorized = "You are already authorized! Send me link and I will save it."
)

//handleCommand method handle telegram commands (message, which starts from "/")
func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}

}

//handleMessage method handle messages (in this case, it sends back a message written by the user before.)
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}

//handleStartCommand method generates authentification link for pocket and sends message with it link
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

//handleUnnkownCommand method sends message which says the command is unknown
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "You entered unknown command")

	_, err := b.bot.Send(msg)
	return err
}
