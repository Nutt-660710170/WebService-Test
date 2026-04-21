package repository

import (
	"oil/internal/domain"

	"gorm.io/gorm"
)

type oilPriceRepository struct {
	db *gorm.DB
}

// NewOilPriceRepository creates a new GORM-backed repository.
func NewOilPriceRepository(db *gorm.DB) domain.OilPriceRepository {
	return &oilPriceRepository{db: db}
}

// Save inserts new oil prices into the database.
func (r *oilPriceRepository) Save(prices []domain.OilPrice) error {
	if len(prices) == 0 {
		return nil
	}
	return r.db.Create(&prices).Error
}

// FindByFilter returns oil prices matching filter criteria, limited to the most recent batch.
func (r *oilPriceRepository) FindByFilter(filter domain.ListFilter) ([]domain.OilPrice, error) {
	var prices []domain.OilPrice

	// 1. Find the latest record to identify the most recent batch's timestamp
	var latestPrice domain.OilPrice
	err := r.db.Order("created_at DESC").First(&latestPrice).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []domain.OilPrice{}, nil
		}
		return nil, err
	}

	// 2. Query for all prices that share that exact latest timestamp
	query := r.db.Where("created_at = ?", latestPrice.CreatedAt)

	// 3. Apply additional filters if provided
	if filter.OilType != "" {
		query = query.Where("oil_type = ?", filter.OilType)
	}

	err = query.Order("oil_type ASC").Find(&prices).Error
	return prices, err
}
