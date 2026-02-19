package models

// Config represents a configuration item.
type Config struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"` // e.g., "yaml", "json", "toml"
	Environment string      `json:"environment"`
	JSONContent interface{} `json:"json_content"`
	Tags        []string    `json:"tags"`
}
