package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"oil/internal/domain"
)

// apiOilPriceResponse represents the outer JSON response from Bangchak API.
type apiOilPriceResponse struct {
	OilPriceID   int    `json:"OilPriceID"`
	OilDateNow   string `json:"OilDateNow"`
	OilPriceDate string `json:"OilPriceDate"`
	OilPriceTime string `json:"OilPriceTime"`
	OilRemark    string `json:"OilRemark"`
	OilRemark2   string `json:"OilRemark2"`
	OilList      string `json:"OilList"` // API sends this as a stringified JSON array
}

// oilListItem represents a single oil type within the OilList JSON string.
type oilListItem struct {
	OilName           string  `json:"OilName"`
	PriceYesterday    float64 `json:"PriceYesterday"`
	PriceToday        float64 `json:"PriceToday"`
	PriceTomorrow     float64 `json:"PriceTomorrow"`
	PriceDifYesterday float64 `json:"PriceDifYesterday"`
	PriceDifTomorrow  float64 `json:"PriceDifTomorrow"`
}

// APIScraper fetches oil prices from Bangchak's API endpoint.
type APIScraper struct {
	apiURL string
}

// NewAPIScraper creates a new API scraper with the given endpoint URL.
func NewAPIScraper(apiURL string) *APIScraper {
	return &APIScraper{apiURL: apiURL}
}

// Scrape fetches oil prices from the Bangchak API and returns domain entities.
func (s *APIScraper) Scrape() ([]domain.OilPrice, error) {
	log.Printf("[APIScraper] Fetching oil prices from: %s", s.apiURL)

	resp, err := http.Get(s.apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse outer JSON array
	var apiResponses []apiOilPriceResponse
	if err := json.Unmarshal(body, &apiResponses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	if len(apiResponses) == 0 {
		return nil, fmt.Errorf("API returned empty response")
	}

	apiResp := apiResponses[0]

	// Parse the date from OilPriceDate
	// Format: "DD/MM/YYYY" (CE)
	priceDate, err := time.Parse("02/01/2006", apiResp.OilPriceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OilPriceDate %q: %w", apiResp.OilPriceDate, err)
	}

	// Double-unmarshal: OilList is a JSON string containing an array
	var oilItems []oilListItem
	if err := json.Unmarshal([]byte(apiResp.OilList), &oilItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OilList: %w", err)
	}

	// Map to domain entities
	var prices []domain.OilPrice
	now := time.Now().Format("2006-01-02")

	for _, item := range oilItems {
		prices = append(prices, domain.OilPrice{
			OilType:   item.OilName,
			Price:     item.PriceToday,
			Date:      priceDate,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	log.Printf("[APIScraper] Successfully scraped %d oil prices", len(prices))
	return prices, nil
}
