package entity

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/locale"
	"github.com/mjiee/world-news/backend/repository/model"
)

// SystemConfig represents the structure of system configuration data stored in the database.
type SystemConfig struct {
	Id        uint
	Key       valueobject.SystemConfigKey
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSystemConfig creates a new SystemConfig entity.
func NewSystemConfig(key string, value any) (*SystemConfig, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &SystemConfig{Key: valueobject.SystemConfigKey(key), Value: string(data)}, err
}

// NewSystemConfigFromModel converts a SystemConfigModel to a SystemConfig entity.
func NewSystemConfigFromModel(m *model.SystemConfig) (*SystemConfig, error) {
	if m == nil {
		return nil, errorx.SystemConfigNotFound
	}

	s := &SystemConfig{
		Id:        m.ID,
		Key:       valueobject.SystemConfigKey(m.Key),
		Value:     m.Value,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return s, nil
}

// ToModel converts the SystemConfig entity to a SystemConfigModel.
func (s *SystemConfig) ToModel() (*model.SystemConfig, error) {
	if s == nil {
		return nil, errorx.SystemConfigNotFound
	}

	return &model.SystemConfig{
		ID:        s.Id,
		Key:       s.Key.String(),
		Value:     s.Value,
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
		var lang string
		if err := s.UnmarshalValue(&lang); err != nil {
			return err
		}

		return locale.SetAppLocalizer(lang)
	}

	return nil
}

// UnmarshalValue unmarshal config value
func (s *SystemConfig) UnmarshalValue(v any) error {
	err := json.Unmarshal([]byte(s.Value), v)

	return errors.WithStack(err)
}

// UnmarshalValue unmarshal config value
func UnmarshalValue[T any](data *SystemConfig, notErr error) (*T, error) {
	if data == nil || data.Id == 0 {
		return nil, notErr
	}

	var v T

	err := json.Unmarshal([]byte(data.Value), &v)

	return &v, errors.WithStack(err)
}
