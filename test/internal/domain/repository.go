package domain

// OilPriceRepository defines the interface for oil price data access.
type OilPriceRepository interface {
	// Save upserts a batch of oil prices (insert or update on conflict).
	Save(prices []OilPrice) error

	// FindByFilter returns oil prices matching the given filter criteria.
	FindByFilter(filter ListFilter) ([]OilPrice, error)
}
