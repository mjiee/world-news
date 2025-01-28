package service

import (
	"github.com/mjiee/world-news/backend/repository"
	"gorm.io/gorm"
)

type SystemSettingsService interface {
}

type systemSettingsService struct {
	newsWebsiteRepo repository.NewsWebsiteRepository
	keywordsRepo    repository.NewsKeywordsRepository
}

func NewSystemSettingsService(db *gorm.DB) SystemSettingsService {
	return &systemSettingsService{
		newsWebsiteRepo: repository.NewNewsWebsiteRepository(db),
		keywordsRepo:    repository.NewNewsKeywordsRepository(db),
	}
}
