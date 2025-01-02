package main

import (
	"os"
	"log"
	"fmt"
	_ "modernc.org/sqlite"
	"database/sql"

	"src/utils"
	repo "src/repository"

    tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"

)

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

func get_api_key() string {
	utils.LoadEnv()
	return os.Getenv("TELEGRAM_TOKEN")
}

func send_message(bot *tg_bot.BotAPI, msg tg_bot.MessageConfig) {
	_, err := bot.Send(msg)
	utils.CheckError("Bot failed to send message", err)
}

func handle_new_poop(db *sql.DB, userId int64, username string, msgId int) {
	err := repo.Log_Poop(db, userId, username, msgId)
	utils.CheckError("Failed to log poop", err)
	log.Println("Poop logged successfully!")
}

func handle_commands(bot *tg_bot.BotAPI, db *sql.DB, update tg_bot.Update, userId int64, msg tg_bot.MessageConfig) {
	log.Println("Command received:", update.Message.Command())
	switch update.Message.Command() {
		case "poop":
			globalPoopCount, errGlobal := repo.Get_Global_Poop_Count(db, userId)
			monthlyPoopCount, errMonthly := repo.Get_Monthly_Poop_Count(db, userId)
			if errGlobal != nil || errMonthly != nil {
				msg.Text = "Sorry, I couldn't retrieve your poop count. Please try again later!"
				send_message(bot, msg)
			}
			monthlyPoopText := "times"
			if monthlyPoopCount == 1 {
				monthlyPoopText = "time"
			}
			globalPoopText := "poops"
			if globalPoopCount == 1 {
				globalPoopText = "poop"
			}
			msg.Text = "This month you've pooped " + fmt.Sprint(monthlyPoopCount) + " " + monthlyPoopText + ".\n" + 
						"You've logged " + fmt.Sprint(globalPoopCount) + " " + globalPoopText + " so far! 💩"
			send_message(bot, msg)
		case "poop_stats":
			monthlyPoopCounts, err := repo.Get_Monthly_Poop_Stats(db, userId)
			utils.CheckError("Failed to get monthly poop stats", err)

			msg.Text = ("Monthly Poop Stats:\n")
			for _, mpc := range monthlyPoopCounts {
				yearMonth := mpc.Month
				month := yearMonth[5:]
				msg.Text += fmt.Sprintf("\t\t\t• %s: %d💩\n", monthNames[month], mpc.PoopCount)
			}
			send_message(bot, msg)
		case "monthly_poodium":
			monthlyPoodium, err := repo.Get_Monthly_Poodium(db)
			if err != nil {
				msg.Text = "Sorry, I couldn't retrieve the monthly poodium. Please try again later!"
				send_message(bot, msg)
			}
			msg.Text = "This month's top poopers are:\n" + build_poodium_message(monthlyPoodium)
			send_message(bot, msg)
		case "yearly_poodium":
			yearlyPoodium, err := repo.Get_Yearly_Poodium(db)
			if err != nil {
				msg.Text = "Sorry, I couldn't retrieve the yearly poodium. Please try again later!"
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
						"\t\t\t\t• /help - Get a list of available commands\n" +
						"\t\t\t\t• /poop - Get your personal poop count\n" +
						"\t\t\t\t• /poop_stats - Get your personal monthly poop statistics\n" +
						"\t\t\t\t• /monthly_poodium - Get the monthly poodium\n" +
						"\t\t\t\t• /yearly_poodium - Get the yearly poodium"
			send_message(bot, msg)
	}
}

func build_poodium_message(topPoopers []repo.UserPoopCount) string {
	return "🥇 " + fmt.Sprint(topPoopers[0].Username) + " - " + fmt.Sprint(topPoopers[0].PoopCount) + "💩\n" +
			"🥈 " + fmt.Sprint(topPoopers[1].Username) + " - " + fmt.Sprint(topPoopers[1].PoopCount) + "💩\n" +
			"🥉 " + fmt.Sprint(topPoopers[2].Username) + " - " + fmt.Sprint(topPoopers[2].PoopCount) + "💩"
}

func send_monthly_poodium(bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
    topPoopers, err := repo.Get_Monthly_Poodium(db)
	utils.CheckError("Failed to get top poopers", err)

    messageText := "🏆 Poodium for this month 🏆\n" + build_poodium_message(topPoopers)
    msg := tg_bot.NewMessage(chatID, messageText)
    send_message(bot, msg)
}

func send_yearly_poodium(bot *tg_bot.BotAPI, db *sql.DB, chatID int64) {
    topPoopers, err := repo.Get_Yearly_Poodium(db)
	utils.CheckError("Failed to get top poopers", err)

    messageText := "🏆 Poodium for this year 🏆\n" + build_poodium_message(topPoopers)
    msg := tg_bot.NewMessage(chatID, messageText)
    send_message(bot, msg)
}

func main() {
	apiKey := get_api_key()
    if apiKey == "" {
        log.Fatal("TELEGRAM_TOKEN is not set")
    }

	bot, err := tg_bot.NewBotAPI(apiKey)
	utils.CheckError("Failed to create new bot instance", err)

	bot.Debug = true

    updateConfig := tg_bot.NewUpdate(0)
    updateConfig.Timeout = 30
    updates := bot.GetUpdatesChan(updateConfig)

	db := repo.Open_DB_Connection()

	// Schedule the poodium message
    monthlyCron := cron.New()
    _, err = monthlyCron.AddFunc("0 0 1 * *", func() {
        chatID := int64(-4734467277) //CHANGE THIS
        send_monthly_poodium(bot, db, chatID)
    })
	utils.CheckError("Failed to schedule monthly poodium message", err)
    monthlyCron.Start()

	// Schedule the poodium message
    yearlyCron := cron.New()
	_, err = yearlyCron.AddFunc("0 0 1 1 *", func() {
		chatID := int64(-4734467277) //CHANGE THIS
		send_yearly_poodium(bot, db, chatID)
	})
	utils.CheckError("Failed to schedule yearly poodium message", err)
    yearlyCron.Start()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tg_bot.NewMessage(update.Message.Chat.ID, update.Message.Text)
		userId := update.Message.From.ID
		username := update.Message.From.UserName
		messageId := update.Message.MessageID

		if update.Message.Text == "💩" || (update.Message.Sticker != nil && update.Message.Sticker.Emoji == "💩") {
			log.Println("New poop detected!")
			handle_new_poop(db, userId, username, messageId)
		}

		if update.Message.Command() != "" {
			handle_commands(bot, db, update, userId, msg)
		}
	}
}