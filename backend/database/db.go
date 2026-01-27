package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", 0755)
	}

	var err error
	DB, err = sql.Open("sqlite", "./data/manager.db")
	if err != nil {
		log.Fatal("Erreur ouverture DB:", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS servers (
		id TEXT PRIMARY KEY,
		name TEXT,
		type TEXT,
		port INTEGER,
		ram INTEGER DEFAULT 2048,
		java_version INTEGER DEFAULT 21,
		version TEXT
	);
	
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password_hash TEXT
	);

	CREATE TABLE IF NOT EXISTS scheduled_tasks (
		id TEXT PRIMARY KEY,
		server_id TEXT,
		name TEXT,
		action TEXT,
		payload TEXT,
		cron_expression TEXT,
		one_shot BOOLEAN,
		last_run DATETIME
	);`

	if _, err := DB.Exec(query); err != nil {
		log.Fatal("Erreur création table:", err)
	}
}
