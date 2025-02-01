package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"src/utils"

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

func get_db_path() string {
	if os.Getenv("DB_PATH") == "" {
		utils.LoadEnv()
	}
	return os.Getenv("DB_PATH")
}

func Log_Poop(db *sql.DB, userID int64, username string, msgId int, timestamp string) error {
	query := `
	INSERT INTO poop_tracker (user_id, username, message_id, timestamp)
	VALUES (?, ?, ?, ?)
	`
	log.Println("Logging poop for user:", username)
	_, err := db.Exec(query, userID, username, msgId, timestamp)
	return err
}

func Get_Global_Poop_Count(db *sql.DB, userID int64) (int, error) {
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

func Get_Monthly_Poop_Count(db *sql.DB, userID int64) (int, error) {
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

func Get_Monthly_Poodium(db *sql.DB) ([]UserPoopCount, error) {
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

func Get_Past_Month_Poodium(db *sql.DB) ([]UserPoopCount, error) {
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

func Get_Yearly_Poodium(db *sql.DB) ([]UserPoopCount, error) {
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

func Get_Monthly_Poop_Stats(db *sql.DB, userID int64) ([]MonthlyPoopCount, error) {
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

func Get_Monthly_Leaderboard(db *sql.DB) ([]UserPoopCount, error) {
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

func Get_Bottom_Poopers(db *sql.DB) ([]UserPoopCount, error) {
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

func Get_Days_Without_Poop(db *sql.DB, userID int64) (int, error) {
    query := `
    WITH all_days AS (
        SELECT julianday('now') - julianday(strftime('%Y-01-01', 'now')) + 1 AS day
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

func Get_Max_Poop_Streak(db *sql.DB, userID int64) (int, error) {
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

func Get_Day_With_Most_Poops(db *sql.DB, userID int64) (string, int, error) {
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

func create_table(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS poop_tracker (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    user_id INTEGER NOT NULL,
	    username TEXT NOT NULL,
		message_id INTEGER UNIQUE NOT NULL,
	    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	return err
}

func Open_DB_Connection() *sql.DB {

    dbPath := get_db_path()
    if dbPath == "" {
        log.Fatal("DB_PATH is not set")
    }

	db, err := sql.Open("sqlite", dbPath)
	utils.CheckError("Failed to connect to SQLite database", err)
	
	err = create_table(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Table created or already exists.")
	return db
}