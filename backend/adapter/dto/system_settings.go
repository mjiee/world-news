package dto

import "github.com/mjiee/world-news/backend/pkg/httpx"

// SystemConfig system config
type SystemConfig struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
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
