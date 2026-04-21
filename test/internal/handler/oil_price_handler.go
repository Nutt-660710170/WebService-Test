package handler

import (
	"fmt"
	"log"

	"oil/internal/domain"
	"oil/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

// Response is the standard API response wrapper.
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// OilPriceHandler handles HTTP requests for oil price endpoints.
type OilPriceHandler struct {
	usecase usecase.OilPriceUsecase
}

// NewOilPriceHandler creates a new handler with the injected use case.
func NewOilPriceHandler(uc usecase.OilPriceUsecase) *OilPriceHandler {
	return &OilPriceHandler{usecase: uc}
}

// RegisterRoutes sets up the API routes on the Fiber app.
func (h *OilPriceHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Get("/list", h.ListOilPrices)
	api.Post("/pull", h.PullOilPrices)
}

// ListOilPrices handles POST /api/v1/list
// Body (all optional): {"oil_type": "...", "date_from": "DD/MM/YYYY", "date_to": "DD/MM/YYYY"}
// Returns today's prices by default when no date filter is provided.
func (h *OilPriceHandler) ListOilPrices(c *fiber.Ctx) error {
	var filter domain.ListFilter

	// Parse body if present (all fields are optional)
	if err := c.QueryParser(&filter); err != nil {
		// If body is empty or invalid, use empty filter (defaults to today)
		log.Printf("[Handler] Query parse info: %v (using default filter)", err)
	}

	prices, err := h.usecase.ListOilPrices(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
	}

	return c.JSON(Response{
		Success: true,
		Data:    prices,
		Message: "OK",
	})
}

func (h *OilPriceHandler) PullOilPrices(c *fiber.Ctx) error {
	var prices []domain.OilPrice
	var err error

	prices, err = h.usecase.PullOilPrices()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
	}

	return c.JSON(Response{
		Success: true,
		Data:    prices,
		Message: fmt.Sprintf("Oil prices pulled and saved successfully"),
	})
}
