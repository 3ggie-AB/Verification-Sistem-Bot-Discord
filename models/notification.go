package models

import "time"

type Notification struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // info, warning, error, success

	// siapa targetnya (optional kalau broadcast)
	UserID    *uint   `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`
}