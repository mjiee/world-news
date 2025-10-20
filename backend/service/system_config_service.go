package service

import (
	"context"
	"slices"

	"github.com/mjiee/gokit"
	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/urlx"
	"github.com/mjiee/world-news/backend/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// SystemConfigService system config service
type SystemConfigService interface {
	SystemConfigInit(ctx context.Context) error
	GetSystemConfig(ctx context.Context, key string) (*entity.SystemConfig, error)
	SaveSystemConfig(ctx context.Context, config *entity.SystemConfig) error
	DeleteSystemConfig(ctx context.Context, key string) error
	GetNewsWebsites(ctx context.Context) ([]*valueobject.NewsWebsite, error)
	UpdateNewsWebsiteWeight(ctx context.Context, website string, step int) error
	SaveNewsWebsites(ctx context.Context, newsWebsites []*valueobject.NewsWebsite) error
}

type systemConfigService struct {
}

func NewSystemConfigService() SystemConfigService {
	return &systemConfigService{}
}

// SystemConfigInit initializes the system configuration.
func (s *systemConfigService) SystemConfigInit(ctx context.Context) error {
	repo := repository.Q.SystemConfig

	// init language
	langConfig, err := s.GetSystemConfig(ctx, valueobject.LanguageKey.String())
	if err != nil {
		return errors.WithStack(err)
	}

	if langConfig.Id != 0 {
		if err := langConfig.UpdateSystemConfig(); err != nil {
			return errors.WithStack(err)
		}
	}

	// init news website collection
	total, err := repo.WithContext(ctx).Count()
	if err != nil {
		return errors.WithStack(err)
	}

	if total > 0 {
		return nil
	}

	sysConfig, err := entity.NewSystemConfig(valueobject.NewsWebsiteCollectionKey.String(), valueobject.NewsWebsiteCollection)
	if err != nil {
		return errors.WithStack(err)
	}

	sysConfigModel, err := sysConfig.ToModel()
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(repository.Q.SystemConfig.WithContext(ctx).Create(sysConfigModel))
}

// GetSystemConfig retrieves the system configuration based on the provided key.
func (s *systemConfigService) GetSystemConfig(ctx context.Context, key string) (*entity.SystemConfig, error) {
	repo := repository.Q.SystemConfig

	config, err := repo.WithContext(ctx).Where(repo.Key.Eq(key)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.NewSystemConfig(key, nil)
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.NewSystemConfigFromModel(config)
}

// SaveSystemConfig saves the provided system configuration.
func (s *systemConfigService) SaveSystemConfig(ctx context.Context, config *entity.SystemConfig) error {
	// get old config
	oldConfig, err := s.GetSystemConfig(ctx, config.Key.String())
	if err != nil {
		return errors.WithStack(err)
	}

	if oldConfig.Id > 0 {
		config.Id = oldConfig.Id
		config.CreatedAt = oldConfig.CreatedAt
	}

	// update system config
	if err = config.UpdateSystemConfig(); err != nil {
		return errors.WithStack(err)
	}

	// save config
	data, err := config.ToModel()
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(repository.Q.SystemConfig.WithContext(ctx).Save(data))
}

// DeleteSystemConfig deletes the system configuration based on the provided key.
func (s *systemConfigService) DeleteSystemConfig(ctx context.Context, key string) error {
	_, err := repository.Q.SystemConfig.WithContext(ctx).Where(repository.Q.SystemConfig.Key.Eq(key)).Delete()

	return errors.WithStack(err)
}

// GetNewsWebsites get the news websites.
func (s *systemConfigService) GetNewsWebsites(ctx context.Context) ([]*valueobject.NewsWebsite, error) {
	var newsWebsites []*valueobject.NewsWebsite

	// get news websites
	config, err := s.GetSystemConfig(ctx, valueobject.NewsWebsiteKey.String())
	if err != nil {
		return nil, err
	}

	if config.Id != 0 {
		if err := config.UnmarshalValue(&newsWebsites); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	slices.SortFunc(newsWebsites, func(a, b *valueobject.NewsWebsite) int {
		return b.Weight - a.Weight
	})

	return newsWebsites, nil
}

// UpdateNewsWebsiteWeight updates the weight of a news website.
func (s *systemConfigService) UpdateNewsWebsiteWeight(ctx context.Context, domain string, step int) error {
	newsWebsites, err := s.GetNewsWebsites(ctx)
	if err != nil {
		return err
	}

	website := gokit.SliceFind(newsWebsites, func(item *valueobject.NewsWebsite) bool {
		return urlx.ExtractSecondLevelDomain(item.Url) == domain
	})

	website.Weight += step

	return s.SaveNewsWebsites(ctx, newsWebsites)
}

// SaveNewsWebsites saves the news websites.
func (s *systemConfigService) SaveNewsWebsites(ctx context.Context, newsWebsites []*valueobject.NewsWebsite) error {
	if len(newsWebsites) == 0 {
		return nil
	}

	slices.SortFunc(newsWebsites, func(a, b *valueobject.NewsWebsite) int {
		return b.Weight - a.Weight
	})

	config, err := entity.NewSystemConfig(valueobject.NewsWebsiteKey.String(), newsWebsites)
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(s.SaveSystemConfig(ctx, config))
}
