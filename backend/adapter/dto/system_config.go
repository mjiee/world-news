package dto

import (
	"encoding/json"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/pkg/httpx"
)

// SystemConfig system config
type SystemConfig struct {
	Key   string `json:"key" binding:"required"`
	Value any    `json:"value" binding:"required"`
}

// NewSystemConfigFromEntity creates a new SystemConfig from the provided entity.
func NewSystemConfigFromEntity(config *entity.SystemConfig) *SystemConfig {
	if config == nil {
		return nil
	}

	var value any
	if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
		return nil
	}

	return &SystemConfig{
		Key:   config.Key.String(),
		Value: value,
	}
}

// ToEntity converts the SystemConfig to an entity.
func (s *SystemConfig) ToEntity() (*entity.SystemConfig, error) {
	return entity.NewSystemConfig(s.Key, s.Value)
}

// GetSystemConfigRequest get system config request
type GetSystemConfigRequest struct {
	Key string `json:"key"`
}

// GetSystemConfigResponse get system config response
type GetSystemConfigResponse struct {
	*httpx.Response
	Result *SystemConfig `json:"result"`
}
