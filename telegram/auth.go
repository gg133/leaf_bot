package telegram

import (
	"context"
	"fmt"

	"github.com/yalagtyarzh/leaf_bot/repository"
)

//generateAuthLink creates redirectURL?chat_id template, gets request token with
//this template and with request token and redirect URL gets authentification URL
func (b *Bot) generateAuthLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectURL(chatID)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthURL(requestToken, redirectURL)
}

//generateRedirectURL generates a redirect URL templale with chat_id
func (b *Bot) generateRedirectURL(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatID)
}
