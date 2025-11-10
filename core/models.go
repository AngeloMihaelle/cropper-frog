package core

// Clip defines the structure for a single crop interval (Req 3.2)
// This is the shared data model between the Go backend and Svelte frontend.
type Clip struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
