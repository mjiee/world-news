package repository

import (
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/repository/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CrawlingRecordRepository is interface for crawling record.
type CrawlingRecordRepository interface {
}

type crawlingRecordRepository struct {
	db *gorm.DB
}

func NewCrawlingRecordRepository(db *gorm.DB) CrawlingRecordRepository {
	return &crawlingRecordRepository{db}
}

// CreateRecord
func (r *crawlingRecordRepository) CreateRecord(data *model.CrawlingRecord) error {
	return errors.WithStack(r.db.Create(&data).Error)
}

// DeleteRecord
func (r *crawlingRecordRepository) DeleteRecord(id uint) error {
	return errors.WithStack(r.db.Delete(&model.CrawlingRecord{Id: id}).Error)
}

// GetRecords
func (r *crawlingRecordRepository) GetRecords(page *httpx.Pagination) ([]*model.CrawlingRecord, int64, error) {
	var (
		records []*model.CrawlingRecord
		total   int64
	)

	if err := r.db.Model(&model.CrawlingRecord{}).Count(&total).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if err := r.db.Model(&model.CrawlingRecord{}).
		Limit(int(page.Limit)).Offset(int(page.GetOffset())).
		Find(&records).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return records, total, nil

}
