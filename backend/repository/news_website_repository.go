package repository

import "gorm.io/gorm"

// NewsWebsiteRepository is the interface for news website.
type NewsWebsiteRepository interface {
}

type newsWebsiteRepository struct {
	db *gorm.DB
}

func NewNewsWebsiteRepository(db *gorm.DB) NewsWebsiteRepository {
	return &newsWebsiteRepository{db}
}
