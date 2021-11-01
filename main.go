package main

import (
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yalagtyarzh/leaf_bot/config"
	"github.com/yalagtyarzh/leaf_bot/pocket"
	"github.com/yalagtyarzh/leaf_bot/repository"
	"github.com/yalagtyarzh/leaf_bot/repository/boltdb"
	"github.com/yalagtyarzh/leaf_bot/server"
	"github.com/yalagtyarzh/leaf_bot/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI("YourAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("YourAwesomePocketToken")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")

	authServer := server.NewAuthServer(pocketClient, tokenRepository, "https://t.me/leafAZ_bot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
