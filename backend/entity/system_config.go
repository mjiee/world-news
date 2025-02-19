package entity

import (
	"encoding/json"
	"time"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/locale"
	"github.com/mjiee/world-news/backend/repository/model"

	"github.com/pkg/errors"
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

	var value any

	if err := json.Unmarshal([]byte(m.Value), &value); err != nil {
		return nil, errors.WithStack(err)
	}

	return &SystemConfig{
		Id:        m.Id,
		Key:       valueobject.SystemConfigKey(m.Key),
		Value:     value,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

// ToModel converts the SystemConfig entity to a SystemConfigModel.
func (s *SystemConfig) ToModel() (*model.SystemConfig, error) {
	if s == nil {
		return nil, errorx.SystemConfigNotFound
	}

	value, err := json.Marshal(s.Value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &model.SystemConfig{
		Id:        s.Id,
		Key:       s.Key.String(),
		Value:     string(value),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
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
