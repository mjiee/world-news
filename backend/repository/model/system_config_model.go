package model

import "time"

// SystemConfig represents the structure of system configuration data stored in the database.
type SystemConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"unique"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (n *SystemConfig) TableName() string {
	return "system_configs"
}
