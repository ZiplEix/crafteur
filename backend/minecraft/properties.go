package minecraft

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// LoadProperties reads a server.properties file and returns a map of key-values
func LoadProperties(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer file.Close()

	props := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
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

// SaveProperties writes the properties map to the file with a standard header
func SaveProperties(path string, props map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write standard header
	_, err = fmt.Fprintf(file, "#Minecraft server properties\n#%s\n", time.Now().Format("Mon Jan 02 15:04:05 MST 2006"))
	if err != nil {
		return err
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		_, err := fmt.Fprintf(file, "%s=%s\n", k, props[k])
		if err != nil {
			return err
		}
	}

	return nil
}
