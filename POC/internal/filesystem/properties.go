package filesystem

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// LoadProperties reads the server.properties file and returns a Map
func LoadProperties(path string) (map[string]string, error) {
	props := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			props[key] = value
		}
	}
	return props, scanner.Err()
}

// SaveProperties writes the Map to the server.properties file
func SaveProperties(path string, props map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString("# Minecraft server properties\n")
	writer.WriteString("# Modified by Crafteur\n")

	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		line := fmt.Sprintf("%s=%s\n", k, props[k])
		writer.WriteString(line)
	}

	return writer.Flush()
}
