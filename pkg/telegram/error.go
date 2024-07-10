package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURL   = errors.New("Invalid URL")
	errUnauthorized = errors.New("Unauthorized user")
	errUnableToSave = errors.New("Unable to save telegram")
)

func (b *Bot) handleError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, "")
	switch err {
	case errInvalidURL:
		msg.Text = "Это не валидная ссылка!"
	case errUnauthorized:
		msg.Text = "Ты не авторизирован!"
	case errUnableToSave:
		msg.Text = "Увы, не удалось сохранить ссылку"

	default:
		msg.Text = "Неизвестная ошибка"
	}
	b.bot.Send(msg)
}
