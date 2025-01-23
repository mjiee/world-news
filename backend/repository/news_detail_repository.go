package repository

import "gorm.io/gorm"

// NewsDetailRepository is the interface for news.
type NewsDetailRepository interface {
}

type newsDetailRepository struct {
	db *gorm.DB
}

func NewNewsDetailRepository(db *gorm.DB) NewsDetailRepository {
	return &newsDetailRepository{db}
}
