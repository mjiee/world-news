package entity

import (
	"time"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/locale"
	"github.com/mjiee/world-news/backend/repository/model"
)

// SystemConfig represents the structure of system configuration data stored in the database.
type SystemConfig struct {
	Id        uint
	Key       valueobject.SystemConfigKey
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSystemConfig creates a new SystemConfig entity.
func NewSystemConfig(key string, value any) *SystemConfig {
	return &SystemConfig{
		Key:   valueobject.SystemConfigKey(key),
		Value: value,
	}
}

// NewSystemConfigFromModel converts a SystemConfigModel to a SystemConfig entity.
func NewSystemConfigFromModel(m *model.SystemConfig) (*SystemConfig, error) {
	if m == nil {
		return nil, errorx.SystemConfigNotFound
	}

	s := &SystemConfig{
		Id:        m.ID,
		Key:       valueobject.SystemConfigKey(m.Key),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	value, err := s.Key.UnmarshalValue(m.Value)
	if err != nil {
		return nil, err
	}

	s.Value = value

	return s, nil
}

// ToModel converts the SystemConfig entity to a SystemConfigModel.
func (s *SystemConfig) ToModel() (*model.SystemConfig, error) {
	if s == nil {
		return nil, errorx.SystemConfigNotFound
	}

	value, err := s.Key.MarshalValue(s.Value)
	if err != nil {
		return nil, err
	}

	return &model.SystemConfig{
		ID:        s.Id,
		Key:       s.Key.String(),
		Value:     string(value),
		CreatedAt: s.CreatedAt,
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateSystemConfig updates the system configuration.
func (s *SystemConfig) UpdateSystemConfig() error {
	if s == nil {
		return errorx.SystemConfigNotFound
	}

	switch s.Key {
	case valueobject.LanguageKey:
		return locale.SetAppLocalizer(s.Value.(string))
	}

	return nil
}
