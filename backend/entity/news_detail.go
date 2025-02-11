package entity

import (
	"encoding/json"
	"time"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/repository/model"

	"github.com/pkg/errors"
)

// NewsDetail represents the detailed information about a news item.
type NewsDetail struct {
	Id          uint
	RecordId    uint // crawling record id
	Title       string
	Link        string
	Contents    []string
	Images      []string
	PublishedAt time.Time
	CreatedAt   time.Time
}

// NewNewsDetailFromModel converts a NewsDetailModel to a NewsDetail entity.
func NewNewsDetailFromModel(m *model.NewsDetail) (*NewsDetail, error) {
	if m == nil {
		return nil, errorx.NewsNotFound
	}

	var (
		contents []string
		images   []string
	)

	if err := json.Unmarshal([]byte(m.Contents), &contents); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal([]byte(m.Images), &images); err != nil {
		return nil, errors.WithStack(err)
	}

	return &NewsDetail{
		Id:          m.ID,
		RecordId:    m.RecordId,
		Title:       m.Title,
		Link:        m.Link,
		Contents:    contents,
		Images:      images,
		PublishedAt: m.PublishedAt,
		CreatedAt:   m.CreatedAt,
	}, nil
}

// ToModel converts the NewsDetail entity to a NewsDetailModel.
func (n *NewsDetail) ToModel() (*model.NewsDetail, error) {
	if n == nil {
		return nil, errorx.NewsNotFound
	}

	contents, err := json.Marshal(n.Contents)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	images, err := json.Marshal(n.Images)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &model.NewsDetail{
		ID:          n.Id,
		RecordId:    n.RecordId,
		Title:       n.Title,
		Link:        n.Link,
		Contents:    string(contents),
		Images:      string(images),
		PublishedAt: n.PublishedAt,
		CreatedAt:   n.CreatedAt,
	}, nil
}

// Validate validates the NewsDetail entity.
func (n *NewsDetail) Validate() error {
	if n == nil {
		return errorx.NewsNotFound
	}

	if n.Title == "" { // title is required
		return errorx.NewsNotFound
	}

	if len(n.Contents) == 0 { // contents is required
		return errorx.NewsNotFound
	}

	return nil
}
