package models

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   uint      `gorm:"index"`
	ReceiverID uint      `gorm:"index"`
	Content    string    `gorm:"type:text"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
}
