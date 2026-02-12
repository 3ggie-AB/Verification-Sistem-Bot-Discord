package models

import (
	"time"
	"gorm.io/datatypes"
)

type Announcement struct {
	ID string `gorm:"primaryKey"`

	Title   string
	Content string `gorm:"type:longtext"`
	Type    string

	Target datatypes.JSON `gorm:"type:json"`

	CreatedBy *string
	CreatedAt time.Time
	UpdatedAt time.Time
}