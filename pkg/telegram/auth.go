package telegram

import (
	"First-TgBot-On-GO/pkg/repository"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) initAuthorizationProcess(cmd *tgbotapi.Message) error {
	// генерируется ссылка в виде %s?chat_id=%d"
	authLink, err := b.generateAuthLink(cmd.Chat.ID)
	if err != nil {
		return err
	}
	// формируется ответное сообщение, с объединением редирект ссылки и темплейта
	// как формируется редирект ссылка в auth.go
	msg := tgbotapi.NewMessage(cmd.Chat.ID, fmt.Sprintf(replyStartTemplate, authLink))
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessToken)
}

// generateAuthLink - создание редирект ссылки
func (b *Bot) generateAuthLink(chatID int64) (string, error) {
	// генерируем ссылку в формает http://localhost/?chat_id=%d
	redirectURL := b.generateRedirectLink(chatID)
	// генерируем токен запроса в покет клиент
	// данный запрос используется для предоставления клиентом аккаунта в тг-боте
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}
	// в базу данных сохраняется чат айди и запрос токена под бакетом RequestToken
	if err = b.tokenRepository.Save(chatID, requestToken, repository.RequestToken); err != nil {
		return "", err
	}
	// получаем ссылку авторизации в покет клиенте
	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectLink(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
