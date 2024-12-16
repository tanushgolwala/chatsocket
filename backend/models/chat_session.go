package models

import "time"

type ChatSession struct {
	ID        uint `gorm:"primaryKey"`
	User1ID   uint
	User2ID   uint
	CreatedAt time.Time
}
