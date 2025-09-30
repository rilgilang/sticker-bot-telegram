package model

import "time"

type Sticker struct {
	ID             string    `json:"id"`
	TelegramFileId string    `json:"telegram_file_id"`
	Word           string    `json:"word"`
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at," sql:"DEFAULT:CURRENT_TIMESTAMP"`
	CreatedAt      time.Time `gorm:"not null" json:"created_at," sql:"DEFAULT:CURRENT_TIMESTAMP"`
}
