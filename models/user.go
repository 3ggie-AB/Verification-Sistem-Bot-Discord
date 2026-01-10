package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"type:varchar(20);default:user"` // default role = user

	MemberExpiredAt *time.Time
	NamaLengkap *string
	IDDiscord *string `gorm:"unique"`
	NamaDiscord *string
	NomorHp *string
	From *string
	Token *string    `gorm:"unique"` 

	CreatedAt time.Time
	UpdatedAt time.Time
}
