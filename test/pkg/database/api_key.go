package database

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"

	"oil/internal/domain"
)

// GenerateAPIKey creates a new random API key string
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32) // 256 bits
	if _, err := crypto_rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SaveAPIKey saves a new API key to the database
func SaveAPIKey(db *gorm.DB, owner string, expiredAt *time.Time) (*domain.APIKey, error) {
	key, err := GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}
	apiKey := &domain.APIKey{
		Key:       key,
		Owner:     owner,
		Status:    "active",
		CreatedAt: time.Now(),
		ExpiredAt: expiredAt,
	}
	if err := db.Create(apiKey).Error; err != nil {
		return nil, fmt.Errorf("failed to save API key: %w", err)
	}
	return apiKey, nil
}
