package repository

import (
	"github.com/mjiee/world-news/backend/repository/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// SystemConfigRepository is the interface for system config.
type SystemConfigRepository interface {
}

type systemConfigRepository struct {
	db *gorm.DB
}

func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &systemConfigRepository{db}
}

// GetSystemConfig get system config
func (s *systemConfigRepository) GetSystemConfig(key string) (*model.SystemConfig, error) {
	var systemConfig *model.SystemConfig

	if err := s.db.Where("key = ?", key).First(&systemConfig).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return systemConfig, nil
		}

		return nil, errors.WithStack(err)
	}

	return systemConfig, nil
}
