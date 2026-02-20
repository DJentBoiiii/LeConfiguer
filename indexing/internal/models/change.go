package models

import "time"

type ConfigChange struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ConfigID    string    `json:"config_id" gorm:"index;not null"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Environment string    `json:"environment"`
	Action      string    `json:"action" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}
