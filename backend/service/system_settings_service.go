package service

import (
	"github.com/mjiee/world-news/backend/repository"
	"gorm.io/gorm"
)

type SystemSettingsService interface {
}

type systemSettingsService struct {
	systemConfigRepo repository.SystemConfigRepository
}

func NewSystemSettingsService(db *gorm.DB) SystemSettingsService {
	return &systemSettingsService{
		systemConfigRepo: repository.NewSystemConfigRepository(db),
	}
}
