package repository

import (
	"src/utils"
	"log"
	"database/sql"
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

func Log_Poop(db *sql.DB, userID int64, username string, msgId int) error {
	query := `
	INSERT INTO poop_tracker (user_id, username, message_id)
	VALUES (?, ?, ?)
	`
	log.Println("Logging poop for user:", username)
	_, err := db.Exec(query, userID, username, msgId)
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
    ORDER BY poop_count DESC
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
    ORDER BY poop_count DESC
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
	db, err := sql.Open("sqlite", "./../../resources/poop_tracker.db")
	utils.CheckError("Failed to connect to SQLite database", err)
	
	err = create_table(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Table created or already exists.")
	return db
}