package models

import (
	"time"
	"gorm.io/datatypes"
)

type AutoMessager struct {
	ID string `gorm:"primaryKey"`

	Name    string
	Message string `gorm:"type:longtext"`
	Image   *string

	BotID    string
	ServerID string
	ChannelID *string

	RunTime *string `gorm:"type:varchar(5)"` 
	// contoh: "08:30"

	DaysOfWeek datatypes.JSON 
	// contoh: ["mon","tue","fri"]

	Timezone string `gorm:"default:Asia/Jakarta"`

	IsActive bool

	LastRunAt *time.Time
	NextRunAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Bot struct {
	ID string `gorm:"primaryKey"`

	Name  string
	Token string

	IsActive bool

	CreatedAt time.Time
}