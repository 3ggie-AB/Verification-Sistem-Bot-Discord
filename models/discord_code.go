package models

import "time"

type DiscordCode struct {
	ID        uint      `gorm:"primaryKey"`
	PaymentID uint      `gorm:"index;not null"`
	Code      string     `gorm:"type:varchar(191);unique;not null"`
	IsUsed    bool       `gorm:"default:false"`
	UsedAt    *time.Time
	CreatedAt time.Time
}
