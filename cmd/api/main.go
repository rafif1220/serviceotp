package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	
	"github.com/rafif1220/serviceotp/internal/config"
)

func main() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Gagal load file .env, pastikan file ada di root directory")
	}

	// 2. Test Koneksi Database & Redis
	// Kalau gagal, fungsi ini bakal log.Fatal dan nge-stop aplikasi
	db := config.InitMySQL()
	redisClient := config.InitRedis()
	
	// (Sementara kita ignore dulu variabel db & redisClient biar gak error 'unused variable')
	_ = db
	_ = redisClient

	// 3. Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		AppName: "Nukar OTP Service v1.0",
	})

	// 4. Bikin route buat Health Check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Pong! Server API OTP Nukar siap tempur 🚀",
		})
	})

	// 5. Jalanin Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Bismillah, Server running on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}