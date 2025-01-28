package repository

import "gorm.io/gorm"

type NewsKeywordsRepository interface {
}

type newsKeywordsRepository struct {
	db *gorm.DB
}

func NewNewsKeywordsRepository(db *gorm.DB) NewsKeywordsRepository {
	return &newsKeywordsRepository{
		db: db,
	}
}
