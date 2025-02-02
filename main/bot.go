package main

import (
	"os"
	"log"
	"fmt"
	_ "modernc.org/sqlite"
	"database/sql"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"strings"

	"src/utils"
	repo "src/repository"

    tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

const telegramAPIBase = "https://api.telegram.org/bot"

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

var CHAT_ID = int64(-1002481034087)

func get_api_key() string {
	if os.Getenv("TELEGRAM_TOKEN") == "" {
		utils.LoadEnv()
	}
	return os.Getenv("TELEGRAM_TOKEN")
}

func send_message(bot *tg_bot.BotAPI, msg tg_bot.MessageConfig) {
	_, err := bot.Send(msg)
	utils.CheckError("Bot failed to send message", err)
}

func handle_new_poop(db *sql.DB, userId int64, username string, msgId int, timestamp int64) {
	t := time.Unix(timestamp, 0).UTC() 
	sqliteTimestamp := t.Format("2006-01-02 15:04:05")

	err := repo.Log_Poop(db, userId, username, msgId, sqliteTimestamp)
	utils.CheckError("Failed to log poop", err)
	log.Println("Poop logged successfully!")
}

func handle_commands(bot *tg_bot.BotAPI, db *sql.DB, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) {
	log.Println("Command received:", update.Message.Command())
	switch update.Message.Command() {
		case "my_poop_log":
			globalPoopCount, errGlobal := repo.Get_Global_Poop_Count(db, userId)
			monthlyPoopCounts, errMonthly := repo.Get_Monthly_Poop_Stats(db, userId)
			daysWithoutPoop, errNoPoop := repo.Get_Days_Without_Poop(db, userId)
			maxStreak, errStreak := repo.Get_Max_Poop_Streak(db, userId)
			day, poops, mostPoopsErr := repo.Get_Day_With_Most_Poops(db, userId)

			if errGlobal != nil || errMonthly != nil || errNoPoop != nil || errStreak != nil || mostPoopsErr != nil {
				msg.Text = "Sorry, I couldn't retrieve your poop log\\. Please try again later\\!"
				send_message(bot, msg)
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

			msg.Text = fmt.Sprintf("*💩 Poop Report for @%s 💩*\n\n", 	strings.ReplaceAll(update.Message.From.UserName, "_", "\\_"))
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

			send_message(bot, msg)
		case "leaderboard":
			monthlyLeaderboard, err := repo.Get_Monthly_Leaderboard(db)
			if err != nil {
				msg.Text = "Sorry, I couldn't retrieve the monthly leaderboard\\. Please try again later\\!"
				send_message(bot, msg)
			}
			msg.Text = "This month's leaderboard:\n"
			for _, user := range monthlyLeaderboard {
				escapedUsername := strings.ReplaceAll(user.Username, "_", "\\_")
				msg.Text += fmt.Sprintf("\t\t\t• %s \\- %d💩\n", escapedUsername, user.PoopCount)
			}
			send_message(bot, msg)
		case "bottom_poopers":
			bottomPoopers, err := repo.Get_Bottom_Poopers(db)
			if err != nil {
				msg.Text = "Sorry, I couldn't retrieve the bottom poopers\\. Please try again later\\!"
				send_message(bot, msg)
			}
			msg.Text = "This month's bottom poopers are:\n" + build_poodium_message(bottomPoopers)
			send_message(bot, msg)
		case "poodium":
			monthlyPoodium, err := repo.Get_Monthly_Poodium(db)
			if err != nil {
				msg.Text = "Sorry, I couldn't retrieve the monthly poodium\\. Please try again later\\!"
				send_message(bot, msg)
			}
			msg.Text = "This month's top poopers are:\n" + build_poodium_message(monthlyPoodium)
			send_message(bot, msg)
		case "poodium_year":
			yearlyPoodium, err := repo.Get_Yearly_Poodium(db)
			if err != nil {
				msg.Text = "Sorry, I couldn't retrieve the yearly poodium\\. Please try again later\\!"
				send_message(bot, msg)
			}
			msg.Text = "This year's top poopers are:\n" + build_poodium_message(yearlyPoodium)
			send_message(bot, msg)
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
			send_message(bot, msg)
	}
}

func build_poodium_message(topPoopers []repo.UserPoopCount) string {
	return "🥇 " + fmt.Sprint(strings.ReplaceAll(topPoopers[0].Username, "_", "\\_")) + " \\- " + fmt.Sprint(topPoopers[0].PoopCount) + "💩\n" +
			"🥈 " + fmt.Sprint(strings.ReplaceAll(topPoopers[1].Username, "_", "\\_")) + " \\- " + fmt.Sprint(topPoopers[1].PoopCount) + "💩\n" +
			"🥉 " + fmt.Sprint(strings.ReplaceAll(topPoopers[2].Username, "_", "\\_")) + " \\- " + fmt.Sprint(topPoopers[2].PoopCount) + "💩"
}

func send_monthly_poodium(bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
    topPoopers, err := repo.Get_Past_Month_Poodium(db)
	utils.CheckError("Failed to get top poopers", err)

    now := time.Now()
    pastMonth := now.AddDate(0, -1, 0)
    _, month, _ := pastMonth.Date()
    monthStr := fmt.Sprintf("%02d", month)

    messageText := "🏆 Poodium for " + monthNames[monthStr] + " 🏆\n" + build_poodium_message(topPoopers)
    msg := tg_bot.NewMessage(chatID, messageText)
    sentMsg, err := bot.Send(msg)
    utils.CheckError("Failed to send message", err)

    bot.Send(tg_bot.PinChatMessageConfig{
        ChatID:              chatID,
        MessageID:           sentMsg.MessageID,
        DisableNotification: true,
    })
}

func send_yearly_poodium(bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
    topPoopers, err := repo.Get_Yearly_Poodium(db)
	utils.CheckError("Failed to get top poopers", err)

	year, _, _ := time.Now().Date()

	messageText := "🏆 Poodium for " + fmt.Sprint(year) + " 🏆\n" + build_poodium_message(topPoopers)
    msg := tg_bot.NewMessage(chatID, messageText)
    send_message(bot, msg)
}

func addReaction(botToken string, chatID int64, messageID int, emoji string) error {
	url := fmt.Sprintf("%s%s/setMessageReaction", telegramAPIBase, botToken)

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
	apiKey := get_api_key()
    if apiKey == "" {
        log.Fatal("TELEGRAM_TOKEN is not set")
    }

	bot, err := tg_bot.NewBotAPI(apiKey)
	utils.CheckError("Failed to create new bot instance", err)

	bot.Debug = true

    updateConfig := tg_bot.NewUpdate(0)
    updateConfig.Timeout = 30
	updateConfig.AllowedUpdates = []string{"message", "message_reaction"}
    updates := bot.GetUpdatesChan(updateConfig)

	db := repo.Open_DB_Connection()

	// Schedule the poodium message
    monthlyCron := cron.New()
    _, err = monthlyCron.AddFunc("0 0 1 * *", func() {
        chatID := int64(CHAT_ID)
        send_monthly_poodium(bot, db, chatID)
    })
	utils.CheckError("Failed to schedule monthly poodium message", err)
	monthlyCron.Start()

	// Schedule the poodium message
    yearlyCron := cron.New()
	_, err = yearlyCron.AddFunc("0 0 1 1 *", func() {
		chatID := int64(CHAT_ID)
		send_yearly_poodium(bot, db, chatID)
	})
	utils.CheckError("Failed to schedule yearly poodium message", err)
	yearlyCron.Start()

	for update := range updates {
		monthlyCron.Start()
		yearlyCron.Start()
		if update.Message == nil {
			continue
		}

		msg := tg_bot.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ParseMode = tg_bot.ModeMarkdownV2
		userId := update.Message.From.ID
		username := update.Message.From.UserName
		messageId := update.Message.MessageID

		if update.Message.Text == "💩" || (update.Message.Sticker != nil && update.Message.Sticker.Emoji == "💩") {
			log.Println("New poop detected!")
			handle_new_poop(db, userId, username, messageId, int64(update.Message.Date))
			err := addReaction(apiKey, update.Message.Chat.ID, update.Message.MessageID, "💩")
			if err != nil {
				log.Printf("Failed to add reaction: %v", err)
			} else {
				log.Println("Reaction added successfully.")
			}
		}

		if update.Message.Command() != "" {
			handle_commands(bot, db, update, userId, msg)
		}
	}

}