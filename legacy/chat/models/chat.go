package models

import "time"

type Chat struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
