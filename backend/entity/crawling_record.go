package entity

import (
	"time"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/repository/model"
)

// CrawlingRecord represents a crawling record.
type CrawlingRecord struct {
	Id         uint
	RecordType valueobject.CrawlingRecordType
	Date       time.Time
	Quantity   int64
	Status     valueobject.CrawlingRecordStatus
	CreatedAt  time.Time
}

// NewCrawlingRecord creates a new CrawlingRecord entity.
func NewCrawlingRecord(recordType valueobject.CrawlingRecordType) *CrawlingRecord {
	return &CrawlingRecord{
		RecordType: recordType,
		Date:       time.Now(),
		Quantity:   0,
		Status:     valueobject.ProcessingCrawlingRecord,
	}
}

// NewCrawlingRecordFromModel converts a CrawlingRecordModel to a CrawlingRecord entity.
func NewCrawlingRecordFromModel(m *model.CrawlingRecord) (*CrawlingRecord, error) {
	if m == nil {
		return nil, errorx.CrawlingRecordNotFound
	}
	return &CrawlingRecord{
		Id:         m.Id,
		RecordType: valueobject.CrawlingRecordType(m.RecordType),
		Date:       m.Date,
		Quantity:   m.Quantity,
		Status:     valueobject.CrawlingRecordStatus(m.Status),
		CreatedAt:  m.CreatedAt,
	}, nil
}

// ToModel converts the CrawlingRecord entity to a CrawlingRecordModel.
func (c *CrawlingRecord) ToModel() (*model.CrawlingRecord, error) {
	if c == nil {
		return nil, errorx.CrawlingRecordNotFound
	}

	return &model.CrawlingRecord{
		Id:         c.Id,
		RecordType: string(c.RecordType),
		Date:       c.Date,
		Quantity:   c.Quantity,
		Status:     string(c.Status),
		CreatedAt:  c.CreatedAt,
	}, nil
}

// CrawlingFailed set the crawling record status to failed.
func (c *CrawlingRecord) CrawlingFailed() {
	c.Status = valueobject.FailedCrawlingRecord
}

// CrawlingCompleted set the crawling record status to completed.
func (c *CrawlingRecord) CrawlingCompleted() {
	if c.Status.IsFailed() {
		return
	}

	c.Status = valueobject.CompletedCrawlingRecord
}
