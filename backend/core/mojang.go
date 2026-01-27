package core

type MojangManifest struct {
	Versions []MojangVersion `json:"versions"`
}

type MojangVersion struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "release" ou "snapshot"
	URL  string `json:"url"`  // URL vers le json de détail
}

type VersionPackage struct {
	Downloads struct {
		Server struct {
			URL  string `json:"url"`
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
		} `json:"server"`
	} `json:"downloads"`
}
