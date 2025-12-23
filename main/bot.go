package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"src/config"
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

var monthNames = map[string]string{
	"01": "January",
	"02": "February",
	"03": "March",
	"04": "April",
	"05": "May",
	"06": "June",
	"07": "July",
	"08": "August",
	"09": "September",
	"10": "October",
	"11": "November",
	"12": "December",
}

func daysInMonth(month string, year int) int {
	monthDays := map[string]int{
		"January":   31,
		"February":  28,
		"March":     31,
		"April":     30,
		"May":       31,
		"June":      30,
		"July":      31,
		"August":    31,
		"September": 30,
		"October":   31,
		"November":  30,
		"December":  31,
	}

	if month == "February" && isLeapYear(year) {
		return 29
	}

	return monthDays[month]
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func sendMessage(bot *tg_bot.BotAPI, msg tg_bot.MessageConfig) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Bot failed to send message: %v", err)
	}
}

func handleNewPoop(db *sql.DB, userId int64, username string, msgId int, timestamp int64) {
	t := time.Unix(timestamp, 0).UTC()
	sqliteTimestamp := t.Format("2006-01-02 15:04:05")

	err := repo.LogPoop(db, userId, username, msgId, sqliteTimestamp, t.Unix())
	if err != nil {
		log.Printf("Failed to log poop: %v", err)
		return
	}
	log.Println("Poop logged successfully!")
}

func handleCommands(bot *tg_bot.BotAPI, db *sql.DB, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) {
	log.Println("Command received:", update.Message.Command())
	switch update.Message.Command() {
	case "my_poop_log":
		globalPoopCount, errGlobal := repo.GetGlobalPoopCount(db, userId)
		monthlyPoopCounts, errMonthly := repo.GetMonthlyPoopStats(db, userId)
		daysWithoutPoop, errNoPoop := repo.GetDaysWithoutPoop(db, userId)
		maxStreak, errStreak := repo.GetMaxPoopStreak(db, userId)
		day, poops, mostPoopsErr := repo.GetDayWithMostPoops(db, userId)

		if errGlobal != nil || errMonthly != nil || errNoPoop != nil || errStreak != nil || mostPoopsErr != nil {
			msg.Text = "Sorry, I couldn't retrieve your poop log\\. Please try again later\\!"
			sendMessage(bot, msg)
		}

		year := time.Now().Year()
		monthlyAverages := make(map[string]float64)
		for _, mpc := range monthlyPoopCounts {
			yearMonth := mpc.Month
			month := yearMonth[5:]
			daysInMonth := daysInMonth(monthNames[month], year)

			if time.Now().Month().String() == monthNames[month] {
				daysInMonth = time.Now().Day()
			}

			monthlyAverages[month] = float64(mpc.PoopCount) / float64(daysInMonth)
		}

		yearlyAverage := float64(globalPoopCount) / float64(time.Now().YearDay())

		msg.Text = fmt.Sprintf("*💩 Poop Report for @%s 💩*\n\n", strings.ReplaceAll(update.Message.From.UserName, "_", "\\_"))
		msg.Text += "*📅 Yearly Overview:*\n"
		msg.Text += fmt.Sprintf("🟤 Total dumps: `%d`\n", globalPoopCount)
		msg.Text += fmt.Sprintf("📊 Average per day: `%.2f`\n", yearlyAverage)
		msg.Text += fmt.Sprintf("🚫 Days without poops: `%d`\n", daysWithoutPoop)
		msg.Text += fmt.Sprintf("🔥 Max poop streak: `%d`\n", maxStreak)
		msg.Text += fmt.Sprintf("💣 Day with most poops: `%s with %d poops`\n\n", day, poops)
		msg.Text += "*📅 Monthly Breakdown:\n*"

		index := 0
		for _, mpc := range monthlyPoopCounts {
			yearMonth := mpc.Month
			month := yearMonth[5:]
			msg.Text += fmt.Sprintf("🗓 %s:  `%d poops`   \\(📊 Avg:   `%.2f per day`\\)\n", monthNames[month], mpc.PoopCount, monthlyAverages[month])
			index++
		}

		sendMessage(bot, msg)
	case "leaderboard":
		monthlyLeaderboard, err := repo.GetMonthlyLeaderboard(db)
		if err != nil {
			msg.Text = "Sorry, I couldn't retrieve the monthly leaderboard\\. Please try again later\\!"
			sendMessage(bot, msg)
		}
		msg.Text = "This month's leaderboard:\n"
		for _, user := range monthlyLeaderboard {
			escapedUsername := strings.ReplaceAll(user.Username, "_", "\\_")
			escapedUsername = strings.ReplaceAll(escapedUsername, "-", "\\-")
			msg.Text += fmt.Sprintf("\t\t\t• %s \\- %d💩\n", escapedUsername, user.PoopCount)
		}
		sendMessage(bot, msg)
	case "bottom_poopers":
		bottomPoopers, err := repo.GetBottomPoopers(db)
		if err != nil {
			msg.Text = "Sorry, I couldn't retrieve the bottom poopers\\. Please try again later\\!"
			sendMessage(bot, msg)
		}
		msg.Text = "This month's bottom poopers are:\n" + buildPoodiumMessage(bottomPoopers)
		sendMessage(bot, msg)
	case "poodium":
		monthlyPoodium, err := repo.GetMonthlyPoodium(db)
		if err != nil {
			msg.Text = "Sorry, I couldn't retrieve the monthly poodium\\. Please try again later\\!"
			sendMessage(bot, msg)
		}
		msg.Text = "This month's top poopers are:\n" + buildPoodiumMessage(monthlyPoodium)
		sendMessage(bot, msg)
	case "poodium_year":
		yearlyPoodium, err := repo.GetYearlyPoodium(db)
		if err != nil {
			msg.Text = "Sorry, I couldn't retrieve the yearly poodium\\. Please try again later\\!"
			sendMessage(bot, msg)
		}
		msg.Text = "This year's top poopers are:\n" + buildPoodiumMessage(yearlyPoodium)
		sendMessage(bot, msg)
	default:
		message := ""
		if update.Message.Command() != "help" {
			message += "Sorry, I don't recognize that command. "
		}
		msg.Text = message + "Here are the commands I understand:\n" +
			"\t\t\t\t• _/help_ \\- Get a list of available commands\n" +
			"\t\t\t\t• _/my\\_poop\\_log_ \\- Get your personal monthly poop statistics\n" +
			"\t\t\t\t• _/leaderboard_ \\- Get the monthly leaderboard\n" +
			"\t\t\t\t• _/bottom\\_poopers_ \\- Get the reverse poodium\n" +
			"\t\t\t\t• _/poodium_ \\- Get the monthly poodium\n" +
			"\t\t\t\t• _/poodium\\_year_ \\- Get the yearly poodium"
		sendMessage(bot, msg)
	}
}

