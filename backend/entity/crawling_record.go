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
	Config     *valueobject.CrawlingRecordConfig
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewCrawlingRecord creates a new CrawlingRecord entity.
func NewCrawlingRecord(recordType valueobject.CrawlingRecordType,
	config *valueobject.CrawlingRecordConfig) *CrawlingRecord {
	return &CrawlingRecord{
		RecordType: recordType,
		Date:       time.Now(),
		Quantity:   0,
		Status:     valueobject.ProcessingCrawlingRecord,
		Config:     config,
	}
}

// NewCrawlingRecordFromModel converts a CrawlingRecordModel to a CrawlingRecord entity.
func NewCrawlingRecordFromModel(m *model.CrawlingRecord) (*CrawlingRecord, error) {
	if m == nil {
		return nil, errorx.CrawlingRecordNotFound
	}

	config, err := valueobject.NewCrawlingRecordConfigFromModel(m.Config)
	if err != nil {
		return nil, err
	}

	return &CrawlingRecord{
		Id:         m.ID,
		RecordType: valueobject.CrawlingRecordType(m.RecordType),
		Date:       m.Date,
		Quantity:   m.Quantity,
		Status:     valueobject.CrawlingRecordStatus(m.Status),
		Config:     config,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}

// ToModel converts the CrawlingRecord entity to a CrawlingRecordModel.
func (c *CrawlingRecord) ToModel() (*model.CrawlingRecord, error) {
	if c == nil {
		return nil, errorx.CrawlingRecordNotFound
	}

	config, err := c.Config.ToModel()
	if err != nil {
		return nil, err
	}

	return &model.CrawlingRecord{
		ID:         c.Id,
		RecordType: string(c.RecordType),
		Date:       c.Date,
		Quantity:   c.Quantity,
		Status:     string(c.Status),
		Config:     config,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  time.Now(),
	}, nil
}

// CrawlingFailed set the crawling record status to failed.
func (c *CrawlingRecord) CrawlingFailed() {
	c.Status = valueobject.FailedCrawlingRecord
}

// CrawlingPaused set the crawling record status to paused.
func (c *CrawlingRecord) CrawlingPaused() {
	c.Status = valueobject.PausedCrawlingRecord
}

// CrawlingCompleted set the crawling record status to completed.
func (c *CrawlingRecord) CrawlingCompleted() {
	if !c.Status.IsProcessing() {
		return
	}

	c.Status = valueobject.CompletedCrawlingRecord
}
