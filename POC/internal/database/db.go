package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type ServerModel struct {
	ID   int
	Name string
	Type string // Vanilla or fabric
	Port int
}

func Init() {
	var err error
	DB, err = sql.Open("sqlite", "./data/crafteur.db")
	if err != nil {
		log.Fatal("Cannot open database: ", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		type TEXT,
		port INTEGER
	);`
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateServer(name string, serverType string, port int) (int64, error) {
	res, err := DB.Exec("INSERT INTO servers (name, type, port) VALUES (?, ?, ?)", name, serverType, port)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetAllServers() ([]ServerModel, error) {
	rows, err := DB.Query("SELECT id, name, type, port FROM servers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []ServerModel
	for rows.Next() {
		var s ServerModel
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Port); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}
