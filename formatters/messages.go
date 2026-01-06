package formatters

import (
	"fmt"
	"log"
	"strings"
	"time"

	repo "src/repository"
)

// daysInMonth calculates the number of days in a given month and year using time package
func daysInMonth(year int, month time.Month) int {
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	return int(firstOfNextMonth.Sub(firstOfMonth).Hours() / 24)
}

// parseMonthString parses a month string (e.g., "01", "02") to time.Month
func parseMonthString(monthStr string) (time.Month, error) {
	switch monthStr {
	case "01":
		return time.January, nil
	case "02":
		return time.February, nil
	case "03":
		return time.March, nil
	case "04":
		return time.April, nil
	case "05":
		return time.May, nil
	case "06":
		return time.June, nil
	case "07":
		return time.July, nil
	case "08":
		return time.August, nil
	case "09":
		return time.September, nil
	case "10":
		return time.October, nil
	case "11":
		return time.November, nil
	case "12":
		return time.December, nil
	default:
		return time.January, fmt.Errorf("invalid month string: %s", monthStr)
	}
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
	now := time.Now()
	year := now.Year()
	monthlyAverages := make(map[string]float64)

	for _, mpc := range monthlyPoopCounts {
		yearMonth := mpc.Month
		monthStr := yearMonth[5:]
		month, err := parseMonthString(monthStr)
		if err != nil {
			log.Printf("Invalid month string: %s", monthStr)
			continue
		}

		daysInMonthCount := daysInMonth(year, month)

		if now.Year() == year && now.Month() == month {
			daysInMonthCount = now.Day()
		}

		monthlyAverages[monthStr] = float64(mpc.PoopCount) / float64(daysInMonthCount)
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
		monthStr := yearMonth[5:]
		month, err := parseMonthString(monthStr)
		if err != nil {
			log.Printf("Invalid month string: %s", monthStr)
			continue
		}
		monthName := month.String()
		msg += fmt.Sprintf("🗓 %s:  `%d poops`   \\(📊 Avg:   `%.2f per day`\\)\n", monthName, mpc.PoopCount, monthlyAverages[monthStr])
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
		"\t\t\t\t• _/poodium\\_year_ \\- Get the yearly poodium\n" +
		"\t\t\t\t• _/poop\\_wrapped_ \\- Get your personalized Poop Wrapped"
	return message
}

func FormatPoodiumTitle(monthName string) string {
	return "🏆 Poodium for " + monthName + " 🏆\n"
}

func FormatYearlyPoodiumTitle(year int) string {
	return fmt.Sprintf("🏆 Poodium for %d 🏆\n", year)
}

func GetMonthName(monthStr string) string {
	month, err := parseMonthString(monthStr)
	if err != nil {
		log.Printf("Invalid month string: %s", monthStr)
		return "Unknown"
	}
	return month.String()
}

func FormatGroupWrappedTitle(period string, year int) string {
	return fmt.Sprintf("🎉 Poop Wrapped Awards %s %d 🎉\n\n", period, year)
}

func BuildGroupAwardsMessage(awards []repo.GroupAward) string {
	if len(awards) == 0 {
		return "Not enough data for awards\\."
	}

	escape := EscapeMarkdownV2
	msg := "🏆 *Awards:*\n\n"

	for _, award := range awards {
		// Ensure there's exactly one space after the emoji
		emoji := award.Emoji
		if emoji != "" && !strings.HasSuffix(emoji, " ") {
			emoji += " "
		}
		msg += fmt.Sprintf("%s*%s*: @%s \\(%s\\)\n",
			emoji,
			escape(award.AwardName),
			escape(award.Winner),
			escape(award.Value))
	}

	return msg
}
