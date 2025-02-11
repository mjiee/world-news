package entity

import (
	"time"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/repository/model"
)

// SystemConfig represents the structure of system configuration data stored in the database.
type SystemConfig struct {
	Id        uint
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSystemConfig creates a new SystemConfig entity.
func NewSystemConfig(key, value string) *SystemConfig {
	return &SystemConfig{
		Key:   key,
		Value: value,
	}
}

// NewSystemConfigFromModel converts a SystemConfigModel to a SystemConfig entity.
func NewSystemConfigFromModel(m *model.SystemConfig) (*SystemConfig, error) {
	if m == nil {
		return nil, errorx.SystemConfigNotFound
	}

	return &SystemConfig{
		Id:        m.Id,
		Key:       m.Key,
		Value:     m.Value,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

// ToModel converts the SystemConfig entity to a SystemConfigModel.
func (s *SystemConfig) ToModel() (*model.SystemConfig, error) {
	if s == nil {
		return nil, errorx.SystemConfigNotFound
	}

	return &model.SystemConfig{
		Id:        s.Id,
		Key:       s.Key,
		Value:     s.Value,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}, nil
}
