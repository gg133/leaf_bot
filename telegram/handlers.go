package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yalagtyarzh/leaf_bot/pocket"
)

//Constants for command handling
const (
	commandStart = "start"
)

//Constants for message templates
const (
	replyStartTemplate     = "Hi! To save links to your Pocket account, first you need to give me access to it. To do this, follow the link:\n%s"
	replyAlreadyAuthorized = "You are already authorized! Send me link and I will save it."
	replyAddSuccess        = "Link saved successfully!"
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

//handleMessage method checks for url validation in tg message, gets access token from db and if everything is url saves in Pocket
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvaildURL
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnauthorized
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, replyAddSuccess)
	_, err = b.bot.Send(msg)
	return err
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
