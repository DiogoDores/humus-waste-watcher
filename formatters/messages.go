package formatters

import (
	"fmt"
	"strings"
	"time"

	repo "src/repository"
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

func EscapeMarkdownV2(s string) string {
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "-", "\\-")
	s = strings.ReplaceAll(s, ".", "\\.")
	s = strings.ReplaceAll(s, "!", "\\!")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	s = strings.ReplaceAll(s, "{", "\\{")
	s = strings.ReplaceAll(s, "}", "\\}")
	s = strings.ReplaceAll(s, "=", "\\=")
	s = strings.ReplaceAll(s, "+", "\\+")
	s = strings.ReplaceAll(s, ">", "\\>")
	s = strings.ReplaceAll(s, "<", "\\<")
	s = strings.ReplaceAll(s, "|", "\\|")
	return s
}

func BuildPoodiumMessage(topPoopers []repo.UserPoopCount) string {
	if len(topPoopers) < 3 {
		return "Not enough data for poodium\\."
	}

	escape := EscapeMarkdownV2
	return "🥇 " + escape(topPoopers[0].Username) + " \\- " + fmt.Sprint(topPoopers[0].PoopCount) + "💩\n" +
		"🥈 " + escape(topPoopers[1].Username) + " \\- " + fmt.Sprint(topPoopers[1].PoopCount) + "💩\n" +
		"🥉 " + escape(topPoopers[2].Username) + " \\- " + fmt.Sprint(topPoopers[2].PoopCount) + "💩"
}

func FormatPoopLog(username string, globalPoopCount int, monthlyPoopCounts []repo.MonthlyPoopCount, daysWithoutPoop int, maxStreak int, day string, poops int) string {
	year := time.Now().Year()
	monthlyAverages := make(map[string]float64)
	for _, mpc := range monthlyPoopCounts {
		yearMonth := mpc.Month
		month := yearMonth[5:]
		daysInMonthCount := daysInMonth(monthNames[month], year)

		if time.Now().Month().String() == monthNames[month] {
			daysInMonthCount = time.Now().Day()
		}

		monthlyAverages[month] = float64(mpc.PoopCount) / float64(daysInMonthCount)
	}

	yearlyAverage := float64(globalPoopCount) / float64(time.Now().YearDay())

	escapedUsername := EscapeMarkdownV2(username)
	msg := fmt.Sprintf("*💩 Poop Report for @%s 💩*\n\n", escapedUsername)
	msg += "*📅 Yearly Overview:*\n"
	msg += fmt.Sprintf("🟤 Total dumps: `%d`\n", globalPoopCount)
	msg += fmt.Sprintf("📊 Average per day: `%.2f`\n", yearlyAverage)
	msg += fmt.Sprintf("🚫 Days without poops: `%d`\n", daysWithoutPoop)
	msg += fmt.Sprintf("🔥 Max poop streak: `%d`\n", maxStreak)
	msg += fmt.Sprintf("💣 Day with most poops: `%s with %d poops`\n\n", EscapeMarkdownV2(day), poops)
	msg += "*📅 Monthly Breakdown:\n*"

	for _, mpc := range monthlyPoopCounts {
		yearMonth := mpc.Month
		month := yearMonth[5:]
		msg += fmt.Sprintf("🗓 %s:  `%d poops`   \\(📊 Avg:   `%.2f per day`\\)\n", monthNames[month], mpc.PoopCount, monthlyAverages[month])
	}

	return msg
}

func FormatLeaderboard(leaderboard []repo.UserPoopCount) string {
	msg := "This month's leaderboard:\n"
	for _, user := range leaderboard {
		escapedUsername := EscapeMarkdownV2(user.Username)
		msg += fmt.Sprintf("\t\t\t• %s \\- %d💩\n", escapedUsername, user.PoopCount)
	}
	return msg
}

func FormatHelpMessage(isUnknownCommand bool) string {
	message := ""
	if isUnknownCommand {
		message += "Sorry, I don't recognize that command\\. "
	}
	message += "Here are the commands I understand:\n" +
		"\t\t\t\t• _/help_ \\- Get a list of available commands\n" +
		"\t\t\t\t• _/my\\_poop\\_log_ \\- Get your personal monthly poop statistics\n" +
		"\t\t\t\t• _/leaderboard_ \\- Get the monthly leaderboard\n" +
		"\t\t\t\t• _/bottom\\_poopers_ \\- Get the reverse poodium\n" +
		"\t\t\t\t• _/poodium_ \\- Get the monthly poodium\n" +
		"\t\t\t\t• _/poodium\\_year_ \\- Get the yearly poodium"
	return message
}

func FormatPoodiumTitle(monthName string) string {
	return "🏆 Poodium for " + monthName + " 🏆\n"
}

func FormatYearlyPoodiumTitle(year int) string {
	return fmt.Sprintf("🏆 Poodium for %d 🏆\n", year)
}

func GetMonthName(monthStr string) string {
	return monthNames[monthStr]
}
