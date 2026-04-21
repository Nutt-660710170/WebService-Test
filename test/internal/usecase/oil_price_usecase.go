package usecase

import (
	"log"

	"oil/internal/domain"
	"oil/internal/scraper"
)

// OilPriceUsecase defines the interface for oil price business logic.
type OilPriceUsecase interface {
	// ListOilPrices returns oil prices based on the given filter.
	ListOilPrices(filter domain.ListFilter) ([]domain.OilPrice, error)

	// PullOilPrices triggers an immediate scrape from the API, saves to DB, and returns results.
	PullOilPrices() ([]domain.OilPrice, error)

	// CronScrape is called by the cron job to scrape and save prices.
	CronScrape()
}

type oilPriceUsecase struct {
	repo       domain.OilPriceRepository
	apiScraper *scraper.APIScraper
}

// NewOilPriceUsecase creates a new use case with injected dependencies.
func NewOilPriceUsecase(
	repo domain.OilPriceRepository,
	apiScraper *scraper.APIScraper,
) OilPriceUsecase {
	return &oilPriceUsecase{
		repo:       repo,
		apiScraper: apiScraper,
	}
}

// ListOilPrices returns oil prices filtered by the given criteria.
func (uc *oilPriceUsecase) ListOilPrices(filter domain.ListFilter) ([]domain.OilPrice, error) {
	return uc.repo.FindByFilter(filter)
}

// PullOilPrices triggers an immediate scrape, saves, and returns latest prices.
func (uc *oilPriceUsecase) PullOilPrices() ([]domain.OilPrice, error) {
	// Scrape from API
	prices, err := uc.apiScraper.Scrape()
	if err != nil {
		return nil, err
	}

	// Save to database
	if err := uc.repo.Save(prices); err != nil {
		return nil, err
	}

	// Return the latest prices from DB
	return nil, err
}

// CronScrape is the cron job handler — scrapes from API and logs the result.
func (uc *oilPriceUsecase) CronScrape() {
	log.Println("[CronJob] Starting scheduled oil price scrape...")

	prices, err := uc.apiScraper.Scrape()
	if err != nil {
		log.Printf("[CronJob] ERROR: Failed to scrape from API: %v", err)
		return
	}

	if err := uc.repo.Save(prices); err != nil {
		log.Printf("[CronJob] ERROR: Failed to save prices: %v", err)
		return
	}

	log.Printf("[CronJob] Successfully scraped and saved %d oil prices", len(prices))
}
