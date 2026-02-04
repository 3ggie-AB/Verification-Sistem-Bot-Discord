package models

type CryptoNews struct {
	ArticleID   string `gorm:"primaryKey;size:100"`
	SourceName  string
	SourceDomain string
	Thumbnail   string
	URL         string
	Title       string
	Description string `gorm:"type:text"`
	Content     string `gorm:"type:text"`
	ContentIndo string `gorm:"type:text"`
	IsUpload 	bool   `gorm:"default:false"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
}
