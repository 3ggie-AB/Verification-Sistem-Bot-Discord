package models

import "time"

type Expense struct {
	ID          uint      `gorm:"primaryKey"`
	Description string    `gorm:"size:255"`  // apa pengeluarannya
	Amount      float64   // nominal pengeluaran
	Category    *string   `gorm:"size:100"`  // optional, misal: operasional, marketing, dll
	SpentAt     time.Time // kapan pengeluaran terjadi
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
