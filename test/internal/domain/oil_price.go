package domain

import "time"

// OilPrice represents the core entity for oil price data.
type OilPrice struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Date      time.Time `json:"date" gorm:"type:date;not null"`
	OilType   string    `json:"oil_type" gorm:"type:varchar(100);not null"`
	Price     float64   `json:"price" gorm:"not null"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// ListFilter represents the filtering options for listing oil prices.
type ListFilter struct {
	OilType string `json:"oil_type" query:"oil_type"` // optional: filter by oil type name
}
