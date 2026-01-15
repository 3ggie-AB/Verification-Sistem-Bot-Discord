package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type ModuleGroup struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`

	Title       string     `gorm:"not null"`
	Description *string
	IsActive    bool       `gorm:"default:true"`

	Modules []Module `gorm:"foreignKey:ModuleGroupID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *ModuleGroup) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}

type Module struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`

	ModuleGroupID uuid.UUID `gorm:"type:char(36);index"`
	ModuleGroup   ModuleGroup `gorm:"constraint:OnDelete:CASCADE"`

	Title       string     `gorm:"not null"`
	Description *string

	YoutubeID   string     `gorm:"not null"`
	IsActive    bool       `gorm:"default:true"`
	PublishedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Module) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}

type ModuleProgress struct {
	ID uint `gorm:"primaryKey"`

	UserID   uint      `gorm:"index;not null"`
	ModuleID uuid.UUID `gorm:"type:char(36);index;not null"`

	Status string `gorm:"type:enum('watching','completed');default:'watching'"`

	LastWatchedAt *time.Time
	CompletedAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time

	// ❗ cegah duplikat progress
	// 1 user = 1 progress per module
	// composite unique
	// mysql compatible
	_ struct{} `gorm:"uniqueIndex:uniq_user_module"`
}
