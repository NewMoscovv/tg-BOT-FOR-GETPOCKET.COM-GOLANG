package telegram

import (
	"First-TgBot-On-GO/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

// Bot структура бота, содержит объект api тг, api покет клиента, интерфейсы бд, ссылку
type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

// NewBot метод сборки бота
func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repository.TokenRepository, redirectURL string) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, redirectURL: redirectURL, tokenRepository: tr}
}

// Start запуск бота
func (b *Bot) Start() error {
	// лог создания бота
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)
	return nil
}

// initUpdatesChannel - создание инициализатора запросов на сервер
func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	// создается запрос
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// ждем получение ответов
	return b.bot.GetUpdatesChan(u)
}

// handleUpdates - получение обновлений с сервера
func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	// проверяются обновления
	for update := range updates {
		// если обновление не содержит сообщения - пропуск
		if update.Message == nil {
			continue
		}
		// если команда, то отправляем на обработчик команд
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		// если сообщение, то отправляем на обработчик сообщений
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}
