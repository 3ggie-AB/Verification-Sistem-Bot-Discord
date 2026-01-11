package models

import "time"

type Payment struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"index;not null"`
	Amount         float64
	OriginalAmount float64
	Method         string    // crypto | ewallet | bank
	Status         string    // pending | paid | failed
	TransactionRef *string
	MonthCount 	   uint
	CouponID       *uint
	Discount       float64
	Bukti 	       *string
	RejectReason   *string
	PaidAt         *time.Time
	CreatedAt      time.Time
	DiscordCode    *DiscordCode `gorm:"foreignKey:PaymentID"`
}
