package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"src/formatters"
	repo "src/repository"
	"time"

	tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PersonalWrappedStats struct {
	UserID                int64
	Username              string
	Year                  int
	TotalPoops            int
	GroupTotal            int
	GroupRank             repo.YearlyRanking
	MaxStreak             int
	DayWithMostPoops      string
	MostPoopsCount        int
	DaysWithoutPoop       int
	HourDistribution      []repo.HourDistribution
	DayOfWeekDistribution []repo.DayOfWeekDistribution
	Personality           repo.PoopPersonality
	MonthlyStats          []repo.MonthlyPoopCount
	WeekDayStats          []repo.WeekDayPoopCount
}

func HandlePersonalWrapped(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userID int64, msg tg_bot.MessageConfig, groupChatID int64) error {
	year := 2025

	yearlyCount, err := r.GetYearlyPoopCount(ctx, userID, year)
	if err != nil {
		return fmt.Errorf("failed to get yearly count: %w", err)
	}

	log.Printf("User %d has %d poops in year %d", userID, yearlyCount, year)

	if yearlyCount < 10 {
		msg.Text = fmt.Sprintf("You need at least 10 poops this year to get your wrapped\\. Come back when you've logged more\\!\n\nYou currently have %d poops in %d\\.", yearlyCount, year)
		_, err := bot.Send(msg)
		return err
	}

	username := update.Message.From.UserName
	if username == "" {
		username = update.Message.From.FirstName
		if username == "" {
			username = fmt.Sprintf("User %d", userID)
		}
	}

	// If sent in group chat, notify user to check private messages
	if update.Message.Chat.ID == groupChatID {
		var groupMsgText string
		if update.Message.From.UserName != "" {
			// Don't escape username in mention - Telegram handles @username mentions specially
			// Usernames can contain underscores and Telegram will parse them correctly
			groupMsgText = fmt.Sprintf("@%s Check your private messages for your Poop Wrapped\\!", update.Message.From.UserName)
		} else {
			// Escape FirstName since it's not in a mention
			escapedFirstName := formatters.EscapeMarkdownV2(update.Message.From.FirstName)
			groupMsgText = fmt.Sprintf("%s Check your private messages for your Poop Wrapped\\!", escapedFirstName)
		}
		groupMsg := tg_bot.NewMessage(groupChatID, groupMsgText)
		groupMsg.ParseMode = tg_bot.ModeMarkdownV2
		_, err := bot.Send(groupMsg)
		if err != nil {
			log.Printf("Failed to send group notification: %v", err)
		}
	}

	stats := PersonalWrappedStats{
		UserID:     userID,
		Username:   username,
		Year:       year,
		TotalPoops: yearlyCount,
	}

	ranking, err := r.GetYearlyRanking(ctx, userID, year)
	if err == nil {
		stats.GroupRank = ranking
		stats.GroupTotal = ranking.TotalUsers
	} else {
		groupStats, err := r.GetGroupYearlyStats(ctx, year)
		if err == nil && len(groupStats) > 0 {
			stats.GroupTotal = len(groupStats)
			for i, user := range groupStats {
				if user.Username == stats.Username {
					stats.GroupRank = repo.YearlyRanking{
						Rank:       i + 1,
						TotalUsers: len(groupStats),
						Percentage: float64(len(groupStats)-i-1) / float64(len(groupStats)) * 100.0,
					}
					break
				}
			}
		}
	}

	stats.MaxStreak, err = r.GetMaxPoopStreak(ctx, userID)
	if err != nil {
		log.Printf("Failed to get streak: %v", err)
	}

	stats.DayWithMostPoops, stats.MostPoopsCount, err = r.GetDayWithMostPoops(ctx, userID)
	if err != nil {
		log.Printf("Failed to get day with most poops: %v", err)
	}

	stats.DaysWithoutPoop, err = r.GetDaysWithoutPoop(ctx, userID)
	if err != nil {
		log.Printf("Failed to get days without poop: %v", err)
	}

	stats.HourDistribution, err = r.GetPoopsByHour(ctx, userID, year)
	if err != nil {
		log.Printf("Failed to get hour distribution: %v", err)
	}

	stats.DayOfWeekDistribution, err = r.GetPoopsByDayOfWeek(ctx, userID, year)
	if err != nil {
		log.Printf("Failed to get day of week distribution: %v", err)
	}

	allMonthlyStats, err := r.GetMonthlyPoopStats(ctx, userID)
	if err == nil {
		yearStr := fmt.Sprintf("%d", year)
		for _, monthly := range allMonthlyStats {
			if len(monthly.Month) >= 4 && monthly.Month[:4] == yearStr {
				stats.MonthlyStats = append(stats.MonthlyStats, monthly)
			}
		}
	} else {
		log.Printf("Failed to get monthly stats: %v", err)
	}

	stats.WeekDayStats, err = r.GetPoopsByWeekAndDay(ctx, userID, year)
	if err != nil {
		log.Printf("Failed to get week/day stats: %v", err)
	}

	stats.Personality = CalculatePersonality(stats.HourDistribution, stats.DayOfWeekDistribution)

	tempDir := filepath.Join("main", "wrapped", "temp")
	os.MkdirAll(tempDir, 0755)

	statsFile := filepath.Join(tempDir, fmt.Sprintf("user_%d_stats.json", userID))
	statsJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	err = os.WriteFile(statsFile, statsJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write stats file: %w", err)
	}

	imagesFile := filepath.Join(tempDir, fmt.Sprintf("user_%d_images.json", userID))
	var imagePaths []string

	// Ensure cleanup happens even if there's an error
	defer func() {
		os.Remove(statsFile)
		os.Remove(imagesFile)
		for _, imgPath := range imagePaths {
			os.Remove(imgPath)
		}
	}()

	pythonScript := filepath.Join("main", "wrapped", "generate_personal_wrapped.py")
	cmd := exec.Command("python3", pythonScript, statsFile)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Python script error: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to generate slides: %w", err)
	}

	imagesJSON, err := os.ReadFile(imagesFile)
	if err != nil {
		return fmt.Errorf("failed to read images file: %w", err)
	}

	err = json.Unmarshal(imagesJSON, &imagePaths)
	if err != nil {
		return fmt.Errorf("failed to unmarshal image paths: %w", err)
	}

	// Try to send the first image to check if we can message the user
	firstPhoto := tg_bot.NewPhoto(userID, tg_bot.FilePath(imagePaths[0]))
	_, err = bot.Send(firstPhoto)
	if err != nil {
		log.Printf("Failed to send wrapped to user %d: %v", userID, err)
		// If we can't send to the user's private chat, notify them in the group
		if update.Message.Chat.ID == groupChatID {
			var errorMsgText string
			if update.Message.From.UserName != "" {
				errorMsgText = fmt.Sprintf("@%s I can't send you a private message\\. Please start a chat with me first by clicking [here](tg://user?id=%d), then try the command again\\!",
					update.Message.From.UserName, userID)
			} else {
				escapedFirstName := formatters.EscapeMarkdownV2(update.Message.From.FirstName)
				errorMsgText = fmt.Sprintf("%s I can't send you a private message\\. Please start a chat with me first by clicking [here](tg://user?id=%d), then try the command again\\!",
					escapedFirstName, userID)
			}
			errorMsg := tg_bot.NewMessage(groupChatID, errorMsgText)
			errorMsg.ParseMode = tg_bot.ModeMarkdownV2
			errorMsg.DisableWebPagePreview = true
			_, sendErr := bot.Send(errorMsg)
			if sendErr != nil {
				log.Printf("Failed to send error message to group: %v", sendErr)
			}
		}
		return fmt.Errorf("cannot send private message to user %d: %w", userID, err)
	}

	// Send remaining images
	for i := 1; i < len(imagePaths); i++ {
		photo := tg_bot.NewPhoto(userID, tg_bot.FilePath(imagePaths[i]))
		_, err := bot.Send(photo)
		if err != nil {
			log.Printf("Failed to send slide %d: %v", i+1, err)
			continue
		}

		if i < len(imagePaths)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil
}

