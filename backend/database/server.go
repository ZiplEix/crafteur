package database

import "github.com/ZiplEix/crafteur/core"

func GetAllServers() ([]core.ServerConfig, error) {
	rows, err := DB.Query("SELECT id, name, type, port, ram, java_version FROM servers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []core.ServerConfig
	for rows.Next() {
		var s core.ServerConfig
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Port, &s.RAM, &s.JavaVersion); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func CreateServer(s *core.ServerConfig) error {
	_, err := DB.Exec(
		"INSERT INTO servers (id, name, type, port, ram, java_version) VALUES (?, ?, ?, ?, ?, ?)",
		s.ID, s.Name, s.Type, s.Port, s.RAM, s.JavaVersion,
	)
	return err
}

func DeleteServer(id string) error {
	_, err := DB.Exec("DELETE FROM servers WHERE id = ?", id)
	return err
}
