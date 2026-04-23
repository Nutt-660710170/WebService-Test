package domain

import (
	"time"
)

// Subscription represents a developer's subscription status for API usage
// Each subscription is linked to an APIKey (one-to-many)
type Subscription struct {
	ID         uint      `gorm:"primaryKey"`
	APIKeyID   uint      `gorm:"not null"`
	APIKey     APIKey    `gorm:"foreignKey:APIKeyID"`
	Status     string    `gorm:"size:20;not null"` // active, expired, cancelled, etc.
	StartAt    time.Time `gorm:"not null"`
	EndAt      time.Time `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
