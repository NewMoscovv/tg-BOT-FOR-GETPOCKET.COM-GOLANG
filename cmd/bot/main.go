package main

import (
	"First-TgBot-On-GO/configs/config"
	"First-TgBot-On-GO/pkg/repository"
	"First-TgBot-On-GO/pkg/repository/boltDB"
	"First-TgBot-On-GO/pkg/server"
	"First-TgBot-On-GO/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	pocketClient, err := pocket.NewClient("111620-3deeded429a5b2977c461da")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := boltDB.NewTokenRepository(db)

	tgBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository,
		"https://t.me/myb00tf0rt3st1ngbot")

	go func() {
		if err := tgBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}
func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestToken))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
