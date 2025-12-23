package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	dotenv "github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DBPath        string
	GroupChatID   int64
	MyChatID      int64

	StickerIDs map[string]string
	APIBaseURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	_ = dotenv.Load(".env")

	cfg := &Config{
		StickerIDs: make(map[string]string),
	}

	cfg.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_TOKEN is not set")
	}

	cfg.DBPath = os.Getenv("DB_PATH")
	if cfg.DBPath == "" {
		return nil, fmt.Errorf("DB_PATH is not set")
	}

	groupChatIDStr := os.Getenv("GROUP_CHAT_ID")
	if groupChatIDStr != "" {
		groupChatID, err := strconv.ParseInt(groupChatIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid GROUP_CHAT_ID: %w", err)
		}
		cfg.GroupChatID = groupChatID
	} else {
		cfg.GroupChatID = -1002481034087
		log.Println("Warning: GROUP_CHAT_ID not set, using default value")
	}

	myChatIDStr := os.Getenv("MY_CHAT_ID")
	if myChatIDStr != "" {
		myChatID, err := strconv.ParseInt(myChatIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid MY_CHAT_ID: %w", err)
		}
		cfg.MyChatID = myChatID
	} else {
		cfg.MyChatID = 824870685
		log.Println("Warning: MY_CHAT_ID not set, using default value")
	}

	cfg.StickerIDs = map[string]string{
		"struggle": "AgADOxkAAgTYWVE",
		"esnoopi":  "AgADRhoAAhq7WVE",
		"jurassic": "AgADQRcAAu99WVE",
		"girly":    "AgADrxkAAtHnYFE",
		"scared":   "AgADfBkAAgwIYVE",
		"sus":      "AgADcxgAAvVG0FE",
	}

	cfg.APIBaseURL = "https://api.telegram.org/bot"

	return cfg, nil
}
