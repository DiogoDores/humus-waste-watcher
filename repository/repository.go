package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"src/config"

	_ "modernc.org/sqlite"
)

type HourDistribution struct {
	Hour      int
	PoopCount int
}

type DayOfWeekDistribution struct {
	DayOfTheWeek string
	PoopCount    int
}

type PoopPersonality struct {
	Type        string
	Description string
	Confidence  float64
}

type GroupAward struct {
	AwardName string
	Winner    string
	Value     string
	Emoji     string
}

type YearlyRanking struct {
	Rank       int
	TotalUsers int
	Percentage float64 // percentage of users below this user
}

type MonthlyPoopCount struct {
	Month     string
	PoopCount int
}

type UserPoopCount struct {
	Username  string
	PoopCount int
}

func LogPoop(ctx context.Context, db *sql.DB, userID int64, username string, msgId int64, timestamp string, unixTimestamp int64) error {
	query := `
	INSERT INTO poop_tracker (user_id, username, message_id, timestamp, created_at_unix)
	VALUES (?, ?, ?, ?, ?)
	`
	log.Println("Logging poop for user:", username)
	_, err := db.ExecContext(ctx, query, userID, username, msgId, timestamp, unixTimestamp)
	return err
}

func GetGlobalPoopCount(ctx context.Context, db *sql.DB, userID int64) (int, error) {
	query := `
    SELECT COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE user_id = ?;
    `
	var poopCount int
	err := db.QueryRowContext(ctx, query, userID).Scan(&poopCount)
	if err != nil {
		return 0, err
	}
	return poopCount, nil
}

func GetMonthlyPoopCount(ctx context.Context, db *sql.DB, userID int64) (int, error) {
	query := `
    SELECT COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE user_id = $1 AND strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now');
    `
	var poopCount int
	err := db.QueryRowContext(ctx, query, userID).Scan(&poopCount)
	if err != nil {
		return 0, err
	}
	return poopCount, nil
}

func GetMonthlyPoodium(ctx context.Context, db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topPoopers []UserPoopCount
	for rows.Next() {
		var upc UserPoopCount
		if err := rows.Scan(&upc.Username, &upc.PoopCount); err != nil {
			return nil, err
		}
		topPoopers = append(topPoopers, upc)
	}
	return topPoopers, nil
}

func GetPastMonthPoodium(ctx context.Context, db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now', '-1 month')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topPoopers []UserPoopCount
	for rows.Next() {
		var upc UserPoopCount
		if err := rows.Scan(&upc.Username, &upc.PoopCount); err != nil {
			return nil, err
		}
		topPoopers = append(topPoopers, upc)
	}
	return topPoopers, nil
}

func GetYearlyPoodium(ctx context.Context, db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y', timestamp) = strftime('%Y', 'now')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topPoopers []UserPoopCount
	for rows.Next() {
		var upc UserPoopCount
		if err := rows.Scan(&upc.Username, &upc.PoopCount); err != nil {
			return nil, err
		}
		topPoopers = append(topPoopers, upc)
	}
	return topPoopers, nil
}

