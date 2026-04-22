package domain

import (
	"time"
)

// APIKey represents an API key for authentication
type APIKey struct {
	ID        uint      `gorm:"primaryKey"`
	Key       string    `gorm:"uniqueIndex;size:64;not null"`
	Owner     string    `gorm:"size:100"` // optional: for identifying the dev
	Status    string    `gorm:"size:20;default:'active'"`
	CreatedAt time.Time
	ExpiredAt *time.Time
}
