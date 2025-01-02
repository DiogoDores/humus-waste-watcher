package main

import (
	"utils"
	"os"
	"log"
    tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func get_api_key() string {
	utils.LoadEnv()
	return os.Getenv("TELEGRAM_TOKEN")
}

func main() {
	apiKey := get_api_key()
    if apiKey == "" {
        log.Fatal("TELEGRAM_TOKEN is not set")
    }

	bot, err := tg_bot.NewBotAPI(apiKey)
	utils.Check(err)

	bot.Debug = true

    updateConfig := tg_bot.NewUpdate(0)
    updateConfig.Timeout = 30
    updates := bot.GetUpdatesChan(updateConfig)

    for update := range updates {
        // Telegram can send many types of updates depending on what your Bot
        // is up to. We only want to look at messages for now, so we can
        // discard any other updates.
        if update.Message == nil {
            continue
        }

        msg := tg_bot.NewMessage(update.Message.Chat.ID, update.Message.Text)
        //msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.Text == "💩" {
			msg.Text = "hehe poopers"
			_, err := bot.Send(msg)
			utils.Check(err)
		} else {
			msg.Text = "NOT A POOP!"
			_, err := bot.Send(msg)
			utils.Check(err)
		}
    }
}