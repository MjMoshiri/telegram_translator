package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db             *sql.DB
	userCache      = make(map[int64]bool)
	userCacheMutex sync.RWMutex
)

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./bot.db")
	if err != nil {
		log.Fatal(err)
	}

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		language TEXT DEFAULT 'es'
	);`

	createUsageTable := `CREATE TABLE IF NOT EXISTS daily_usage (
		date TEXT,
		user_id INTEGER,
		tokens INTEGER,
		PRIMARY KEY (date, user_id)
	);`

	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(createUsageTable); err != nil {
		log.Fatal(err)
	}

	// Load users into cache
	rows, err := db.Query("SELECT id FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	userCacheMutex.Lock()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		userCache[id] = true
	}
	userCacheMutex.Unlock()

	if err := AddUser(OwnerID); err != nil {
		log.Println("Error adding owner to database:", err)
	}
}

func AddUser(id int64) error {
	_, err := db.Exec("INSERT OR IGNORE INTO users (id) VALUES (?)", id)
	if err != nil {
		return err
	}
	userCacheMutex.Lock()
	userCache[id] = true
	userCacheMutex.Unlock()
	return nil
}

func RemoveUser(id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	userCacheMutex.Lock()
	delete(userCache, id)
	userCacheMutex.Unlock()
	return nil
}

func IsUserAllowed(id int64) bool {
	userCacheMutex.RLock()
	allowed := userCache[id]
	userCacheMutex.RUnlock()
	return allowed
}

func SetUserLanguage(id int64, lang string) error {
	_, err := db.Exec("UPDATE users SET language = ? WHERE id = ?", lang, id)
	return err
}

func GetUserLanguage(id int64) string {
	var lang string
	err := db.QueryRow("SELECT language FROM users WHERE id = ?", id).Scan(&lang)
	if err != nil {
		return "es" // Default
	}
	return lang
}

func UpdateDailyUsage(date string, userID int64, tokens int) error {
	_, err := db.Exec(`
		INSERT INTO daily_usage (date, user_id, tokens) 
		VALUES (?, ?, ?) 
		ON CONFLICT(date, user_id) 
		DO UPDATE SET tokens = tokens + ?`,
		date, userID, tokens, tokens)
	return err
}