func GetMonthlyPoopStats(ctx context.Context, db *sql.DB, userID int64) ([]MonthlyPoopCount, error) {
	query := `
    SELECT strftime('%Y-%m', timestamp) AS month, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE user_id = ?
    GROUP BY month
    ORDER BY month;
    `

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MonthlyPoopCount
	for rows.Next() {
		var mpc MonthlyPoopCount
		if err := rows.Scan(&mpc.Month, &mpc.PoopCount); err != nil {
			return nil, err
		}
		results = append(results, mpc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetMonthlyLeaderboard(ctx context.Context, db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC;
    `
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []UserPoopCount
	for rows.Next() {
		var upc UserPoopCount
		if err := rows.Scan(&upc.Username, &upc.PoopCount); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, upc)
	}
	return leaderboard, nil
}

func GetBottomPoopers(ctx context.Context, db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now')
    GROUP BY user_id
    ORDER BY poop_count ASC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bottomPoopers []UserPoopCount
	for rows.Next() {
		var upc UserPoopCount
		if err := rows.Scan(&upc.Username, &upc.PoopCount); err != nil {
			return nil, err
		}
		bottomPoopers = append(bottomPoopers, upc)
	}
	return bottomPoopers, nil
}

func GetDaysWithoutPoop(ctx context.Context, db *sql.DB, userID int64) (int, error) {
	query := `
    WITH all_days AS (
        SELECT date('now', '-' || (julianday('now') - julianday(date('now', 'start of year'))) || ' days') AS day
        UNION ALL
        SELECT date(day, '+1 day')
        FROM all_days
        WHERE day < date('now')
    ),
    pooped_days AS (
        SELECT DISTINCT date(timestamp) AS day
        FROM poop_tracker
        WHERE user_id = ?
    )
    SELECT COUNT(*)
    FROM all_days
    WHERE day NOT IN (SELECT day FROM pooped_days);
    `
	var daysWithoutPoop int
	err := db.QueryRowContext(ctx, query, userID).Scan(&daysWithoutPoop)

	if err != nil {
		return 0, err
	}
	return daysWithoutPoop, nil
}

func GetMaxPoopStreak(ctx context.Context, db *sql.DB, userID int64) (int, error) {
	query := `
    WITH daily_poops AS (
        SELECT 
            date(timestamp) AS day, 
            COUNT(*) AS poops
        FROM poop_tracker
        WHERE user_id = ?
        GROUP BY day
        ORDER BY day
    ),
    streaks AS (
        SELECT 
            day, 
            poops,
            -- Identify streak groups
            ROW_NUMBER() OVER (ORDER BY day) - 
            ROW_NUMBER() OVER (PARTITION BY poops ORDER BY day) AS streak_group
        FROM daily_poops
    )
    SELECT MAX(streak_count) AS max_streak
    FROM (
        SELECT COUNT(*) AS streak_count
        FROM streaks
        GROUP BY streak_group
    );
    `
	var maxStreak int
	err := db.QueryRowContext(ctx, query, userID).Scan(&maxStreak)
	if err != nil {
		return 0, err
	}
	return maxStreak, nil
}

func GetDayWithMostPoops(ctx context.Context, db *sql.DB, userID int64) (string, int, error) {
	query := `
    SELECT 
        date(timestamp) AS day, 
        COUNT(*) AS dumps
    FROM poop_tracker
    WHERE user_id = ?
    GROUP BY day
    ORDER BY dumps DESC
    LIMIT 1;
    `
	var day string
	var dumps int
	err := db.QueryRowContext(ctx, query, userID).Scan(&day, &dumps)
	if err != nil {
		return "", 0, err
	}
	return day, dumps, nil
}

func GetYearlyPoopCount(ctx context.Context, db *sql.DB, userID int64, year int) (int, error) {
	query := `
	SELECT COUNT(*) AS poop_count
	FROM poop_tracker
	WHERE user_id = ? AND strftime('%Y', timestamp) = ?;
	`
	var poopCount int
	err := db.QueryRowContext(ctx, query, userID, strconv.Itoa(year)).Scan(&poopCount)
	if err != nil {
		return 0, err
	}
	return poopCount, nil
}

func GetPoopsByHour(ctx context.Context, db *sql.DB, userID int64, year int) ([]HourDistribution, error) {
	query := `
	SELECT CAST(strftime('%H', timestamp) AS INTEGER) AS hour, COUNT(*) AS poop_count
	FROM poop_tracker
	WHERE user_id = ? AND strftime('%Y', timestamp) = ?
	GROUP BY hour
	ORDER BY hour;
	`

	rows, err := db.QueryContext(ctx, query, userID, strconv.Itoa(year))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hourMap := make(map[int]int)
	for i := 0; i < 24; i++ {
		hourMap[i] = 0
	}

	for rows.Next() {
		var hd HourDistribution
		if err := rows.Scan(&hd.Hour, &hd.PoopCount); err != nil {
			return nil, err
		}
		hourMap[hd.Hour] = hd.PoopCount
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var results []HourDistribution
	for hour := 0; hour < 24; hour++ {
		results = append(results, HourDistribution{
			Hour:      hour,
			PoopCount: hourMap[hour],
		})
	}

	return results, nil
}

func GetPoopsByDayOfWeek(ctx context.Context, db *sql.DB, userID int64, year int) ([]DayOfWeekDistribution, error) {
	query := `
	SELECT 
    CASE CAST(strftime('%w', timestamp) AS INTEGER)
        WHEN 0 THEN 'Sunday'
        WHEN 1 THEN 'Monday'
        WHEN 2 THEN 'Tuesday'
        WHEN 3 THEN 'Wednesday'
        WHEN 4 THEN 'Thursday'
        WHEN 5 THEN 'Friday'
        WHEN 6 THEN 'Saturday'
    END AS day_of_week,
    COUNT(*) AS poop_count
	FROM poop_tracker
	WHERE user_id = ? AND strftime('%Y', timestamp) = ?
	GROUP BY day_of_week
	ORDER BY 
    CASE CAST(strftime('%w', timestamp) AS INTEGER)
        WHEN 0 THEN 7
        WHEN 1 THEN 1
        WHEN 2 THEN 2
        WHEN 3 THEN 3
        WHEN 4 THEN 4
        WHEN 5 THEN 5
        WHEN 6 THEN 6
    END;
	`

	rows, err := db.QueryContext(ctx, query, userID, strconv.Itoa(year))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dayMap := make(map[string]int)
	dayOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range dayOrder {
		dayMap[day] = 0
	}

	for rows.Next() {
		var dotw DayOfWeekDistribution
		if err := rows.Scan(&dotw.DayOfTheWeek, &dotw.PoopCount); err != nil {
			return nil, err
		}
		dayMap[dotw.DayOfTheWeek] = dotw.PoopCount
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var results []DayOfWeekDistribution
	for _, day := range dayOrder {
		results = append(results, DayOfWeekDistribution{
			DayOfTheWeek: day,
			PoopCount:    dayMap[day],
		})
	}

	return results, nil
}

func GetYearlyRanking(ctx context.Context, db *sql.DB, userID int64, year int) (YearlyRanking, error) {
	query := `
	WITH user_stats AS (
		SELECT user_id, COUNT(*) AS poop_count
		FROM poop_tracker
		WHERE strftime('%Y', timestamp) = ?
		GROUP BY user_id	
	),
	user_rank AS (
		SELECT 
			user_id,
			poop_count,
			ROW_NUMBER() OVER (ORDER BY poop_count DESC) AS rank,
			COUNT(*) OVER () AS total_users
		FROM user_stats
	)
	SELECT rank, total_users,
		CAST ((total_users - rank) AS FLOAT) / CAST (total_users AS FLOAT) * 100.0 AS percentage
	FROM user_rank
	WHERE user_id = ?;
	`

	var yr YearlyRanking
	err := db.QueryRowContext(ctx, query, strconv.Itoa(year), userID).Scan(&yr.Rank, &yr.TotalUsers, &yr.Percentage)
	if err != nil {
		return YearlyRanking{}, err
	}
	return yr, nil
}

func GetGroupYearlyStats(ctx context.Context, db *sql.DB, year int) ([]UserPoopCount, error) {
	query := `
	SELECT username, COUNT(*) AS poop_count
	FROM poop_tracker
	WHERE strftime('%Y', timestamp) = ?
	GROUP BY user_id
	ORDER BY poop_count DESC;
	`
	rows, err := db.QueryContext(ctx, query, strconv.Itoa(year))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []UserPoopCount
	for rows.Next() {
		var upc UserPoopCount
		if err := rows.Scan(&upc.Username, &upc.PoopCount); err != nil {
			return nil, err
		}
		results = append(results, upc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetGroupAwards(ctx context.Context, db *sql.DB, year int) ([]GroupAward, error) {
	var awards []GroupAward
	yearStr := strconv.Itoa(year)

	// Award 1: Early Bird (Most poops 05:00-08:00)
	earlyBirdQuery := `
	SELECT username, COUNT(*) AS count
	FROM poop_tracker
	WHERE strftime('%Y', timestamp) = ?
	  AND CAST(strftime('%H', timestamp) AS INTEGER) BETWEEN 5 AND 8
	GROUP BY user_id
	ORDER BY count DESC
	LIMIT 1;
	`
	var earlyBirdWinner string
	var earlyBirdCount int
	err := db.QueryRowContext(ctx, earlyBirdQuery, yearStr).Scan(&earlyBirdWinner, &earlyBirdCount)
	if err == nil && earlyBirdCount > 0 {
		awards = append(awards, GroupAward{
			AwardName: "Early Bird",
			Winner:    earlyBirdWinner,
			Value:     strconv.Itoa(earlyBirdCount),
			Emoji:     "☀️",
		})
	}

	// Award 2: Night Owl (Most poops 23:00-04:00)
	nightOwlQuery := `
	SELECT username, COUNT(*) AS count
	FROM poop_tracker
	WHERE strftime('%Y', timestamp) = ?
	  AND (CAST(strftime('%H', timestamp) AS INTEGER) >= 23 
	       OR CAST(strftime('%H', timestamp) AS INTEGER) <= 4)
	GROUP BY user_id
	ORDER BY count DESC
	LIMIT 1;
	`
	var nightOwlWinner string
	var nightOwlCount int
	err = db.QueryRowContext(ctx, nightOwlQuery, yearStr).Scan(&nightOwlWinner, &nightOwlCount)
	if err == nil && nightOwlCount > 0 {
		awards = append(awards, GroupAward{
			AwardName: "Night Owl",
			Winner:    nightOwlWinner,
			Value:     strconv.Itoa(nightOwlCount),
			Emoji:     "🦉",
		})
	}

	// Award 3: Machine Gun (Most poops in single day)
	machineGunQuery := `
	SELECT username, MAX(daily_count) AS max_poops
	FROM (
		SELECT user_id, username, date(timestamp) AS day, COUNT(*) AS daily_count
		FROM poop_tracker
		WHERE strftime('%Y', timestamp) = ?
		GROUP BY user_id, day
	) AS daily_stats
	GROUP BY user_id
	ORDER BY max_poops DESC
	LIMIT 1;
	`
	var machineGunWinner string
	var machineGunCount int
	err = db.QueryRowContext(ctx, machineGunQuery, yearStr).Scan(&machineGunWinner, &machineGunCount)
	if err == nil && machineGunCount > 0 {
		awards = append(awards, GroupAward{
			AwardName: "Machine Gun",
			Winner:    machineGunWinner,
			Value:     strconv.Itoa(machineGunCount),
			Emoji:     "🔫",
		})
	}

	// Award 4: Consistency King (Longest streak)
	// Calculate max streak for each user by finding consecutive days
	consistencyQuery := `
	WITH daily_poops AS (
		SELECT 
			user_id,
			username,
			date(timestamp) AS day
		FROM poop_tracker
		WHERE strftime('%Y', timestamp) = ?
		GROUP BY user_id, day
	),
	streaks AS (
		SELECT 
			user_id,
			username,
			day,
			ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY day) - 
			julianday(day) AS streak_group
		FROM daily_poops
	),
	streak_lengths AS (
		SELECT 
			user_id,
			username,
			streak_group,
			COUNT(*) AS streak_count
		FROM streaks
		GROUP BY user_id, username, streak_group
	),
	max_streaks AS (
		SELECT 
			user_id,
			username,
			MAX(streak_count) AS max_streak
		FROM streak_lengths
		GROUP BY user_id, username
	)
	SELECT username, max_streak
	FROM max_streaks
	ORDER BY max_streak DESC
	LIMIT 1;
	`
	var consistencyWinner string
	var consistencyStreak int
	err = db.QueryRowContext(ctx, consistencyQuery, yearStr).Scan(&consistencyWinner, &consistencyStreak)
	if err == nil && consistencyStreak > 0 {
		awards = append(awards, GroupAward{
			AwardName: "Consistency King",
			Winner:    consistencyWinner,
			Value:     strconv.Itoa(consistencyStreak),
			Emoji:     "👑",
		})
	}

	// Award 5: Weekend Warrior (Highest % on Sat/Sun)
	weekendWarriorQuery := `
	WITH user_stats AS (
		SELECT 
			user_id,
			username,
			COUNT(*) AS total_poops,
			SUM(CASE WHEN CAST(strftime('%w', timestamp) AS INTEGER) IN (0, 6) THEN 1 ELSE 0 END) AS weekend_poops
		FROM poop_tracker
		WHERE strftime('%Y', timestamp) = ?
		GROUP BY user_id
	)
	SELECT username, CAST(weekend_poops AS FLOAT) / CAST(total_poops AS FLOAT) * 100.0 AS weekend_percentage
	FROM user_stats
	WHERE total_poops > 0
	ORDER BY weekend_percentage DESC
	LIMIT 1;
	`
	var weekendWinner string
	var weekendPercentage float64
	err = db.QueryRowContext(ctx, weekendWarriorQuery, yearStr).Scan(&weekendWinner, &weekendPercentage)
	if err == nil && weekendPercentage > 0 {
		awards = append(awards, GroupAward{
			AwardName: "Weekend Warrior",
			Winner:    weekendWinner,
			Value:     fmt.Sprintf("%.1f%%", weekendPercentage),
			Emoji:     "🎉",
		})
	}

	// Award 6: Boss makes a dollar, I make a dime (Most poops 09:00-18:00)
	companyTimeQuery := `
	SELECT username, COUNT(*) AS count
	FROM poop_tracker
	WHERE strftime('%Y', timestamp) = ?
	  AND CAST(strftime('%H', timestamp) AS INTEGER) BETWEEN 9 AND 18
	GROUP BY user_id
	ORDER BY count DESC
	LIMIT 1;
	`
	var companyTimeWinner string
	var companyTimeCount int
	err = db.QueryRowContext(ctx, companyTimeQuery, yearStr).Scan(&companyTimeWinner, &companyTimeCount)
	if err == nil && companyTimeCount > 0 {
		awards = append(awards, GroupAward{
			AwardName: "Boss makes a dollar, I make a dime",
			Winner:    companyTimeWinner,
			Value:     strconv.Itoa(companyTimeCount),
			Emoji:     "💰",
		})
	}

	return awards, nil
}

func createTable(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS poop_tracker (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    user_id INTEGER NOT NULL,
	    username TEXT NOT NULL,
		message_id INTEGER UNIQUE NOT NULL,
	    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        created_at_unix INTEGER NOT NULL
	);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

func OpenDBConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	ctx := context.Background()
	err = createTable(ctx, db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Table created or already exists.")
	return db, nil
}

func NewRepository(db *sql.DB) Repository {
	return NewSQLiteRepository(db)
}

func HealthCheck(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}
