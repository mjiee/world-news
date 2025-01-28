package service

import (
	"github.com/mjiee/world-news/backend/repository"
	"gorm.io/gorm"
)

type NewsDetailService interface {
}

type newsDetailService struct {
	newsDetailRepo repository.NewsDetailRepository
}

func NewNewsDetailService(db *gorm.DB) NewsDetailService {
	return &newsDetailService{
		newsDetailRepo: repository.NewNewsDetailRepository(db),
	}
}
