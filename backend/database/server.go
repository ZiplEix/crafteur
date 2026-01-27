package database

import (
	"database/sql"

	"github.com/ZiplEix/crafteur/core"
)

func GetAllServers() ([]core.ServerConfig, error) {
	rows, err := DB.Query("SELECT id, name, type, port, ram, java_version, version, jar_name FROM servers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []core.ServerConfig
	for rows.Next() {
		var s core.ServerConfig
		var jarName sql.NullString // Handle potential nulls safely for old rows if migration missed (though default takes care)
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Port, &s.RAM, &s.JavaVersion, &s.Version, &jarName); err != nil {
			return nil, err
		}
		if jarName.Valid {
			s.JarName = jarName.String
		} else {
			s.JarName = "server.jar"
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func GetServer(id string) (*core.ServerConfig, error) {
	var s core.ServerConfig
	var jarName sql.NullString
	err := DB.QueryRow("SELECT id, name, type, port, ram, java_version, version, jar_name FROM servers WHERE id = ?", id).Scan(&s.ID, &s.Name, &s.Type, &s.Port, &s.RAM, &s.JavaVersion, &s.Version, &jarName)
	if err != nil {
		return nil, err
	}
	if jarName.Valid {
		s.JarName = jarName.String
	} else {
		s.JarName = "server.jar"
	}
	return &s, nil
}

func CreateServer(s *core.ServerConfig) error {
	_, err := DB.Exec(
		"INSERT INTO servers (id, name, type, port, ram, java_version, version, jar_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		s.ID, s.Name, s.Type, s.Port, s.RAM, s.JavaVersion, s.Version, s.JarName,
	)
	return err
}

func DeleteServer(id string) error {
	_, err := DB.Exec("DELETE FROM servers WHERE id = ?", id)
	return err
}
