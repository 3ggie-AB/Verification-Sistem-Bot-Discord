package models

import "time"

type Coupon struct {
	ID           uint      `gorm:"primaryKey"`
	Code         string    `gorm:"uniqueIndex;size:50"`
	Type         string    // percent | fixed
	Value        float64   // 10 (%) atau 50000 (rupiah)
	MaxDiscount  *float64  // optional (buat percent)
	Quota        uint      // berapa kali bisa dipakai
	UsedCount    uint
	Trigger      *string
	ExpiredAt    *time.Time
	IsActive     bool
	MinMonth     uint     `gorm:"default:0"`
	CreatedAt    time.Time
}
