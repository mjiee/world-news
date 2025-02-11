package entity

import "time"

// SystemConfig represents the structure of system configuration data stored in the database.
type SystemConfig struct {
	Id        uint
	Key       string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}
