package minecraft

type PlayerCacheEntry struct {
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	ExpiresOn string `json:"expiresOn"`
}

type OpEntry struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type BanEntry struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Created string `json:"created"`
	Source  string `json:"source"`
	Expires string `json:"expires"`
	Reason  string `json:"reason"`
}
