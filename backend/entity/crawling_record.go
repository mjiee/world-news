package entity

import (
	"time"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/repository/model"
)

// CrawlingRecord represents a crawling record.
type CrawlingRecord struct {
	Id        uint
	Date      string
	Quantity  int64
	Status    string
	CreatedAt time.Time
}

// NewCrawlingRecordFromModel converts a CrawlingRecordModel to a CrawlingRecord entity.
func NewCrawlingRecordFromModel(m *model.CrawlingRecord) (*CrawlingRecord, error) {
	if m == nil {
		return nil, errorx.CrawlingRecordNotFound
	}
	return &CrawlingRecord{
		Id:        m.Id,
		Date:      m.Date,
		Quantity:  m.Quantity,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}, nil
}

// ToModel converts the CrawlingRecord entity to a CrawlingRecordModel.
func (c *CrawlingRecord) ToModel() (*model.CrawlingRecord, error) {
	if c == nil {
		return nil, errorx.CrawlingRecordNotFound
	}

	return &model.CrawlingRecord{
		Id:        c.Id,
		Date:      c.Date,
		Quantity:  c.Quantity,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
	}, nil
}
