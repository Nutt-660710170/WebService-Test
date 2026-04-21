package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"oil/internal/handler"
	"oil/internal/repository"
	"oil/internal/scraper"
	"oil/internal/usecase"
	"oil/pkg/config"
	"oil/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize layers (dependency injection)
	repo := repository.NewOilPriceRepository(db)
	apiScraper := scraper.NewAPIScraper(cfg.BangchakAPIURL)
	uc := usecase.NewOilPriceUsecase(repo, apiScraper)
	h := handler.NewOilPriceHandler(uc)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Oil Price API v1.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Register API routes
	h.RegisterRoutes(app)

	// Setup cron job — runs daily at 1:00 AM
	c := cron.New()
	_, err = c.AddFunc("0 1 * * *", func() {
		uc.CronScrape()
	})
	if err != nil {
		log.Fatalf("Failed to setup cron job: %v", err)
	}
	c.Start()
	log.Println("[CronJob] Scheduled daily oil price scrape at 1:00 AM")

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down...")
		c.Stop()
		_ = app.Shutdown()
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("[Server] Starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
