package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvaildURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "An unknown error has occurred")

	switch err {
	case errInvaildURL:
		msg.Text = "This is an invalid link!"
		b.bot.Send(msg)
	case errUnauthorized:
		msg.Text = "You are not logged in! Use the command /start"
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = "Oops, the link could not be saved. Please try again later."
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}

}
