package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"

	"src/config"
	"src/formatters"
	"src/handlers"
	repo "src/repository"

	tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type AddReactionRequest struct {
	ChatID    int64          `json:"chat_id"`
	MessageID int            `json:"message_id"`
	Reaction  []ReactionType `json:"reaction"`
	IsBig     bool           `json:"is_big"`
}

type ReactionType struct {
	Type          string `json:"type"`
	Emoji         string `json:"emoji,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

func sendMessage(bot *tg_bot.BotAPI, msg tg_bot.MessageConfig) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Bot failed to send message: %v", err)
	}
}

func handleNewPoop(ctx context.Context, db *sql.DB, userId int64, username string, msgId int, timestamp int64) {
	t := time.Unix(timestamp, 0).UTC()
	sqliteTimestamp := t.Format("2006-01-02 15:04:05")

	err := repo.LogPoop(ctx, db, userId, username, msgId, sqliteTimestamp, t.Unix())
	if err != nil {
		log.Printf("Failed to log poop: %v", err)
		return
	}
	log.Println("Poop logged successfully!")
}

func sendMonthlyPoodium(ctx context.Context, bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
	topPoopers, err := repo.GetPastMonthPoodium(ctx, db)
	if err != nil {
		log.Printf("Failed to get top poopers for monthly poodium: %v", err)
		return
	}

	now := time.Now()
	pastMonth := now.AddDate(0, -1, 0)
	_, month, _ := pastMonth.Date()
	monthStr := fmt.Sprintf("%02d", month)
	monthName := formatters.GetMonthName(monthStr)

	messageText := formatters.FormatPoodiumTitle(monthName) + formatters.BuildPoodiumMessage(topPoopers)
	msg := tg_bot.NewMessage(chatID, messageText)
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send monthly poodium message: %v", err)
		return
	}

	_, err = bot.Send(tg_bot.PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           sentMsg.MessageID,
		DisableNotification: true,
	})
	if err != nil {
		log.Printf("Failed to pin monthly poodium message: %v", err)
	}
}

func sendYearlyPoodium(ctx context.Context, bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
	topPoopers, err := repo.GetYearlyPoodium(ctx, db)
	if err != nil {
		log.Printf("Failed to get top poopers for yearly poodium: %v", err)
		return
	}

	year, _, _ := time.Now().Date()
	messageText := formatters.FormatYearlyPoodiumTitle(year) + formatters.BuildPoodiumMessage(topPoopers)
	msg := tg_bot.NewMessage(chatID, messageText)
	sendMessage(bot, msg)
}

func handleReactions(cfg *config.Config, chatID int64, messageID int, sticker *tg_bot.Sticker) {
	var reactEmoji = "💩"

	if sticker != nil && sticker.Emoji == "💩" {
		if sticker.SetName == "Poopers2" {
			switch sticker.FileUniqueID {
			case cfg.StickerIDs["struggle"]:
				reactEmoji = "😢"
			case cfg.StickerIDs["esnoopi"]:
				reactEmoji = "🎉"
			case cfg.StickerIDs["jurassic"]:
				reactEmoji = "🏆"
			case cfg.StickerIDs["girly"]:
				reactEmoji = "💅"
			case cfg.StickerIDs["scared"]:
				reactEmoji = "🫡"
			case cfg.StickerIDs["sus"]:
				reactEmoji = "🤨"
			}
		}
	}

	err := addReaction(cfg, chatID, messageID, reactEmoji)
	if err != nil {
		log.Printf("Failed to add reaction: %v", err)
	} else {
		log.Println("Reaction added successfully.")
	}
}

func addReaction(cfg *config.Config, chatID int64, messageID int, emoji string) error {
	url := fmt.Sprintf("%s%s/setMessageReaction", cfg.APIBaseURL, cfg.TelegramToken)

	reactionRequest := AddReactionRequest{
		ChatID:    chatID,
		MessageID: messageID,
		Reaction:  []ReactionType{{Type: "emoji", Emoji: emoji}},
		IsBig:     false,
	}

	jsonBody, err := json.Marshal(reactionRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal reaction request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add reaction, status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	fmt.Print("Starting bot...\n")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	bot, err := tg_bot.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create new bot instance: %v", err)
	}

	bot.Debug = true

	updateConfig := tg_bot.NewUpdate(0)
	updateConfig.Timeout = 30
	updateConfig.AllowedUpdates = []string{"message", "message_reaction"}
	updates := bot.GetUpdatesChan(updateConfig)

	db, err := repo.OpenDBConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx := context.Background()

	// Schedule the monthly poodium message
	monthlyCron := cron.New()
	_, err = monthlyCron.AddFunc("0 0 1 * *", func() {
		sendMonthlyPoodium(ctx, bot, db, cfg.GroupChatID)
	})
	if err != nil {
		log.Fatalf("Failed to schedule monthly poodium message: %v", err)
	}
	monthlyCron.Start()

	// Schedule the yearly poodium message
	yearlyCron := cron.New()
	_, err = yearlyCron.AddFunc("0 0 1 1 *", func() {
		sendYearlyPoodium(ctx, bot, db, cfg.GroupChatID)
	})
	if err != nil {
		log.Fatalf("Failed to schedule yearly poodium message: %v", err)
	}
	yearlyCron.Start()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tg_bot.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ParseMode = tg_bot.ModeMarkdownV2
		userID := update.Message.From.ID
		username := update.Message.From.UserName
		messageID := update.Message.MessageID
		chatID := update.Message.Chat.ID

		switch chatID {
		case cfg.GroupChatID:
			if update.Message.Text == "💩" || (update.Message.Sticker != nil && update.Message.Sticker.Emoji == "💩") {
				log.Println("New poop detected!")
				handleNewPoop(ctx, db, userID, username, messageID, int64(update.Message.Date))
				handleReactions(cfg, chatID, messageID, update.Message.Sticker)
			}

			if update.Message.Command() != "" {
				handlers.HandleCommand(ctx, bot, db, update, userID, msg)
			}
		case cfg.MyChatID:
			if update.Message.Text == "💩" || (update.Message.Sticker != nil && update.Message.Sticker.Emoji == "💩") {
				userID = update.Message.ForwardFrom.ID
				username = update.Message.ForwardFrom.UserName
				messageID = -update.Message.MessageID
				handleNewPoop(ctx, db, userID, username, messageID, int64(update.Message.ForwardDate))
				handleReactions(cfg, chatID, update.Message.MessageID, update.Message.Sticker)
			}

			if update.Message.Command() != "" {
				handlers.HandleCommand(ctx, bot, db, update, userID, msg)
			}
		default:
			if update.Message.Command() != "" {
				msg.Text = "Sorry, I only respond to commands in the group chat\\."
				sendMessage(bot, msg)
			}
		}
	}

}