func CalculatePersonality(hourDist []repo.HourDistribution, dayDist []repo.DayOfWeekDistribution) repo.PoopPersonality {
	hourMap := make(map[int]int)
	total := 0
	for _, h := range hourDist {
		hourMap[h.Hour] = h.PoopCount
		total += h.PoopCount
	}

	if total == 0 {
		return repo.PoopPersonality{
			Type:        "Unknown",
			Description: "Not enough data",
			Confidence:  0.0,
		}
	}

	// Calculate morning ratio (5-10)
	morningCount := 0
	for h := 5; h <= 10; h++ {
		morningCount += hourMap[h]
	}
	morningRatio := float64(morningCount) / float64(total)

	// Calculate afternoon ratio (12-17)
	afternoonCount := 0
	for h := 12; h <= 17; h++ {
		afternoonCount += hourMap[h]
	}
	afternoonRatio := float64(afternoonCount) / float64(total)

	// Calculate evening ratio (18-22)
	eveningCount := 0
	for h := 18; h <= 22; h++ {
		eveningCount += hourMap[h]
	}
	eveningRatio := float64(eveningCount) / float64(total)

	// Calculate night ratio (23-4)
	nightCount := 0
	for h := 23; h <= 23; h++ {
		nightCount += hourMap[h]
	}
	for h := 0; h <= 4; h++ {
		nightCount += hourMap[h]
	}
	nightRatio := float64(nightCount) / float64(total)

	// Calculate weekend ratio
	dayMap := make(map[string]int)
	weekendTotal := 0
	for _, d := range dayDist {
		dayMap[d.DayOfTheWeek] = d.PoopCount
		if d.DayOfTheWeek == "Saturday" || d.DayOfTheWeek == "Sunday" {
			weekendTotal += d.PoopCount
		}
	}
	weekendRatio := float64(weekendTotal) / float64(total)

	// Prioritize the most distinctive trait (highest ratio wins)
	var personality repo.PoopPersonality

	// Find the highest ratio among all patterns
	maxRatio := morningRatio
	winner := "morning"
	if afternoonRatio > maxRatio {
		maxRatio = afternoonRatio
		winner = "afternoon"
	}
	if eveningRatio > maxRatio {
		maxRatio = eveningRatio
		winner = "evening"
	}
	if nightRatio > maxRatio {
		maxRatio = nightRatio
		winner = "night"
	}
	if weekendRatio > maxRatio {
		maxRatio = weekendRatio
		winner = "weekend"
	}

	// Apply thresholds - lower for small groups to ensure variety
	if winner == "morning" && morningRatio >= 0.3 {
		personality = repo.PoopPersonality{
			Type:        "Morning Menace",
			Description: fmt.Sprintf("You're a creature of habit. Your bowels wake up before you do.\n(%.0f%% of your poops happen between 5-10 AM)", morningRatio*100),
			Confidence:  morningRatio,
		}
	} else if winner == "afternoon" && afternoonRatio >= 0.3 {
		personality = repo.PoopPersonality{
			Type:        "Afternoon Pooper",
			Description: fmt.Sprintf("You're a lunch break champion. Midday is your time to shine.\n(%.0f%% of your poops happen between 12-5 PM)", afternoonRatio*100),
			Confidence:  afternoonRatio,
		}
	} else if winner == "evening" && eveningRatio >= 0.3 {
		personality = repo.PoopPersonality{
			Type:        "Evening Pooper",
			Description: fmt.Sprintf("You unwind after work. The evening toilet is your sanctuary.\n(%.0f%% of your poops happen between 6-10 PM)", eveningRatio*100),
			Confidence:  eveningRatio,
		}
	} else if winner == "night" && nightRatio >= 0.25 {
		personality = repo.PoopPersonality{
			Type:        "Night Pooper",
			Description: fmt.Sprintf("You prefer the quiet hours. The toilet is your midnight companion.\n(%.0f%% of your poops happen between 11 PM-4 AM)", nightRatio*100),
			Confidence:  nightRatio,
		}
	} else if winner == "weekend" && weekendRatio >= 0.3 {
		personality = repo.PoopPersonality{
			Type:        "Weekend Warrior",
			Description: fmt.Sprintf("You save it for the weekend. Work can wait.\n(%.0f%% of your poops happen on Saturday and Sunday)", weekendRatio*100),
			Confidence:  weekendRatio,
		}
	} else {
		maxWindowCount := 0
		var maxWindowStart int
		for startHour := 0; startHour < 24; startHour++ {
			windowCount := 0
			for offset := 0; offset < 3; offset++ {
				hour := (startHour + offset) % 24
				windowCount += hourMap[hour]
			}
			if windowCount > maxWindowCount {
				maxWindowCount = windowCount
				maxWindowStart = startHour
			}
		}
		concentrationRatio := float64(maxWindowCount) / float64(total)

		// Clockwork Pooper is the default
		personality = repo.PoopPersonality{
			Type:        "Clockwork Pooper",
			Description: fmt.Sprintf("Regular. Reliable. Your toilet fears you.\n(%.0f%% of your poops happen between %d:00-%d:00)", concentrationRatio*100, maxWindowStart, (maxWindowStart+2)%24),
			Confidence:  concentrationRatio,
		}
	}

	return personality
}
