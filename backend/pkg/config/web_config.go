//go:build web

package config

import "os"

// WebConfig is the configuration for the web application.
type WebConfig struct {
	Host    string
	DBAddr  string
	LogFile string
}

// NewWebConfig creates a new WebConfig.
func NewWebConfig() *WebConfig {
	var (
		config = &WebConfig{
			Host:   "0.0.0.0:9010",
			DBAddr: "host=localhost port=5432 user=world_news password=world_news dbname=world_news sslmode=disable TimeZone=Asia/Shanghai",
		}
		host    = os.Getenv("WORLD_NEWS_HOST")
		dBAddr  = os.Getenv("WORLD_NEWS_DB_ADDR")
		LogFile = os.Getenv("WORLD_NEWS_LOG_FILE")
	)

	if host != "" {
		config.Host = host
	}

	if dBAddr != "" {
		config.DBAddr = dBAddr
	}

	if LogFile != "" {
		config.LogFile = LogFile
	}

	return config
}
