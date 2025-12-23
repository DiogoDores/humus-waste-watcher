package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// setupTestDB creates a temporary in-memory database for testing
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	ctx := context.Background()
	err = createTable(ctx, db)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// insertTestData inserts sample data for testing
func insertTestData(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// User 1: Alice - Regular pooper, morning person
	// 5 poops in January 2025, mostly morning (8-10 AM)
	aliceID := int64(1001)
	aliceUsername := "alice"
	for i := 0; i < 5; i++ {
		timestamp := baseTime.AddDate(0, 0, i*7).Add(time.Hour * time.Duration(8+i%3))
		err := LogPoop(ctx, db, aliceID, aliceUsername, int64(1000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// User 2: Bob - Night owl, more active
	// 10 poops in January 2025, mostly night (11 PM - 2 AM), some consecutive days
	bobID := int64(1002)
	bobUsername := "bob"
	// Create a 3-day streak
	for day := 0; day < 3; day++ {
		timestamp := baseTime.AddDate(0, 0, day).Add(time.Hour * 23) // 11 PM each day
		err := LogPoop(ctx, db, bobID, bobUsername, int64(2000+day), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}
	// Add more poops throughout January
	for i := 0; i < 7; i++ {
		timestamp := baseTime.AddDate(0, 0, 5+i*3).Add(time.Hour * time.Duration(1+i%2)) // 1-2 AM
		err := LogPoop(ctx, db, bobID, bobUsername, int64(2010+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// User 3: Charlie - Weekend warrior, machine gun (multiple poops in one day)
	// 8 poops in January 2025, mostly weekends
	charlieID := int64(1003)
	charlieUsername := "charlie"
	// Weekend poops (Saturday = 6, Sunday = 0)
	weekendDays := []int{4, 5, 11, 12, 18, 19} // Saturdays and Sundays in January 2025
	for i, day := range weekendDays {
		timestamp := baseTime.AddDate(0, 0, day).Add(time.Hour * 10) // 10 AM
		err := LogPoop(ctx, db, charlieID, charlieUsername, int64(3000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}
	// Add 2 more weekend poops
	for i := 0; i < 2; i++ {
		timestamp := baseTime.AddDate(0, 0, weekendDays[i]).Add(time.Hour * 15) // 3 PM same day
		err := LogPoop(ctx, db, charlieID, charlieUsername, int64(3010+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// User 4: Machine Gun - Most poops in single day (5 poops on Jan 15)
	machineGunID := int64(1004)
	machineGunUsername := "machine_gun"
	bigDay := baseTime.AddDate(0, 0, 14) // Jan 15
	for i := 0; i < 5; i++ {
		timestamp := bigDay.Add(time.Hour * time.Duration(8+i*3)) // 8 AM, 11 AM, 2 PM, 5 PM, 8 PM
		err := LogPoop(ctx, db, machineGunID, machineGunUsername, int64(4000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// User 5: Early Bird - Most poops between 5-8 AM
	earlyBirdID := int64(1005)
	earlyBirdUsername := "early_bird"
	for i := 0; i < 6; i++ {
		timestamp := baseTime.AddDate(0, 0, i*4).Add(time.Hour * time.Duration(6+i%3)) // 6-8 AM
		err := LogPoop(ctx, db, earlyBirdID, earlyBirdUsername, int64(5000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// User 6: Company Time Pooper - Most poops during work hours (9 AM - 6 PM)
	companyTimeID := int64(1006)
	companyTimeUsername := "company_time"
	companyHours := []int{9, 10, 11, 12, 13, 14, 15, 16} // 8 hours within 9-18 range
	for i := 0; i < 8; i++ {
		hour := companyHours[i%len(companyHours)]
		timestamp := baseTime.AddDate(0, 0, i*3).Add(time.Hour * time.Duration(hour))
		err := LogPoop(ctx, db, companyTimeID, companyTimeUsername, int64(6000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// Add some data for 2024 to test year filtering
	baseTime2024 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		timestamp := baseTime2024.AddDate(0, 0, i*10).Add(time.Hour * 12)
		err := LogPoop(ctx, db, aliceID, aliceUsername, int64(7000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}
}

func TestGetYearlyPoopCount(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	insertTestData(t, db)

	ctx := context.Background()

	tests := []struct {
		name     string
		userID   int64
		year     int
		expected int
	}{
		{"Alice 2025", 1001, 2025, 5},
		{"Bob 2025", 1002, 2025, 10},
		{"Charlie 2025", 1003, 2025, 8},
		{"Alice 2024", 1001, 2024, 3},
		{"Non-existent user", 9999, 2025, 0},
		{"User with no data in year", 1001, 2023, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetYearlyPoopCount(ctx, db, tt.userID, tt.year)
			if err != nil {
				t.Fatalf("GetYearlyPoopCount() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("GetYearlyPoopCount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPoopsByHour(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	insertTestData(t, db)

	ctx := context.Background()

	result, err := GetPoopsByHour(ctx, db, 1001, 2025) // Alice
	if err != nil {
		t.Fatalf("GetPoopsByHour() error = %v", err)
	}

	// Should return all 24 hours
	if len(result) != 24 {
		t.Errorf("GetPoopsByHour() returned %d hours, want 24", len(result))
	}

	// Check that hours are in order
	for i, hd := range result {
		if hd.Hour != i {
			t.Errorf("GetPoopsByHour() hour at index %d = %d, want %d", i, hd.Hour, i)
		}
	}

	// Alice should have poops in hours 8-10 (morning)
	totalPoops := 0
	for _, hd := range result {
		totalPoops += hd.PoopCount
	}
	if totalPoops != 5 {
		t.Errorf("GetPoopsByHour() total poops = %d, want 5", totalPoops)
	}

	// Check that hours 8-10 have some poops
	hasMorningPoops := false
	for _, hd := range result {
		if hd.Hour >= 8 && hd.Hour <= 10 && hd.PoopCount > 0 {
			hasMorningPoops = true
			break
		}
	}
	if !hasMorningPoops {
		t.Error("GetPoopsByHour() expected morning poops (8-10 AM) but found none")
	}
}

func TestGetPoopsByDayOfWeek(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	insertTestData(t, db)

	ctx := context.Background()

	result, err := GetPoopsByDayOfWeek(ctx, db, 1003, 2025) // Charlie (weekend warrior)
	if err != nil {
		t.Fatalf("GetPoopsByDayOfWeek() error = %v", err)
	}

	// Should return all 7 days
	if len(result) != 7 {
		t.Errorf("GetPoopsByDayOfWeek() returned %d days, want 7", len(result))
	}

	// Check day order (Monday first, Sunday last)
	expectedOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i, dotw := range result {
		if dotw.DayOfTheWeek != expectedOrder[i] {
			t.Errorf("GetPoopsByDayOfWeek() day at index %d = %s, want %s", i, dotw.DayOfTheWeek, expectedOrder[i])
		}
	}

	// Charlie should have more weekend poops
	totalPoops := 0
	weekendPoops := 0
	for _, dotw := range result {
		totalPoops += dotw.PoopCount
		if dotw.DayOfTheWeek == "Saturday" || dotw.DayOfTheWeek == "Sunday" {
			weekendPoops += dotw.PoopCount
		}
	}
	if totalPoops != 8 {
		t.Errorf("GetPoopsByDayOfWeek() total poops = %d, want 8", totalPoops)
	}
	if weekendPoops == 0 {
		t.Error("GetPoopsByDayOfWeek() expected weekend poops but found none")
	}
}

func TestGetYearlyRanking(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	insertTestData(t, db)

	ctx := context.Background()

	// Bob has 10 poops, should rank higher than Alice (5 poops)
	bobRanking, err := GetYearlyRanking(ctx, db, 1002, 2025) // Bob
	if err != nil {
		t.Fatalf("GetYearlyRanking() error = %v", err)
	}

	aliceRanking, err := GetYearlyRanking(ctx, db, 1001, 2025) // Alice
	if err != nil {
		t.Fatalf("GetYearlyRanking() error = %v", err)
	}

	// Bob should rank higher (lower rank number = better)
	if bobRanking.Rank >= aliceRanking.Rank {
		t.Errorf("Bob rank (%d) should be better than Alice rank (%d)", bobRanking.Rank, aliceRanking.Rank)
	}

	// Both should have same total users
	if bobRanking.TotalUsers != aliceRanking.TotalUsers {
		t.Errorf("TotalUsers mismatch: Bob=%d, Alice=%d", bobRanking.TotalUsers, aliceRanking.TotalUsers)
	}

	// Percentage should be valid (0-100)
	if bobRanking.Percentage < 0 || bobRanking.Percentage > 100 {
		t.Errorf("Bob percentage = %f, should be between 0 and 100", bobRanking.Percentage)
	}

	// Top user should have highest percentage
	if bobRanking.Percentage <= aliceRanking.Percentage {
		t.Errorf("Bob percentage (%f) should be higher than Alice (%f)", bobRanking.Percentage, aliceRanking.Percentage)
	}
}

func TestGetGroupYearlyStats(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	insertTestData(t, db)

	ctx := context.Background()

	result, err := GetGroupYearlyStats(ctx, db, 2025)
	if err != nil {
		t.Fatalf("GetGroupYearlyStats() error = %v", err)
	}

	// Should have 6 users (added company_time user)
	if len(result) != 6 {
		t.Errorf("GetGroupYearlyStats() returned %d users, want 6", len(result))
	}

	// Should be sorted by poop count descending
	for i := 0; i < len(result)-1; i++ {
		if result[i].PoopCount < result[i+1].PoopCount {
			t.Errorf("GetGroupYearlyStats() not sorted: %d < %d at index %d", result[i].PoopCount, result[i+1].PoopCount, i)
		}
	}

	// Bob should be first (10 poops)
	if result[0].Username != "bob" || result[0].PoopCount != 10 {
		t.Errorf("GetGroupYearlyStats() first user = %s (%d poops), want bob (10 poops)", result[0].Username, result[0].PoopCount)
	}
}

func TestGetGroupAwards(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	insertTestData(t, db)

	ctx := context.Background()

	awards, err := GetGroupAwards(ctx, db, 2025)
	if err != nil {
		t.Fatalf("GetGroupAwards() error = %v", err)
	}

	// Should have at least some awards
	if len(awards) == 0 {
		t.Error("GetGroupAwards() returned no awards")
	}

	// Check for specific awards
	awardMap := make(map[string]GroupAward)
	for _, award := range awards {
		awardMap[award.AwardName] = award
	}

	// Early Bird award
	if earlyBird, ok := awardMap["Early Bird"]; ok {
		if earlyBird.Winner != "early_bird" {
			t.Errorf("Early Bird winner = %s, want early_bird", earlyBird.Winner)
		}
		if earlyBird.Emoji != "☀️" {
			t.Errorf("Early Bird emoji = %s, want ☀️", earlyBird.Emoji)
		}
	}

	// Machine Gun award
	if machineGun, ok := awardMap["Machine Gun"]; ok {
		if machineGun.Winner != "machine_gun" {
			t.Errorf("Machine Gun winner = %s, want machine_gun", machineGun.Winner)
		}
		if machineGun.Value != "5" {
			t.Errorf("Machine Gun value = %s, want 5", machineGun.Value)
		}
	}

	// Consistency King (should be Bob with 3-day streak)
	if consistency, ok := awardMap["Consistency King"]; ok {
		if consistency.Winner != "bob" {
			t.Errorf("Consistency King winner = %s, want bob", consistency.Winner)
		}
		if consistency.Value != "3" {
			t.Errorf("Consistency King value = %s, want 3", consistency.Value)
		}
	}

	// Weekend Warrior (should be Charlie)
	if weekend, ok := awardMap["Weekend Warrior"]; ok {
		if weekend.Winner != "charlie" {
			t.Errorf("Weekend Warrior winner = %s, want charlie", weekend.Winner)
		}
	}

	// Boss makes a dollar, I make a dime (should be company_time)
	if companyTime, ok := awardMap["Boss makes a dollar, I make a dime"]; ok {
		if companyTime.Winner != "company_time" {
			t.Errorf("Company Time winner = %s, want company_time", companyTime.Winner)
		}
		if companyTime.Emoji != "💰" {
			t.Errorf("Company Time emoji = %s, want 💰", companyTime.Emoji)
		}
		if companyTime.Value != "8" {
			t.Errorf("Company Time value = %s, want 8", companyTime.Value)
		}
	} else {
		t.Error("Company Time award not found in results")
	}
}

func TestCompanyTimeAward(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	insertTestData(t, db)

	ctx := context.Background()

	// Test the company time award specifically
	awards, err := GetGroupAwards(ctx, db, 2025)
	if err != nil {
		t.Fatalf("GetGroupAwards() error = %v", err)
	}

	// Find the company time award
	var companyTimeAward *GroupAward
	for i := range awards {
		if awards[i].AwardName == "Boss makes a dollar, I make a dime" {
			companyTimeAward = &awards[i]
			break
		}
	}

	if companyTimeAward == nil {
		t.Fatal("Company Time award not found")
	}

	// Verify the winner
	if companyTimeAward.Winner != "company_time" {
		t.Errorf("Company Time winner = %s, want company_time", companyTimeAward.Winner)
	}

	// Verify the count (should be 8 poops during 9-18)
	if companyTimeAward.Value != "8" {
		t.Errorf("Company Time value = %s, want 8", companyTimeAward.Value)
	}

	// Verify emoji
	if companyTimeAward.Emoji != "💰" {
		t.Errorf("Company Time emoji = %s, want 💰", companyTimeAward.Emoji)
	}
}

func TestCompanyTimeAward_EdgeCases(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test with user who has poops ONLY during company time
	userID := int64(2000)
	for i := 0; i < 5; i++ {
		// All poops at 10 AM (company time)
		timestamp := baseTime.AddDate(0, 0, i).Add(time.Hour * 10)
		err := LogPoop(ctx, db, userID, "workaholic", int64(10000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// Test with user who has poops ONLY outside company time
	userID2 := int64(2001)
	for i := 0; i < 3; i++ {
		// All poops at 7 AM (before company time)
		timestamp := baseTime.AddDate(0, 0, i).Add(time.Hour * 7)
		err := LogPoop(ctx, db, userID2, "early_riser", int64(20000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	awards, err := GetGroupAwards(ctx, db, 2025)
	if err != nil {
		t.Fatalf("GetGroupAwards() error = %v", err)
	}

	// Find the company time award
	var companyTimeAward *GroupAward
	for i := range awards {
		if awards[i].AwardName == "Boss makes a dollar, I make a dime" {
			companyTimeAward = &awards[i]
			break
		}
	}

	if companyTimeAward == nil {
		t.Fatal("Company Time award not found")
	}

	// Should be workaholic (5 poops) not early_riser (3 poops, but outside company time)
	if companyTimeAward.Winner != "workaholic" {
		t.Errorf("Company Time winner = %s, want workaholic", companyTimeAward.Winner)
	}
	if companyTimeAward.Value != "5" {
		t.Errorf("Company Time value = %s, want 5", companyTimeAward.Value)
	}
}

func TestGetYearlyRanking_EdgeCases(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Test with single user
	aliceID := int64(1001)
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	err := LogPoop(ctx, db, aliceID, "alice", 1, baseTime.Format("2006-01-02 15:04:05"), baseTime.Unix())
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	ranking, err := GetYearlyRanking(ctx, db, aliceID, 2025)
	if err != nil {
		t.Fatalf("GetYearlyRanking() error = %v", err)
	}

	if ranking.Rank != 1 {
		t.Errorf("Single user rank = %d, want 1", ranking.Rank)
	}
	if ranking.TotalUsers != 1 {
		t.Errorf("Single user totalUsers = %d, want 1", ranking.TotalUsers)
	}
	if ranking.Percentage != 0.0 {
		t.Errorf("Single user percentage = %f, want 0.0", ranking.Percentage)
	}
}

func TestGetPoopsByHour_AllHoursPresent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Insert one poop at each hour
	userID := int64(2000)
	for hour := 0; hour < 24; hour++ {
		timestamp := baseTime.Add(time.Hour * time.Duration(hour))
		err := LogPoop(ctx, db, userID, "testuser", int64(10000+hour), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	result, err := GetPoopsByHour(ctx, db, userID, 2025)
	if err != nil {
		t.Fatalf("GetPoopsByHour() error = %v", err)
	}

	// Should have all 24 hours
	if len(result) != 24 {
		t.Errorf("GetPoopsByHour() returned %d hours, want 24", len(result))
	}

	// Each hour should have exactly 1 poop
	for i, hd := range result {
		if hd.Hour != i {
			t.Errorf("Hour at index %d = %d, want %d", i, hd.Hour, i)
		}
		if hd.PoopCount != 1 {
			t.Errorf("Hour %d poop count = %d, want 1", i, hd.PoopCount)
		}
	}
}

func TestGetPoopsByDayOfWeek_AllDaysPresent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // Jan 1, 2025 is a Wednesday

	// Insert one poop on each day of the week
	userID := int64(3000)
	daysToAdd := []int{0, 1, 2, 3, 4, 5, 6} // Wed, Thu, Fri, Sat, Sun, Mon, Tue
	for i, dayOffset := range daysToAdd {
		timestamp := baseTime.AddDate(0, 0, dayOffset)
		err := LogPoop(ctx, db, userID, "testuser", int64(20000+i), timestamp.Format("2006-01-02 15:04:05"), timestamp.Unix())
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}

	result, err := GetPoopsByDayOfWeek(ctx, db, userID, 2025)
	if err != nil {
		t.Fatalf("GetPoopsByDayOfWeek() error = %v", err)
	}

	// Should have all 7 days
	if len(result) != 7 {
		t.Errorf("GetPoopsByDayOfWeek() returned %d days, want 7", len(result))
	}

	// Check that all days are present and in correct order
	expectedOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i, dotw := range result {
		if dotw.DayOfTheWeek != expectedOrder[i] {
			t.Errorf("Day at index %d = %s, want %s", i, dotw.DayOfTheWeek, expectedOrder[i])
		}
	}
}

// Benchmark tests
func BenchmarkGetYearlyPoopCount(b *testing.B) {
	db, cleanup := setupTestDB(&testing.T{})
	defer cleanup()
	insertTestData(&testing.T{}, db)

	ctx := context.Background()
	userID := int64(1001)
	year := 2025

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetYearlyPoopCount(ctx, db, userID, year)
	}
}

func BenchmarkGetPoopsByHour(b *testing.B) {
	db, cleanup := setupTestDB(&testing.T{})
	defer cleanup()
	insertTestData(&testing.T{}, db)

	ctx := context.Background()
	userID := int64(1001)
	year := 2025

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetPoopsByHour(ctx, db, userID, year)
	}
}

func BenchmarkGetGroupAwards(b *testing.B) {
	db, cleanup := setupTestDB(&testing.T{})
	defer cleanup()
	insertTestData(&testing.T{}, db)

	ctx := context.Background()
	year := 2025

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetGroupAwards(ctx, db, year)
	}
}
