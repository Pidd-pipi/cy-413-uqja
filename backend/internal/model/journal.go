package model

import "time"

type Journal struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	MoodLevel int       `json:"mood_level"`
	Weather   string    `json:"weather"`
	Prompt    string    `gorm:"type:text" json:"prompt"`
	IsPrivate bool      `gorm:"default:true" json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