func buildPoodiumMessage(topPoopers []repo.UserPoopCount) string {
	escape := func(s string) string {
		s = strings.ReplaceAll(s, "_", "\\_")
		s = strings.ReplaceAll(s, "-", "\\-")
		return s
	}
	return "🥇 " + escape(fmt.Sprint(topPoopers[0].Username)) + " \\- " + fmt.Sprint(topPoopers[0].PoopCount) + "💩\n" +
		"🥈 " + escape(fmt.Sprint(topPoopers[1].Username)) + " \\- " + fmt.Sprint(topPoopers[1].PoopCount) + "💩\n" +
		"🥉 " + escape(fmt.Sprint(topPoopers[2].Username)) + " \\- " + fmt.Sprint(topPoopers[2].PoopCount) + "💩"
}

func sendMonthlyPoodium(bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
	topPoopers, err := repo.GetPastMonthPoodium(db)
	if err != nil {
		log.Printf("Failed to get top poopers for monthly poodium: %v", err)
		return
	}

	now := time.Now()
	pastMonth := now.AddDate(0, -1, 0)
	_, month, _ := pastMonth.Date()
	monthStr := fmt.Sprintf("%02d", month)

	messageText := "🏆 Poodium for " + monthNames[monthStr] + " 🏆\n" + buildPoodiumMessage(topPoopers)
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

func sendYearlyPoodium(bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
	topPoopers, err := repo.GetYearlyPoodium(db)
	if err != nil {
		log.Printf("Failed to get top poopers for yearly poodium: %v", err)
		return
	}

	year, _, _ := time.Now().Date()

	messageText := "🏆 Poodium for " + fmt.Sprint(year) + " 🏆\n" + buildPoodiumMessage(topPoopers)
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

	db := repo.OpenDBConnection(cfg)

	// Schedule the monthly poodium message
	monthlyCron := cron.New()
	_, err = monthlyCron.AddFunc("0 0 1 * *", func() {
		sendMonthlyPoodium(bot, db, cfg.GroupChatID)
	})
	if err != nil {
		log.Fatalf("Failed to schedule monthly poodium message: %v", err)
	}
	monthlyCron.Start()

	// Schedule the yearly poodium message
	yearlyCron := cron.New()
	_, err = yearlyCron.AddFunc("0 0 1 1 *", func() {
		sendYearlyPoodium(bot, db, cfg.GroupChatID)
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
				handleNewPoop(db, userID, username, messageID, int64(update.Message.Date))
				handleReactions(cfg, chatID, messageID, update.Message.Sticker)
			}

			if update.Message.Command() != "" {
				handleCommands(bot, db, update, userID, msg)
			}
		case cfg.MyChatID:
			if update.Message.Text == "💩" || (update.Message.Sticker != nil && update.Message.Sticker.Emoji == "💩") {
				userID = update.Message.ForwardFrom.ID
				username = update.Message.ForwardFrom.UserName
				messageID = -update.Message.MessageID
				handleNewPoop(db, userID, username, messageID, int64(update.Message.ForwardDate))
				handleReactions(cfg, chatID, update.Message.MessageID, update.Message.Sticker)
			}

			if update.Message.Command() != "" {
				handleCommands(bot, db, update, userID, msg)
			}
		default:
			if update.Message.Command() != "" {
				msg.Text = "Sorry, I only respond to commands in the group chat\\."
				sendMessage(bot, msg)
			}
		}
	}

}
