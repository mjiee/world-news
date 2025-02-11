package service

import (
	"context"
	"encoding/json"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SystemSettingsService interface {
	SystemConfigInit(ctx context.Context) error
	GetSystemConfig(ctx context.Context, key string) (*entity.SystemConfig, error)
	SaveSystemConfig(ctx context.Context, config *entity.SystemConfig) error
}

type systemSettingsService struct {
}

func NewSystemSettingsService() SystemSettingsService {
	return &systemSettingsService{}
}

const (
	NewsWebsiteCollectionKey = "news_website_collection" // news website collection
	NewsWebsiteKey           = "news_website"            // news website
)

// SystemConfigInit initializes the system configuration.
func (s *systemSettingsService) SystemConfigInit(ctx context.Context) error {
	repo := repository.Q.SystemConfig

	total, err := repo.WithContext(ctx).Count()
	if err != nil {
		return errors.WithStack(err)
	}

	if total > 0 {
		return nil
	}

	collectionValue, _ := json.Marshal(valueobject.NewsWebsiteCollection)
	sysConfig, _ := entity.NewSystemConfig(NewsWebsiteCollectionKey, string(collectionValue)).ToModel()

	return errors.WithStack(repository.Q.SystemConfig.WithContext(ctx).Create(sysConfig))
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
	oldConfig, err := s.GetSystemConfig(ctx, config.Key)
	if err != nil {
		return errors.WithStack(err)
	}

	if oldConfig.Id > 0 {
		config.Id = oldConfig.Id
		config.CreatedAt = oldConfig.CreatedAt
	}

	data, err := config.ToModel()
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(repository.Q.SystemConfig.WithContext(ctx).Save(data))
}
