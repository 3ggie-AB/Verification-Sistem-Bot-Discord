package models

import (
	"time"
)


type RulePricing struct {
	ID uint `gorm:"primaryKey"`

	MinMonth int
	MaxMonth *int

	TotalPrice float64
	IsActive   bool

	CreatedAt time.Time
	UpdatedAt time.Time
}