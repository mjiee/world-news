package service

import (
	"context"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SystemSettingsService interface {
	GetSystemConfig(ctx context.Context, key string) (*entity.SystemConfig, error)
	SaveSystemConfig(ctx context.Context, config *entity.SystemConfig) error
}

type systemSettingsService struct {
}

func NewSystemSettingsService() SystemSettingsService {
	return &systemSettingsService{}
}

// GetSystemConfig retrieves the system configuration based on the provided key.
func (s *systemSettingsService) GetSystemConfig(ctx context.Context, key string) (*entity.SystemConfig, error) {
	repo := repository.Q.SystemConfig

	config, err := repo.WithContext(ctx).Where(repo.Key.Eq(key)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.SystemConfig{Key: key}, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.NewSystemConfigFromModel(config)
}

// SaveSystemConfig saves the provided system configuration.
func (s *systemSettingsService) SaveSystemConfig(ctx context.Context, config *entity.SystemConfig) error {
	data, err := config.ToModel()
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(repository.Q.SystemConfig.WithContext(ctx).Save(data))
}
