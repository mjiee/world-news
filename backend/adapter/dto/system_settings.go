package dto

import (
	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/pkg/httpx"
)

// SystemConfig system config
type SystemConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewSystemConfigFromEntity creates a new SystemConfig from the provided entity.
func NewSystemConfigFromEntity(config *entity.SystemConfig) *SystemConfig {
	return &SystemConfig{
		Key:   config.Key,
		Value: config.Value,
	}
}

// ToEntity converts the SystemConfig to an entity.
func (s *SystemConfig) ToEntity() *entity.SystemConfig {
	return &entity.SystemConfig{
		Key:   s.Key,
		Value: s.Value,
	}
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
