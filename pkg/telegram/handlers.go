package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"net/url"
)

const (
	commandStart       = "start"
	replyStartTemplate = "Чтобы" +
		" сохранять свои ссылки в своем Pocket аккаунте, перейди по ссылке:\n%s"
	replyAlreadyAuthorized = "Вы уже зарегистрированы)"
)

// handleMessage - обработчик СООБЩЕНИЙ
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	// формируем ответ
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена")
	// лог сообщения (время, кто прислал, что прислал)
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)

	if err != nil {
		return errUnauthorized
	}

	if err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}
	_, err = b.bot.Send(msg)
	return err
}

// handleCommand - обработчик КОМАНД
func (b *Bot) handleCommand(command *tgbotapi.Message) error {
	// лог команды (время, кто прислал, какую команду)
	log.Printf("[%s] %s", command.From.UserName, command.Text)
	//проверка команд
	switch command.Command() {
	case commandStart:
		return b.handleStartCommand(command)
	default:
		return b.handleUnknownCommand(command)
	}
}

// обработчик КОМАНДЫ "/start"
func (b *Bot) handleStartCommand(cmd *tgbotapi.Message) error {
	_, err := b.getAccessToken(cmd.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(cmd)
	}
	msg := tgbotapi.NewMessage(cmd.Chat.ID, replyAlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

// обработчик НЕИЗВЕСТНОЙ КОМАНДЫ
func (b *Bot) handleUnknownCommand(cmd *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(cmd.Chat.ID, "Я не знаю такой команды")
	_, err := b.bot.Send(msg)
	return err
}
