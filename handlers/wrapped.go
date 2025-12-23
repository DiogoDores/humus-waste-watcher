package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	repo "src/repository"
	"time"

	tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PersonalWrappedStats struct {
	UserID 			 int64
	Username 		 string
	Year 			 int
	TotalPoops 		 int
	GroupTotal 		 int
	GroupRank 		 repo.YearlyRanking
	MaxStreak 		 int
	DayWithMostPoops string
	MostPoopsCount   int
	DaysWithoutPoop  int
}

func HandlePersonalWrapped(ctx context.Context, bot *tg_bot.BotAPI, r repo.Repository, update tg_bot.Update, userID int64, msg tg_bot.MessageConfig) error {
	year := 2025

	yearlyCount, err := r.GetYearlyPoopCount(ctx, userID, year)
	if err != nil {
		return fmt.Errorf("failed to get yearly count: %w", err)
	}

	stats := PersonalWrappedStats {
		UserID: userID,
		Username: update.Message.From.UserName,
		Year: year,
		TotalPoops: yearlyCount,
	}

	groupStats, err := r.GetGroupYearlyStats(ctx, year)
	if err == nil && len(groupStats) > 0 {
		stats.GroupTotal = len(groupStats)
		for i, user := range groupStats {
			if user.Username == stats.Username {
				stats.GroupRank = repo.YearlyRanking{
					Rank: i + 1,
					TotalUsers: len(groupStats),
					Percentage: float64(len(groupStats) - i - 1) / float64(len(groupStats)) * 100.0,
				}
				break
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

	pythonScript := filepath.Join("main", "wrapped", "generate_personal_wrapped.py")
	cmd := exec.Command("python3", pythonScript, statsFile)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Python script error: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to generate slides: %w", err)
	}

	imagesFile := filepath.Join(tempDir, fmt.Sprintf("user_%d_images.json", userID))
	imagesJSON, err := os.ReadFile(imagesFile)
	if err != nil {
		return fmt.Errorf("failed to read images file: %w", err)
	}

	var imagePaths []string
	err = json.Unmarshal(imagesJSON, &imagePaths)
	if err != nil {
		return fmt.Errorf("failed to unmarshal image paths: %w", err)
	}

	for i, imagePath := range imagePaths {
		photo := tg_bot.NewPhoto(update.Message.Chat.ID, tg_bot.FilePath(imagePath))
		_, err := bot.Send(photo)
		if err != nil {
			log.Printf("Failed to send slide %d: %v", i+1, err)
			continue
		}
		
		if i < len(imagePaths)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	defer func() {
		os.Remove(statsFile)
		os.Remove(imagesFile)
		for _, imgPath := range imagePaths {
			os.Remove(imgPath)
		}
	}()

	return nil
}