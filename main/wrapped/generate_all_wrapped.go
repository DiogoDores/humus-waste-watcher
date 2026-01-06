package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"src/config"
	"src/handlers"
	repo "src/repository"

	_ "modernc.org/sqlite"
)

type User struct {
	UserID   int64
	Username string
}

func getAllUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	query := `
	SELECT DISTINCT user_id, username
	FROM poop_tracker
	ORDER BY user_id;
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.Username); err != nil {
			log.Printf("Failed to scan user: %v", err)
			continue
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func generateWrappedForUser(ctx context.Context, r repo.Repository, user User, year int, outputDir string) error {
	yearlyCount, err := r.GetYearlyPoopCount(ctx, user.UserID, year)
	if err != nil {
		return fmt.Errorf("failed to get yearly count: %w", err)
	}

	if yearlyCount < 10 {
		log.Printf("Skipping user %s (%d): only %d poops this year (minimum 10 required)", user.Username, user.UserID, yearlyCount)
		return fmt.Errorf("skipped")
	}

	log.Printf("Generating wrapped for user: %s (%d) - %d poops", user.Username, user.UserID, yearlyCount)

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

	stats := PersonalWrappedStats{
		UserID:     user.UserID,
		Username:   user.Username,
		Year:       year,
		TotalPoops: yearlyCount,
	}

	ranking, err := r.GetYearlyRanking(ctx, user.UserID, year)
	if err == nil {
		stats.GroupRank = ranking
		stats.GroupTotal = ranking.TotalUsers
	} else {
		groupStats, err := r.GetGroupYearlyStats(ctx, year)
		if err == nil && len(groupStats) > 0 {
			stats.GroupTotal = len(groupStats)
			for i, u := range groupStats {
				if u.Username == stats.Username {
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

	stats.MaxStreak, err = r.GetMaxPoopStreak(ctx, user.UserID)
	if err != nil {
		log.Printf("Failed to get streak for %s: %v", user.Username, err)
	}

	stats.DayWithMostPoops, stats.MostPoopsCount, err = r.GetDayWithMostPoops(ctx, user.UserID)
	if err != nil {
		log.Printf("Failed to get day with most poops for %s: %v", user.Username, err)
	}

	stats.DaysWithoutPoop, err = r.GetDaysWithoutPoop(ctx, user.UserID)
	if err != nil {
		log.Printf("Failed to get days without poop for %s: %v", user.Username, err)
	}

	stats.HourDistribution, err = r.GetPoopsByHour(ctx, user.UserID, year)
	if err != nil {
		log.Printf("Failed to get hour distribution for %s: %v", user.Username, err)
	}

	stats.DayOfWeekDistribution, err = r.GetPoopsByDayOfWeek(ctx, user.UserID, year)
	if err != nil {
		log.Printf("Failed to get day of week distribution for %s: %v", user.Username, err)
	}

	allMonthlyStats, err := r.GetMonthlyPoopStats(ctx, user.UserID)
	if err == nil {
		yearStr := fmt.Sprintf("%d", year)
		for _, monthly := range allMonthlyStats {
			if len(monthly.Month) >= 4 && monthly.Month[:4] == yearStr {
				stats.MonthlyStats = append(stats.MonthlyStats, monthly)
			}
		}
	} else {
		log.Printf("Failed to get monthly stats for %s: %v", user.Username, err)
	}

	stats.WeekDayStats, err = r.GetPoopsByWeekAndDay(ctx, user.UserID, year)
	if err != nil {
		log.Printf("Failed to get week/day stats for %s: %v", user.Username, err)
	}

	stats.Personality = handlers.CalculatePersonality(stats.HourDistribution, stats.DayOfWeekDistribution)

	tempDir := filepath.Join("main", "wrapped", "temp")
	os.MkdirAll(tempDir, 0755)

	statsFile := filepath.Join(tempDir, fmt.Sprintf("user_%d_stats.json", user.UserID))
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
		log.Printf("Python script error for %s: %v, output: %s", user.Username, err, string(output))
		return fmt.Errorf("failed to generate slides: %w", err)
	}

	imagesFile := filepath.Join(tempDir, fmt.Sprintf("user_%d_images.json", user.UserID))
	imagesJSON, err := os.ReadFile(imagesFile)
	if err != nil {
		return fmt.Errorf("failed to read images file: %w", err)
	}

	var imagePaths []string
	err = json.Unmarshal(imagesJSON, &imagePaths)
	if err != nil {
		return fmt.Errorf("failed to unmarshal image paths: %w", err)
	}

	// Copy images to output directory (don't delete them)
	userOutputDir := filepath.Join(outputDir, user.Username)
	os.MkdirAll(userOutputDir, 0755)

	for i, imagePath := range imagePaths {
		srcPath := imagePath
		destPath := filepath.Join(userOutputDir, fmt.Sprintf("slide_%02d.png", i+1))

		// Read source file
		data, err := os.ReadFile(srcPath)
		if err != nil {
			log.Printf("Failed to read image %s: %v", srcPath, err)
			continue
		}

		// Write to destination
		err = os.WriteFile(destPath, data, 0644)
		if err != nil {
			log.Printf("Failed to write image %s: %v", destPath, err)
			continue
		}

		log.Printf("  Saved slide %d: %s", i+1, destPath)
	}

	log.Printf("✓ Successfully generated wrapped for %s (%d slides)", user.Username, len(imagePaths))
	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repo.OpenDBConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repository := repo.NewSQLiteRepository(db)
	ctx := context.Background()

	year := 2025
	outputDir := filepath.Join("main", "wrapped", "all_wrapped_output")
	os.MkdirAll(outputDir, 0755)

	log.Printf("Fetching all users from database...")
	users, err := getAllUsers(ctx, db)
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	log.Printf("Found %d users. Generating wrapped slides...\n", len(users))

	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, user := range users {
		log.Printf("\n[%d/%d] Processing user: %s (%d)", i+1, len(users), user.Username, user.UserID)

		err := generateWrappedForUser(ctx, repository, user, year, outputDir)
		if err != nil {
			if err.Error() == "skipped" {
				skipCount++
			} else {
				log.Printf("✗ Error generating wrapped for %s: %v", user.Username, err)
				errorCount++
			}
		} else {
			successCount++
		}
	}

	log.Printf("\n" + strings.Repeat("=", 50))
	log.Printf("Summary:")
	log.Printf("  Total users: %d", len(users))
	log.Printf("  Successfully generated: %d", successCount)
	log.Printf("  Skipped (< 10 poops): %d", skipCount)
	log.Printf("  Errors: %d", errorCount)
	log.Printf("\nOutput directory: %s", outputDir)
	log.Printf(strings.Repeat("=", 50))
}
