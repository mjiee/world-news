package model

// NewsKeywords
type NewsKeywords struct {
	Id      uint   `gorm:"primaryKey"`
	Keyword string `gorm:"type:varchar(255);unique;NOT NULL"`
}
