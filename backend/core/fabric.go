package core

type FabricInstaller struct {
	Url    string `json:"url"`
	Maven  string `json:"maven"`
	Stable bool   `json:"stable"`
}
