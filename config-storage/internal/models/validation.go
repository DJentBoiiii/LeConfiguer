package models

import (
	"errors"
	"strings"
)

var (
	ErrEmptyID          = errors.New("id cannot be empty")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrEmptyType        = errors.New("type cannot be empty")
	ErrInvalidType      = errors.New("type must be one of: yaml, json, toml, xml, ini")
	ErrEmptyEnvironment = errors.New("environment cannot be empty")
	ErrEmptyContent     = errors.New("json_content cannot be empty")
)

var validTypes = map[string]bool{
	"yaml": true,
	"json": true,
	"toml": true,
	"xml":  true,
	"ini":  true,
}

// ValidateForCreate checks if the Config has valid fields for creation (ID is optional).
func (c *Config) ValidateForCreate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrEmptyName
	}

	if strings.TrimSpace(c.Type) == "" {
		return ErrEmptyType
	}

	if !validTypes[strings.ToLower(c.Type)] {
		return ErrInvalidType
	}

	if strings.TrimSpace(c.Environment) == "" {
		return ErrEmptyEnvironment
	}

	if c.JSONContent == nil {
		return ErrEmptyContent
	}

	return nil
}

// Validate checks if the Config has valid fields (including ID).
func (c *Config) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return ErrEmptyID
	}

	return c.ValidateForCreate()
}
