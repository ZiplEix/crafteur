package core

import "time"

type ModrinthSearchResponse struct {
	Hits      []ModrinthProject `json:"hits"`
	Offset    int               `json:"offset"`
	Limit     int               `json:"limit"`
	TotalHits int               `json:"total_hits"`
}

type ModrinthProject struct {
	ProjectID     string   `json:"project_id"`
	Slug          string   `json:"slug"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Categories    []string `json:"categories"` // e.g., ["fabric", "forge"]
	IconURL       string   `json:"icon_url"`
	ClientSide    string   `json:"client_side"`
	ServerSide    string   `json:"server_side"`
	Downloads     int      `json:"downloads"`
	Follows       int      `json:"follows"`
	ProjectType   string   `json:"project_type"` // mod, modpack, etc.
	Author        string   `json:"author"`
	Versions      []string `json:"versions"`
	LatestVersion string   `json:"latest_version,omitempty"`
	License       string   `json:"license"`
	Gallery       []string `json:"gallery"`
	DateCreated   string   `json:"date_created"`
	DateModified  string   `json:"date_modified"`
}

type ModrinthVersion struct {
	ID            string         `json:"id"`
	ProjectID     string         `json:"project_id"`
	AuthorID      string         `json:"author_id"`
	Featured      bool           `json:"featured"`
	Name          string         `json:"name"`
	VersionNumber string         `json:"version_number"`
	Changelog     string         `json:"changelog"`
	DatePublished time.Time      `json:"date_published"`
	Downloads     int            `json:"downloads"`
	VersionType   string         `json:"version_type"` // release, beta, alpha
	Files         []ModrinthFile `json:"files"`
	Dependencies  []struct {
		VersionID string `json:"version_id"`
		ProjectID string `json:"project_id"`
	} `json:"dependencies"`
	GameVersions []string `json:"game_versions"`
	Loaders      []string `json:"loaders"`
}

type ModrinthFile struct {
	Hashes   map[string]string `json:"hashes"`
	Url      string            `json:"url"`
	Filename string            `json:"filename"`
	Primary  bool              `json:"primary"`
	Size     int               `json:"size"`
}
