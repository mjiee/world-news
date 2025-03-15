package translate

import (
	"context"
	"fmt"

	"github.com/mjiee/world-news/backend/pkg/errorx"
)

// Translater is the interface for the translator
type Translater interface {
	Translate(ctx context.Context, toLang string, texts ...string) ([]string, error)
}

// platform is the platform for the translator
type platform string

const (
	aliyunPlatform    platform = "aliyun"
	baiduPlatform     platform = "baidu"
	googlePlatform    platform = "google"
	microsoftPlatform platform = "microsoft"
)

// Config is the config for the translator
type Config struct {
	Platform  platform `json:"platform"`
	AppId     string   `json:"appId"`
	AppSecret string   `json:"appSecret"`
}

// NewTranslater creates a new translator
func NewTranslater(cfg *Config) (Translater, error) {
	switch cfg.Platform {
	case baiduPlatform:
		return newBaiduTranslater(cfg.AppId, cfg.AppSecret)
	case googlePlatform:
		return newGoogleTranslater(cfg.AppId)
	case microsoftPlatform:
		return newMicrosoftTranslater(cfg.AppId)
	case aliyunPlatform:
		return newAliyunTranslater(cfg.AppId, cfg.AppSecret)
	default:
		return nil, errorx.InternalError.SetMessage(fmt.Sprintf("unsupported translation platform: %s", cfg.Platform))
	}
}
