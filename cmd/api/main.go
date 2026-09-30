package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/rafif1220/serviceotp/internal/config"
	"github.com/rafif1220/serviceotp/internal/delivery/http"
	"github.com/rafif1220/serviceotp/internal/delivery/middleware"
	"github.com/rafif1220/serviceotp/internal/provider/barantum"
	"github.com/rafif1220/serviceotp/internal/provider/smtp"
	mysqlrepo "github.com/rafif1220/serviceotp/internal/repository/mysql"
	redisrepo "github.com/rafif1220/serviceotp/internal/repository/redis"
	"github.com/rafif1220/serviceotp/internal/usecase"
)

func main() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Gagal load file .env")
	}

	// 2. Inisialisasi Koneksi Database & Redis
	db := config.InitMySQL()
	redisClient := config.InitRedis()

	// 3. Inisialisasi Repository Layer
	mysqlRepo := mysqlrepo.NewOTPRepository(db)
	redisRepo := redisrepo.NewCacheRepository(redisClient)

	// 4. Inisialisasi Provider Layer
	barantumProvider := barantum.NewBarantumProvider()
	smtpProvider := smtp.NewSMTPProvider()

	// 5. Inisialisasi Usecase Layer (Suntik semua dependency ke sini)
	otpUsecase := usecase.NewOtpUseCase(
		mysqlRepo,
		redisRepo,
		barantumProvider,
		smtpProvider,
	)

	// 6. Inisialisasi Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Nukar OTP Service v1.0",
	})

	// Route Health Check (Tanpa Middleware)
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "success", "message": "Pong!"})
	})

	// 7. Setup Route dengan Middleware
	// Bikin grup /api/v1 biar rapi
	api := app.Group("/api/v1")
	
	// Bikin grup /otp yang dijagain sama Middleware RequireAPIKey
	otpGroup := api.Group("/otp", middleware.RequireAPIKey())

	// 8. Inisialisasi Handler dan daftarin routenya
	otpHandler := http.NewOtpHandler(otpUsecase)
	otpHandler.Route(otpGroup) // Endpoint jadinya: POST /api/v1/otp/request

	// 9. Jalanin Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Bismillah, Server API OTP running on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}