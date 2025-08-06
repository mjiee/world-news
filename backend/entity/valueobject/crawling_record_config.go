package valueobject

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// CrawlingRecordConfig represents the configuration for a crawling record.
type CrawlingRecordConfig struct {
	Sources []*NewsWebsite `json:"sources,omitempty"`
	Topics  []string       `json:"topics,omitempty"`
}

// NewCrawlingRecordConfig creates a new CrawlingRecordConfig.
func NewCrawlingRecordConfig(sources []*NewsWebsite, topics []string) *CrawlingRecordConfig {
	return &CrawlingRecordConfig{
		Sources: sources,
		Topics:  topics,
	}
}

// NewCrawlingRecordConfigFromModel creates a new CrawlingRecordConfig from a model.
func NewCrawlingRecordConfigFromModel(data string) (*CrawlingRecordConfig, error) {
	if data == "" {
		return &CrawlingRecordConfig{}, nil
	}

	var config CrawlingRecordConfig

	err := json.Unmarshal([]byte(data), &config)

	return &config, errors.WithStack(err)
}

// ToModel converts the CrawlingRecordConfig to a model.
func (c *CrawlingRecordConfig) ToModel() (string, error) {
	if c == nil {
		return "", nil
	}

	data, err := json.Marshal(c)

	return string(data), errors.WithStack(err)
}
