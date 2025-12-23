package repository

import (
	"database/sql"
	"fmt"
	"log"

	"src/config"

	_ "modernc.org/sqlite"
)

type MonthlyPoopCount struct {
	Month     string
	PoopCount int
}

type UserPoopCount struct {
	Username  string
	PoopCount int
}

func LogPoop(db *sql.DB, userID int64, username string, msgId int, timestamp string, unixTimestamp int64) error {
	query := `
	INSERT INTO poop_tracker (user_id, username, message_id, timestamp, created_at_unix)
	VALUES (?, ?, ?, ?, ?)
	`
	log.Println("Logging poop for user:", username)
	_, err := db.Exec(query, userID, username, msgId, timestamp, unixTimestamp)
	return err
}

func GetGlobalPoopCount(db *sql.DB, userID int64) (int, error) {
	query := `
    SELECT COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE user_id = ?;
    `
	var poopCount int
	err := db.QueryRow(query, userID).Scan(&poopCount)
	if err != nil {
		return 0, err
	}
	return poopCount, nil
}

func GetMonthlyPoopCount(db *sql.DB, userID int64) (int, error) {
	query := `
    SELECT COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE user_id = $1 AND strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now');
    `
	var poopCount int
	err := db.QueryRow(query, userID).Scan(&poopCount)
	if err != nil {
		return 0, err
	}
	return poopCount, nil
}

func GetMonthlyPoodium(db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.Query(query)
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

func GetPastMonthPoodium(db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now', '-1 month')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.Query(query)
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

func GetYearlyPoodium(db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y', timestamp) = strftime('%Y', 'now')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.Query(query)
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

func GetMonthlyPoopStats(db *sql.DB, userID int64) ([]MonthlyPoopCount, error) {
	query := `
    SELECT strftime('%Y-%m', timestamp) AS month, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE user_id = ?
    GROUP BY month
    ORDER BY month;
    `

	rows, err := db.Query(query, userID)
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

func GetMonthlyLeaderboard(db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now')
    GROUP BY user_id
    ORDER BY poop_count DESC, MAX(timestamp) ASC;
    `
	rows, err := db.Query(query)
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

func GetBottomPoopers(db *sql.DB) ([]UserPoopCount, error) {
	query := `
    SELECT username, COUNT(*) AS poop_count
    FROM poop_tracker
    WHERE strftime('%Y-%m', timestamp) = strftime('%Y-%m', 'now')
    GROUP BY user_id
    ORDER BY poop_count ASC, MAX(timestamp) ASC
    LIMIT 3;
    `
	rows, err := db.Query(query)
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

func GetDaysWithoutPoop(db *sql.DB, userID int64) (int, error) {
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
	err := db.QueryRow(query, userID).Scan(&daysWithoutPoop)

	fmt.Println(daysWithoutPoop)
	if err != nil {
		return 0, err
	}
	return daysWithoutPoop, nil
}

func GetMaxPoopStreak(db *sql.DB, userID int64) (int, error) {
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
	err := db.QueryRow(query, userID).Scan(&maxStreak)
	if err != nil {
		return 0, err
	}
	return maxStreak, nil
}

func GetDayWithMostPoops(db *sql.DB, userID int64) (string, int, error) {
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
	err := db.QueryRow(query, userID).Scan(&day, &dumps)
	if err != nil {
		return "", 0, err
	}
	return day, dumps, nil
}

func createTable(db *sql.DB) error {
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
	_, err := db.Exec(query)
	return err
}

func OpenDBConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	err = createTable(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Table created or already exists.")
	return db, nil
}
