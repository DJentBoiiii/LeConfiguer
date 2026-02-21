package models

import "errors"

// Config represents a configuration item.
type Config struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"` // e.g., "yaml", "json", "toml"
	Environment string      `json:"environment"`
	JSONContent interface{} `json:"json_content"`
}

// Validate validates a config for update.
func (c *Config) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Type == "" {
		return errors.New("type is required")
	}
	if c.Environment == "" {
		return errors.New("environment is required")
	}
	return nil
}
